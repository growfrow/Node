package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/naturalselectionlabs/rss3-node/internal/constant"
	"github.com/naturalselectionlabs/rss3-node/internal/database"
	"github.com/naturalselectionlabs/rss3-node/internal/engine"
	"github.com/naturalselectionlabs/rss3-node/internal/engine/source"
	"github.com/naturalselectionlabs/rss3-node/internal/engine/worker"
	"github.com/naturalselectionlabs/rss3-node/schema"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
)

type Server struct {
	id             string
	config         *engine.Config
	source         engine.Source
	worker         engine.Worker
	databaseClient database.Client
}

func (s *Server) Run(ctx context.Context) error {
	var (
		// TODO Develop a more effective solution to implement back pressure.
		tasksChan = make(chan []engine.Task, 1)
		errorChan = make(chan error)
	)

	zap.L().Info("start node", zap.String("version", constant.BuildVersion()))

	s.source.Start(ctx, tasksChan, errorChan)

	for {
		select {
		case tasks := <-tasksChan:
			if err := s.handleTasks(ctx, tasks); err != nil {
				return fmt.Errorf("handle tasks: %w", err)
			}
		case err := <-errorChan:
			if err != nil {
				return fmt.Errorf("an error occurred in the source: %w", err)
			}

			return nil
		}
	}
}

func (s *Server) handleTasks(ctx context.Context, tasks []engine.Task) error {
	if len(tasks) == 0 {
		return nil
	}

	resultPool := pool.NewWithResults[*schema.Feed]().WithMaxGoroutines(lo.Ternary(len(tasks) < 20*runtime.NumCPU(), len(tasks), 20*runtime.NumCPU()))

	for _, task := range tasks {
		task := task

		resultPool.Go(func() *schema.Feed {
			zap.L().Debug("start match task", zap.String("task.id", task.ID()))

			matched, err := s.worker.Match(ctx, task)
			if err != nil {
				zap.L().Error("matched task", zap.String("task.id", task.ID()))

				return nil
			}

			if !matched {
				zap.L().Warn("unmatched task", zap.String("task.id", task.ID()))

				return nil
			}

			zap.L().Debug("start transform task", zap.String("task.id", task.ID()))

			feed, err := s.worker.Transform(ctx, task)
			if err != nil {
				zap.L().Error("transform task", zap.String("task.id", task.ID()))

				return nil
			}

			return feed
		})
	}

	// Filter failed feeds.
	feeds := lo.Filter(resultPool.Wait(), func(feed *schema.Feed, _ int) bool {
		return feed != nil
	})

	// Save feeds and checkpoint to the database.
	return s.databaseClient.WithTransaction(ctx, func(ctx context.Context, client database.Client) error {
		if err := client.SaveFeeds(ctx, feeds); err != nil {
			return fmt.Errorf("save %d feeds: %w", len(feeds), err)
		}

		checkpoint := engine.Checkpoint{
			ID:      s.id,
			Network: s.source.Network(),
			Worker:  s.worker.Name(),
			State:   s.source.State(),
		}

		zap.L().Info("save checkpoint", zap.Any("checkpoint", checkpoint))

		if err := client.SaveCheckpoint(ctx, &checkpoint); err != nil {
			return fmt.Errorf("save checkpoint: %w", err)
		}

		return nil
	})
}

func NewServer(ctx context.Context, config *engine.Config, databaseClient database.Client) (server *Server, err error) {
	instance := Server{
		id:             fmt.Sprintf("%s.%s", config.Network, config.Worker),
		config:         config,
		databaseClient: databaseClient,
	}

	// Initialize worker.
	if instance.worker, err = worker.New(instance.config, databaseClient); err != nil {
		return nil, fmt.Errorf("new worker: %w", err)
	}

	// Load checkpoint for initialize the source.
	checkpoint, err := instance.databaseClient.LoadCheckpoint(ctx, instance.id, config.Network, instance.worker.Name())
	if err != nil {
		return nil, fmt.Errorf("loca checkpoint: %w", err)
	}

	// Unmarshal checkpoint state to map for print it in log.
	var state map[string]any
	if err := json.Unmarshal(checkpoint.State, &state); err != nil {
		return nil, fmt.Errorf("unmarshal checkpoint state: %w", err)
	}

	zap.L().Info("load checkpoint", zap.String("checkpoint.id", checkpoint.ID), zap.String("checkpoint.network", checkpoint.Network.String()), zap.String("checkpoint.worker", checkpoint.Worker), zap.Any("checkpoint.state", state))

	// Initialize source.
	if instance.source, err = source.New(instance.config, checkpoint); err != nil {
		return nil, fmt.Errorf("new source: %w", err)
	}

	return &instance, nil
}
