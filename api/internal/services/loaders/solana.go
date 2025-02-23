package loaders

import (
	grpc "pulse/internal/proto/solana_service"
	"pulse/internal/proto/solana_service/generated"

	"gorm.io/gorm"
)

// SolanaLoader uses Solana JSON RPC endpoints for wallet balances/token accounts,
// CoinGecko/DexScreener for price data, and also loads & stores transactions via GORM.
type SolanaLoader struct {
	Client generated.WalletServiceClient
	Db     *gorm.DB
}

// NewSolanaLoader constructs a new loader using the provided config and a *gorm.DB connection.
func NewSolanaLoader(db *gorm.DB) *SolanaLoader {
	client, _, err := grpc.NewSolanaGRPCClient()
	if err != nil {
		return nil
	}
	return &SolanaLoader{
		Client: client,
		Db:     db,
	}
}
