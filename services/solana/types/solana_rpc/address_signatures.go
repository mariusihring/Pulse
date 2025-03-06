package solana_types

type SignaturesForAddressResponse struct {
	JsonRPC string                          `json:"jsonrpc"`
	Result  []WalletTransactionHashResponse `json:"result"`
	Id      int64                           `json:"id"`
}

type WalletTransactionHashResponse struct {
	Err       *interface{} `json:"err"`
	Memo      string       `json:"memo"`
	Signature string       `json:"signature"`
	Slot      int64        `json:"slot"`
	BlockTime int64        `json:"blockTime"`
}
