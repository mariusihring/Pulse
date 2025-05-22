package solana_types

type TokenAccountsByOwnerResponse struct {
	JsonRPC string                        `json:"jsonrpc"`
	Result  GetTokenAccountsByOwnerResult `json:"result"`
}

type GetTokenAccountsByOwnerResult struct {
	Context GetTokenAccountsByOwnerContext `json:"context"`
	Value   []TokenAccount                 `json:"value"`
}

type GetTokenAccountsByOwnerContext struct {
	ApiVersion string `json:"apiVersion"`
	Slot       int64  `json:"slot"`
}

type TokenAccount struct {
	Account Account `json:"account"`
	Pubkey  string  `json:"pubkey"`
}

type Account struct {
	Data       AccountData `json:"data"`
	Executable bool        `json:"executable"`
	Lamports   int64       `json:"lamports"`
	Owner      string      `json:"owner"`
	RentEpoch  uint64      `json:"rentEpoch"`
	Space      int         `json:"space"`
}

type AccountData struct {
	Parsed  ParsedData `json:"parsed"`
	Program string     `json:"program"`
	Space   int        `json:"space"`
}

type ParsedData struct {
	Info AccountInfo `json:"info"`
	Type string      `json:"type"`
}

type AccountInfo struct {
	IsNative    bool        `json:"isNative"`
	Mint        string      `json:"mint"`
	Owner       string      `json:"owner"`
	State       string      `json:"state"`
	TokenAmount TokenAmount `json:"tokenAmount"`
}

type TokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       int     `json:"decimals"`
	UIAmount       float64 `json:"uiAmount"`
	UIAmountString string  `json:"uiAmountString"`
}
