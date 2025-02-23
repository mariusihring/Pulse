package coingecko_types

// OHLCVSResponse represents the entire JSON response.
type OHLCVSResponse struct {
	Data OHLCVData `json:"data"`
	Meta OHLCVMeta `json:"meta"`
}

// OHLCVData holds the id, type, and attributes for the OHLCV request.
type OHLCVData struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes OHLCVAttributes `json:"attributes"`
}

// OHLCVAttributes contains the list of OHLCV records.
type OHLCVAttributes struct {
	// Each inner slice represents a record in the order:
	// [timestamp, open, high, low, close, volume]
	OHLCVList [][]float64 `json:"ohlcv_list"`
}

// OHLCVMeta holds metadata for the response.
type OHLCVMeta struct {
	Base  CoinMeta `json:"base"`
	Quote CoinMeta `json:"quote"`
}

// CoinMeta represents the coin metadata in the meta section.
type CoinMeta struct {
	Address         string `json:"address"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	CoinGeckoCoinID string `json:"coingecko_coin_id"`
}
