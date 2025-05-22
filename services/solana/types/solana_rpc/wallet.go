package solana_types

type WalletResponse struct {
	JsonRPC string          `json:"jsonrpc"`
	Result  GetWalletResult `json:"result"`
	Id      int64           `json:"id"`
}

type GetWalletResult struct {
	Context interface{} `json:"context"`
	Value   int64       `json:"value"`
}
