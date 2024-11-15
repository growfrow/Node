package mirror

import (
	"context"
	"fmt"

	"github.com/gabriel-vasile/mimetype"
	"github.com/rss3-network/node/config"
	"github.com/rss3-network/node/internal/database"
	"github.com/rss3-network/node/internal/engine"
	source "github.com/rss3-network/node/internal/engine/protocol/arweave"
	"github.com/rss3-network/node/internal/engine/worker/decentralized/contract/mirror/model"
	"github.com/rss3-network/node/provider/arweave"
	"github.com/rss3-network/node/provider/arweave/contract/mirror"
	"github.com/rss3-network/node/provider/ipfs"
	"github.com/rss3-network/node/schema/worker/decentralized"
	"github.com/rss3-network/protocol-go/schema"
	activityx "github.com/rss3-network/protocol-go/schema/activity"
	"github.com/rss3-network/protocol-go/schema/metadata"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
	"github.com/rss3-network/protocol-go/schema/typex"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// make sure worker implements engine.Worker
var _ engine.Worker = (*worker)(nil)

type worker struct {
	config         *config.Module
	arweaveClient  arweave.Client
	ipfsClient     ipfs.HTTPClient
	databaseClient database.Client
}

func (w *worker) Name() string {
	return decentralized.Mirror.String()
}

func (w *worker) Platform() string {
	return decentralized.PlatformMirror.String()
}

func (w *worker) Network() []network.Network {
	return []network.Network{
		network.Arweave,
	}
}

func (w *worker) Tags() []tag.Tag {
	return []tag.Tag{
		tag.Social,
	}
}

func (w *worker) Types() []schema.Type {
	return []schema.Type{
		typex.SocialPost,
		typex.SocialRevise,
	}
}

// Filter returns a filter for protocol.
func (w *worker) Filter() engine.DataSourceFilter {
	return &source.Filter{OwnerAddresses: []string{mirror.AddressMirror}}
}

// Transform returns an activity  with the action of the task.
func (w *worker) Transform(ctx context.Context, task engine.Task) (*activityx.Activity, error) {
	// Cast the task to an Arweave task.
	arweaveTask, ok := task.(*source.Task)
	if !ok {
		return nil, fmt.Errorf("invalid task type: %T", task)
	}

	zap.L().Debug("transforming mirror task", zap.String("task_id", arweaveTask.ID()))

	// Build the activity.
	activity, err := task.BuildActivity(activityx.WithActivityPlatform(w.Platform()))
	if err != nil {
		return nil, fmt.Errorf("build activity: %w", err)
	}

	// Get actions from the transaction.
	actions, err := w.transformMirrorAction(ctx, arweaveTask)
	if err != nil {
		return nil, fmt.Errorf("handle arweave mirror transaction: %w", err)
	}

	activity.To = mirror.AddressMirror

	// Activity type should be inferred from the action (if it's revise)
	if actions != nil {
		activity.Type = actions[0].Type
		activity.Actions = append(activity.Actions, actions...)
	}

	zap.L().Debug("successfully transformed mirror task")

	return activity, nil
}

// transformPostOrReviseAction Returns the actions of mirror post or revise.
func (w *worker) transformMirrorAction(ctx context.Context, task *source.Task) ([]*activityx.Action, error) {
	var (
		contentDigest       string
		originContentDigest string
		emptyOriginDigest   bool
	)

	zap.L().Debug("transforming mirror action", zap.String("transaction_id", task.Transaction.ID))

	for _, tag := range task.Transaction.Tags {
		tagName, err := arweave.Base64Decode(tag.Name)
		if err != nil {
			return nil, fmt.Errorf("base64 decode tag name: %w", err)
		}

		tagValue, err := arweave.Base64Decode(tag.Value)
		if err != nil {
			return nil, err
		}

		switch string(tagName) {
		case "Content-Digest":
			contentDigest = string(tagValue)
			zap.L().Debug("found content digest", zap.String("digest", contentDigest))
		case "Original-Content-Digest":
			originContentDigest = string(tagValue)
			zap.L().Debug("found original content digest", zap.String("digest", originContentDigest))

			if len(string(tagValue)) == 0 {
				emptyOriginDigest = true

				zap.L().Debug("original content digest is empty")
			}
		}
	}

	// Construct content URI from tx id
	contentURI := fmt.Sprintf("ar://%s", task.Transaction.ID)
	zap.L().Debug("constructed content URI", zap.String("uri", contentURI))

	// Get detailed post info from transaction data
	transactionData, err := arweave.Base64Decode(task.Transaction.Data)
	if err != nil {
		return nil, fmt.Errorf("invalid foramt of transaction data: %w", err)
	}

	mirrorData := gjson.ParseBytes(transactionData)

	author := mirrorData.Get("authorship.contributor").String()
	zap.L().Debug("found post author", zap.String("author", author))

	var media []metadata.Media

	// Get mirror nft as media
	address := mirrorData.Get("wnft.imageURI").String()
	if address != "" {
		zap.L().Debug("fetching NFT image", zap.String("address", address))

		file, err := w.ipfsClient.Fetch(ctx, fmt.Sprintf("/ipfs/%s", address), ipfs.FetchModeQuick)
		if err != nil {
			return nil, fmt.Errorf("fetch ipfs: %w", err)
		}

		defer lo.Try(file.Close)

		// Get nft mimetype
		result, err := mimetype.DetectReader(file)
		if err != nil {
			return nil, fmt.Errorf("detect mimetype: %w", err)
		}

		if result == nil {
			return nil, fmt.Errorf("empty result")
		}

		media = append(media, metadata.Media{
			Address:  fmt.Sprintf("ipfs://%s", address),
			MimeType: result.String(),
		})

		zap.L().Debug("added NFT media", zap.String("mimetype", result.String()))
	}

	var publicationID string

	if contentDigest == "" {
		publicationID = mirrorData.Get("digest").String()
		zap.L().Debug("using digest as publication ID", zap.String("publication_id", publicationID))
	} else {
		publicationID = contentDigest
		zap.L().Debug("using content digest as publication ID", zap.String("publication_id", publicationID))
	}

	if originContentDigest != "" {
		publicationID = originContentDigest
		zap.L().Debug("using original content digest as publication ID", zap.String("publication_id", publicationID))
	}

	// Construct mirror Metadata
	mirrorMetadata := &metadata.SocialPost{
		Title:         mirrorData.Get("content.title").String(),
		Body:          mirrorData.Get("content.body").String(),
		ContentURI:    contentURI,
		PublicationID: publicationID,
		Media:         media,
		Timestamp:     mirrorData.Get("content.timestamp").Uint(),
	}
	zap.L().Debug("constructed mirror metadata", zap.String("title", mirrorMetadata.Title))

	// Build the post or revise action
	action, err := w.buildMirrorAction(ctx, task.Transaction.ID, author, mirror.AddressMirror, mirrorMetadata, emptyOriginDigest, originContentDigest)
	if err != nil {
		return nil, fmt.Errorf("build post action: %w", err)
	}

	// Save Dataset Mirror Post
	post := &model.DatasetMirrorPost{
		TransactionID:       task.Transaction.ID,
		OriginContentDigest: originContentDigest,
	}

	if err := w.databaseClient.SaveDatasetMirrorPost(context.TODO(), post); err != nil {
		return nil, fmt.Errorf("save dataset mirror post: %w", err)
	}

	zap.L().Info("saved dataset mirror post", zap.String("transaction_id", task.Transaction.ID))

	actions := []*activityx.Action{
		action,
	}

	zap.L().Debug("successfully transformed mirror action")

	return actions, nil
}

// buildArweaveTransactionTransferAction Returns the native transfer transaction action.
func (w *worker) buildMirrorAction(ctx context.Context, txID, from, to string, mirrorMetadata *metadata.SocialPost, emptyOriginDigest bool, originContentDigest string) (*activityx.Action, error) {
	zap.L().Debug("building mirror action",
		zap.String("transaction_id", txID),
		zap.String("from", from),
		zap.String("to", to),
		zap.Any("metadata", mirrorMetadata),
		zap.Bool("empty_origin_digest", emptyOriginDigest),
		zap.String("origin_content_digest", originContentDigest))

	// Default action type is post.
	filterType :=
		typex.SocialPost

	// if the origin digest is empty, the action type should be revise.
	if emptyOriginDigest {
		zap.L().Debug("empty origin digest detected, setting action type to revise")

		filterType = typex.SocialRevise
	}

	// If the origin digest is not empty, check if the origin digest is the first mirror post.
	if originContentDigest != "" {
		zap.L().Debug("loading dataset mirror post", zap.String("origin_content_digest", originContentDigest))

		post, err := w.databaseClient.LoadDatasetMirrorPost(ctx, originContentDigest)
		if err != nil {
			return nil, fmt.Errorf("load dataset mirror post: %w", err)
		}

		if post != nil && txID != post.TransactionID {
			zap.L().Debug("post exists with different transaction ID, setting action type to revise",
				zap.String("post_transaction_id", post.TransactionID),
				zap.String("current_transaction_id", txID))

			filterType = typex.SocialRevise
		}
	}

	// Construct action
	action := activityx.Action{
		Type:     filterType,
		Tag:      tag.Social,
		Platform: w.Platform(),
		From:     from,
		To:       to,
		Metadata: mirrorMetadata,
	}

	zap.L().Debug("successfully built mirror action")

	return &action, nil
}

// NewWorker returns a new Arweave worker.
func NewWorker(config *config.Module, databaseClient database.Client) (engine.Worker, error) {
	var instance = worker{
		config: config,
	}

	var err error

	if instance.arweaveClient, err = arweave.NewClient(); err != nil {
		return nil, fmt.Errorf("new arweave client: %w", err)
	}

	// Initialize ipfs client.
	if instance.ipfsClient, err = ipfs.NewHTTPClient(ipfs.WithGateways(config.IPFSGateways)); err != nil {
		return nil, fmt.Errorf("new ipfs client: %w", err)
	}

	instance.databaseClient = databaseClient

	return &instance, nil
}
