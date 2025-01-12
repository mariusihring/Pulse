package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"pulse/graph/graphql_model"
	"pulse/internal/auth"
	"pulse/internal/db/models"
)

type WalletService struct {
	db *gorm.DB
}

func NewWalletService(db *gorm.DB) *WalletService {
	return &WalletService{db: db}
}

func toGQLWallet(w *models.Wallet) *graphql_model.Wallet {
	if w == nil {
		return nil
	}

	// Convert subwallets if they exist
	subwallets := make([]*graphql_model.Subwallet, len(w.Subwallets))
	for i, s := range w.Subwallets {
		subwallets[i] = toGQLSubwallet(&s)
	}

	return &graphql_model.Wallet{
		ID:         w.ID,
		CreatedAt:  w.CreatedAt,
		UpdatedAt:  w.UpdatedAt,
		Name:       w.Name,
		Subwallets: subwallets,
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
		ID:        s.ID,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		Name:      s.Name,
		Chain:     chain,
	}
}

func (s *WalletService) CreateWallet(ctx context.Context, input *graphql_model.CreateWalletInput) (*graphql_model.Wallet, error) {
	userID := auth.UserFromContext(ctx)
	fmt.Println(userID)
	wallet := &models.Wallet{
		Name:   input.Name,
		UserID: *userID,
	}

	if err := s.db.Create(wallet).Error; err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	// Reload the wallet with all relationships
	if err := s.db.Preload("Subwallets").Preload("Subwallets.Chain").First(wallet, "id = ?", wallet.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload wallet: %w", err)
	}

	return toGQLWallet(wallet), nil
}

func (s *WalletService) GetWallet(ctx context.Context, id uuid.UUID) (*graphql_model.Wallet, error) {

	var wallet models.Wallet
	if err := s.db.Preload("Subwallets").Preload("Subwallets.Chain").First(&wallet, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return toGQLWallet(&wallet), nil
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

	// Verify wallet exists
	var wallet models.Wallet
	if err := s.db.First(&wallet, "id = ?", input.WalletID).Error; err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}

	// Verify chain exists
	var chain models.Chain
	if err := s.db.First(&chain, "id = ?", input.ChainID).Error; err != nil {
		return nil, fmt.Errorf("chain not found: %w", err)
	}

	subwallet := &models.Subwallet{
		Name:     input.Name,
		WalletID: input.WalletID,
		ChainID:  input.ChainID,
	}

	if err := s.db.Create(subwallet).Error; err != nil {
		return nil, fmt.Errorf("failed to create subwallet: %w", err)
	}

	// Reload the subwallet with chain
	if err := s.db.Preload("Chain").First(subwallet, "id = ?", subwallet.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload subwallet: %w", err)
	}

	return toGQLSubwallet(subwallet), nil
}

func (s *WalletService) GetSubwallet(ctx context.Context, id uuid.UUID) (*graphql_model.Subwallet, error) {

	var subwallet models.Subwallet
	if err := s.db.Preload("Chain").First(&subwallet, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get subwallet: %w", err)
	}

	return toGQLSubwallet(&subwallet), nil
}

func (s *WalletService) GetSubwallets(ctx context.Context, walletID string) ([]*graphql_model.Subwallet, error) {
	id, err := uuid.Parse(walletID)
	if err != nil {
		return nil, fmt.Errorf("invalid wallet ID: %w", err)
	}

	var subwallets []models.Subwallet
	if err := s.db.Where("wallet_id = ?", id).Preload("Chain").Find(&subwallets).Error; err != nil {
		return nil, fmt.Errorf("failed to get subwallets: %w", err)
	}

	result := make([]*graphql_model.Subwallet, len(subwallets))
	for i, s := range subwallets {
		result[i] = toGQLSubwallet(&s)
	}
	return result, nil
}
