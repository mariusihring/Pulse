package types

import "time"

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
type TransactionListResponse struct {
	Total int64               `json:"total"`
	Data  []TransactionDetail `json:"data"`
	Meta  TransactionListMeta `json:"meta"`
}

type TransactionListMeta struct {
	Paging PagingInfo `json:"paging"`
}

type PagingInfo struct {
	NextPageToken string `json:"next_page_token"`
}

// TransactionDetail represents a single transaction
type TransactionDetail struct {
	ID            string          `json:"id"`
	BlockID       string          `json:"block_id"`
	Date          int64           `json:"date"`
	Status        string          `json:"status"`
	NumEvents     int             `json:"num_events"`
	Meta          TransactionMeta `json:"meta"`
	BlockNumber   int64           `json:"block_number"`
	Confirmations int64           `json:"confirmations"`
	Events        []Event         `json:"events"`
}

type TransactionMeta struct {
	Index int64 `json:"index"`
	VSize int64 `json:"vsize"`
}

// Event represents a transaction event (input, output, or fee)
type Event struct {
	ID            string     `json:"id"`
	TransactionID string     `json:"transaction_id"`
	Type          string     `json:"type"`
	Denomination  string     `json:"denomination"`
	Source        string     `json:"source,omitempty"`
	Destination   string     `json:"destination,omitempty"`
	Meta          *EventMeta `json:"meta"`
	Date          int64      `json:"date"`
	Amount        int64      `json:"amount"`
	Decimals      int        `json:"decimals"`
}

type EventMeta struct {
	Addresses  []string `json:"addresses,omitempty"`
	Index      int64    `json:"index,omitempty"`
	Script     string   `json:"script,omitempty"`
	ScriptType string   `json:"script_type,omitempty"`
}

// Helper types for processing transactions
type ProcessedTransaction struct {
	ID            string
	BlockNumber   int64
	Date          time.Time
	Confirmations int64
	Type          string  // "SEND" or "RECEIVE"
	Amount        float64 // Net amount for the address
}

type ProcessedEvent struct {
	Address string
	Amount  float64
}
type DexScreenerResponse struct {
	SchemaVersion string    `json:"schemaVersion"`
	Pairs         []DexPair `json:"pairs"`
}

// DexPair represents each object in the "pairs" array.
type DexPair struct {
	ChainID       string         `json:"chainId"`
	DexID         string         `json:"dexId"`
	URL           string         `json:"url"`
	PairAddress   string         `json:"pairAddress"`
	Labels        []string       `json:"labels,omitempty"`
	BaseToken     DexToken       `json:"baseToken"`
	QuoteToken    DexToken       `json:"quoteToken"`
	PriceNative   string         `json:"priceNative"`
	PriceUSD      string         `json:"priceUsd"`
	Txns          DexTxns        `json:"txns"`
	Volume        DexVolume      `json:"volume"`
	PriceChange   DexPriceChange `json:"priceChange"`
	Liquidity     DexLiquidity   `json:"liquidity"`
	FDV           float64        `json:"fdv"`           // Changed from int64 to float64
	MarketCap     float64        `json:"marketCap"`     // Also changed this to be consistent
	PairCreatedAt int64          `json:"pairCreatedAt"` // This can stay as int64 since it's a timestamp
	Info          DexInfo        `json:"info"`
}

// DexToken holds fields for baseToken and quoteToken.
type DexToken struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
}

// DexTxns represents the "txns" object with subfields (m5, h1, h6, h24).
type DexTxns struct {
	M5  DexTxnStats `json:"m5"`
	H1  DexTxnStats `json:"h1"`
	H6  DexTxnStats `json:"h6"`
	H24 DexTxnStats `json:"h24"`
}

// DexTxnStats holds the "buys" and "sells" counts.
type DexTxnStats struct {
	Buys  int64 `json:"buys"`
	Sells int64 `json:"sells"`
}

// DexVolume represents volume metrics (h24, h6, h1, m5).
type DexVolume struct {
	H24 float64 `json:"h24"`
	H6  float64 `json:"h6"`
	H1  float64 `json:"h1"`
	M5  float64 `json:"m5"`
}

// DexPriceChange represents price change metrics (m5, h1, h6, h24).
type DexPriceChange struct {
	M5  float64 `json:"m5,omitempty"`
	H1  float64 `json:"h1,omitempty"`
	H6  float64 `json:"h6,omitempty"`
	H24 float64 `json:"h24,omitempty"`
}

// DexLiquidity represents the liquidity object.
type DexLiquidity struct {
	USD   float64 `json:"usd"`
	Base  float64 `json:"base"`
	Quote float64 `json:"quote"`
}

// DexInfo represents the "info" object containing image URLs, websites, socials, etc.
type DexInfo struct {
	ImageURL  string       `json:"imageUrl"`
	Header    string       `json:"header"`
	OpenGraph string       `json:"openGraph"`
	Websites  []DexWebsite `json:"websites"`
	Socials   []DexSocial  `json:"socials"`
}

// DexWebsite represents each object in "info.websites".
type DexWebsite struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// DexSocial represents each object in "info.socials".
type DexSocial struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}
