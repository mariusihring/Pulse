package services

import (
	"pulse/internal/services/loaders"

	"gorm.io/gorm"
)

type WalletService struct {
	db      *gorm.DB
	loaders *Loaders
}

type Loaders struct {
	Solana *loaders.SolanaLoader
}
type TokenMetrics struct {
	Symbol        string
	AssetPath     string
	Balance       float64
	CurrentPrice  float64
	CurrentValue  float64
	InvestedValue float64
	PnL           float64
	PnLPercentage float64
}

func NewWalletService(db *gorm.DB, solanaLoader *loaders.SolanaLoader) *WalletService {
	loaders := Loaders{Solana: solanaLoader}
	return &WalletService{db: db, loaders: &loaders}
}

// func (s *WalletService) CreateWallet(ctx context.Context, input *graphql_model.CreateWalletInput) (*graphql_model.Wallet, error) {
// 	request := &solana_grpc.WalletRequest{WalletAddress: input.Name}
// 	s.loaders.Solana.Client.AddWallet(ctx, request)
// 	// 	return toGQLWallet(wallet), nil
// 	// }

// 	// func (s *WalletService) GetWallet(ctx context.Context, id uuid.UUID) (*graphql_model.Wallet, error) {
// 	// 	var wallet models.Wallet
// 	// 	if err := s.db.Preload("Subwallets").Preload("Subwallets.Chain").First(&wallet, "id = ?", id).Error; err != nil {
// 	// 		return nil, fmt.Errorf("failed to get wallet: %w", err)
// 	// 	}

// 	// 	return toGQLWallet(&wallet), nil
// 	return nil, nil
// }
