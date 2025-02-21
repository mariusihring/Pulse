package loaders

import (
	"net/http"
	"time"

	"gorm.io/gorm"
)

// SolanaLoader uses Solana JSON RPC endpoints for wallet balances/token accounts,
// CoinGecko/DexScreener for price data, and also loads & stores transactions via GORM.
type SolanaLoader struct {
	client *http.Client
	db     *gorm.DB
}

// NewSolanaLoader constructs a new loader using the provided config and a *gorm.DB connection.
func NewSolanaLoader(db *gorm.DB) *SolanaLoader {
	client := &http.Client{Timeout: 15 * time.Second}
	return &SolanaLoader{
		client: client,
		db:     db,
	}
}
