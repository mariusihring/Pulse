package solana_types

type AccountInfoResponse struct {
	JsonRPC string               `json:"jsonrpc"`
	Result  GetAccountInfoResult `json:"result"`
	Id      int16                `json:"id"`
}

type GetAccountInfoResult struct {
	Context GetAccountInfoContext `json:"context"`
	Value   GetAccountInfoValue   `json:"value"`
}

type GetAccountInfoContext struct {
	ApiVersion string `json:"apiVersion"`
	Slot       int64  `json:"slot"`
}
type GetAccountInfoValue struct {
	Data       string `json:"data"`
	Executable bool   `json:"executable"`
	Lamports   int64  `json:"lamports"`
	Owner      string `json:"owner"`
	RentEpoch  uint64 `json:"rentEpoch"`
	Space      int64  `json:"space"`
}
