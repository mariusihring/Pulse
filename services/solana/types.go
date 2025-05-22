package main

type GetBalanceResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		// Value is the balance in lamports (1 SOL = 1e9 lamports).
		Value uint64 `json:"value"`
	} `json:"result"`
	Id int `json:"id"`
}

type GetTokenAccountsResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		Value []struct {
			Pubkey  string `json:"pubkey"`
			Account struct {
				Data struct {
					Parsed struct {
						Info struct {
							Mint        string `json:"mint"`
							Owner       string `json:"owner"`
							TokenAmount struct {
								Amount         string  `json:"amount"`
								Decimals       int     `json:"decimals"`
								UiAmount       float64 `json:"uiAmount"`
								UiAmountString string  `json:"uiAmountString"`
							} `json:"tokenAmount"`
						} `json:"info"`
						Type string `json:"type"`
					} `json:"parsed"`
					Program string `json:"program"`
					Space   int    `json:"space"`
				} `json:"data"`
				Executable bool   `json:"executable"`
				Lamports   int64  `json:"lamports"`
				Owner      string `json:"owner"`
				// Use uint64 to accommodate very large values.
				RentEpoch uint64 `json:"rentEpoch"`
			} `json:"account"`
		} `json:"value"`
	} `json:"result"`
	Id int `json:"id"`
}

type TokenInfo struct {
	Mint         string  `json:"mint"`
	Amount       float64 `json:"amount"`
	Decimals     int     `json:"decimals"`
	TokenAccount string  `json:"token_account"`
	UsdValue     float64 `json:"usd_value"`
}

type BirdeyeTokenPriceResponse struct {
	Data    Data `json:"data"`
	Success bool `json:"success"`
}

type Data struct {
	Value           float64 `json:"value"`
	UpdateUnixTime  int64   `json:"updateUnixTime"`
	UpdateHumanTime string  `json:"updateHumanTime"`
	PriceChange24h  float64 `json:"priceChange24h"`
	PriceInNative   float64 `json:"priceInNative"`
}
