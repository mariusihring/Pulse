package types

type PriceResponse struct {
	Solana struct {
		USD float64 `json:"usd"`
	} `json:"solana"`
}

type TokenBalance struct {
	Currency         Currency `json:"currency"`
	ConfirmedBalance string   `json:"confirmed_balance"`
	ConfirmedBlock   int64    `json:"confirmed_block"`
}

type Currency struct {
	AssetPath string `json:"asset_path"`
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	Decimals  int    `json:"decimals"`
	Type      string `json:"type"`
}
