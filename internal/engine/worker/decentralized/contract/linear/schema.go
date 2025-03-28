package linear

type FunctionCallArgs struct {
	ReceiverID string `json:"receiver_id"`
	Amount     string `json:"amount"`
	Msg        string `json:"msg"`
}

type Msg struct {
	Force   int64    `json:"force"`
	Actions []Action `json:"actions"`
}

type Action struct {
	PoolID       int64  `json:"pool_id"`
	TokenIn      string `json:"token_in"`
	TokenOut     string `json:"token_out"`
	AmountIn     string `json:"amount_in,omitempty"`
	MinAmountOut string `json:"min_amount_out"`
}

type Event struct {
	Standard     string `json:"standard"`
	Version      string `json:"version"`
	Event        string `json:"event"`
	Data         []Data `json:"data"`
	TokenAddress string `json:"token_address"`
}

type Data struct {
	OldOwnerID string `json:"old_owner_id"`
	NewOwnerID string `json:"new_owner_id"`
	Amount     string `json:"amount"`
}
