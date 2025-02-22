package types

type CoinGeckoPriceResponse struct {
	Data TokenPriceData `json:"data"`
}

type TokenPriceData struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Attributes TokenPriceAttributes `json:"attributes"`
}

type TokenPriceAttributes struct {
	TokenPrices map[string]string `json:"token_prices"`
}
