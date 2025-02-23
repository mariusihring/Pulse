package coingecko_types

type PoolResponse struct {
	Data []Pool `json:"data"`
}

// Pool represents each pool object in the data array.
type Pool struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Attributes    Attributes    `json:"attributes"`
	Relationships Relationships `json:"relationships"`
}

// Attributes contains all the pricing, volume, and transaction details for a pool.
type Attributes struct {
	BaseTokenPriceUSD             string                `json:"base_token_price_usd"`
	BaseTokenPriceNativeCurrency  string                `json:"base_token_price_native_currency"`
	QuoteTokenPriceUSD            string                `json:"quote_token_price_usd"`
	QuoteTokenPriceNativeCurrency string                `json:"quote_token_price_native_currency"`
	BaseTokenPriceQuoteToken      string                `json:"base_token_price_quote_token"`
	QuoteTokenPriceBaseToken      string                `json:"quote_token_price_base_token"`
	Address                       string                `json:"address"`
	Name                          string                `json:"name"`
	PoolCreatedAt                 string                `json:"pool_created_at"`
	TokenPriceUSD                 string                `json:"token_price_usd"`
	FDVUSD                        string                `json:"fdv_usd"`
	MarketCapUSD                  string                `json:"market_cap_usd"`
	PriceChangePercentage         PriceChangePercentage `json:"price_change_percentage"`
	Transactions                  Transactions          `json:"transactions"`
	VolumeUSD                     VolumeUSD             `json:"volume_usd"`
	ReserveInUSD                  string                `json:"reserve_in_usd"`
}

// PriceChangePercentage contains percentage changes over different time frames.
type PriceChangePercentage struct {
	M5  string `json:"m5"`
	H1  string `json:"h1"`
	H6  string `json:"h6"`
	H24 string `json:"h24"`
}

// Transactions holds transaction counts for different time periods.
type Transactions struct {
	M5  TransactionPeriod `json:"m5"`
	M15 TransactionPeriod `json:"m15"`
	M30 TransactionPeriod `json:"m30"`
	H1  TransactionPeriod `json:"h1"`
	H24 TransactionPeriod `json:"h24"`
}

// TransactionPeriod represents the counts of buys, sells, buyers, and sellers.
type TransactionPeriod struct {
	Buys    int `json:"buys"`
	Sells   int `json:"sells"`
	Buyers  int `json:"buyers"`
	Sellers int `json:"sellers"`
}

// VolumeUSD holds volume data over various time frames.
type VolumeUSD struct {
	M5  string `json:"m5"`
	H1  string `json:"h1"`
	H6  string `json:"h6"`
	H24 string `json:"h24"`
}

// Relationships holds related objects for a pool.
type Relationships struct {
	BaseToken  RelationshipItem `json:"base_token"`
	QuoteToken RelationshipItem `json:"quote_token"`
	Dex        RelationshipItem `json:"dex"`
}

// RelationshipItem wraps the relationship data.
type RelationshipItem struct {
	Data Reference `json:"data"`
}

// Reference represents a related entity with an id and type.
type Reference struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
