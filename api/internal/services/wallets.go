package services

import (
	"context"
	"fmt"
	"pulse/graph/graphql_model"
	"pulse/internal/db/models"
	"pulse/internal/services/loaders"

	"github.com/google/uuid"
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

func toGQLWallet(w *models.Wallet) *graphql_model.Wallet {
	if w == nil {
		return nil
	}

	// Convert subwallets if they exist
	subwallets := make([]*graphql_model.Subwallet, len(w.Subwallets))
	totalBalance := float64(0)
	for i, s := range w.Subwallets {
		totalBalance += s.CurrentValue
		subwallets[i] = toGQLSubwallet(&s)
	}

	return &graphql_model.Wallet{
		ID:           w.ID,
		CreatedAt:    w.CreatedAt,
		UpdatedAt:    w.UpdatedAt,
		Name:         w.Name,
		Subwallets:   subwallets,
		TotalBalance: totalBalance,
	}
}

func toGQLSubwallet(s *models.Subwallet) *graphql_model.Subwallet {
	if s == nil {
		return nil
	}

	// Convert chain if it exists
	var chain *graphql_model.Chain
	if s.Chain.ID != uuid.Nil {
		chain = &graphql_model.Chain{
			ID:        s.Chain.ID,
			CreatedAt: s.Chain.CreatedAt,
			UpdatedAt: s.Chain.UpdatedAt,
			Name:      s.Chain.Name,
		}
	}

	return &graphql_model.Subwallet{
		ID:           s.ID,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
		Name:         s.Name,
		Address:      s.Address,
		Chain:        chain,
		CurrentValue: s.CurrentValue,
	}
}

func (s *WalletService) CreateWallet(ctx context.Context, input *graphql_model.CreateWalletInput) (*graphql_model.Wallet, error) {
	// _ := auth.UserFromContext(ctx)
	// 	wallet := &models.Wallet{
	// 		Name:   input.Name,
	// 		UserID: *userID,
	// 	}

	// 	if err := s.db.Create(wallet).Error; err != nil {
	// 		return nil, fmt.Errorf("failed to create wallet: %w", err)
	// 	}

	// 	// Reload the wallet with all relationships
	// 	if err := s.db.Preload("Subwallets").Preload("Subwallets.Chain").First(wallet, "id = ?", wallet.ID).Error; err != nil {
	// 		return nil, fmt.Errorf("failed to reload wallet: %w", err)
	// 	}

	// 	return toGQLWallet(wallet), nil
	// }

	// func (s *WalletService) GetWallet(ctx context.Context, id uuid.UUID) (*graphql_model.Wallet, error) {
	// 	var wallet models.Wallet
	// 	if err := s.db.Preload("Subwallets").Preload("Subwallets.Chain").First(&wallet, "id = ?", id).Error; err != nil {
	// 		return nil, fmt.Errorf("failed to get wallet: %w", err)
	// 	}

	// 	return toGQLWallet(&wallet), nil
	return nil, nil
}

func (s *WalletService) GetWallets(ctx context.Context) ([]*graphql_model.Wallet, error) {
	var wallets []models.Wallet
	if err := s.db.Preload("Subwallets").Preload("Subwallets.Chain").Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}

	result := make([]*graphql_model.Wallet, len(wallets))
	for i, w := range wallets {
		result[i] = toGQLWallet(&w)
	}
	return result, nil
}

func (s *WalletService) CreateSubwallet(ctx context.Context, input *graphql_model.CreateSubwalletInput) (*graphql_model.Subwallet, error) {

	return nil, nil
}

func (s *WalletService) GetSubwallet(ctx context.Context, id uuid.UUID) (*graphql_model.Subwallet, error) {
	var subwallet models.Subwallet
	return toGQLSubwallet(&subwallet), nil
}

func (s *WalletService) GetSubwallets(ctx context.Context, walletID string) ([]*graphql_model.Subwallet, error) {
	var subwallets []models.Subwallet

	result := make([]*graphql_model.Subwallet, len(subwallets))
	for i, s := range subwallets {
		result[i] = toGQLSubwallet(&s)
	}
	return result, nil
}
