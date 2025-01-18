package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"pulse/graph/graphql_model"
	"pulse/internal/auth"
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
		Address:   s.Address,
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

	// Create subwallet
	subwallet := &models.Subwallet{
		Name:         input.Name,
		Address:      input.Address,
		WalletID:     input.WalletID,
		ChainID:      input.ChainID,
		CurrentValue: 0,
	}

	// Load token balances
	tokens, err := s.loaders.Solana.LoadWallet(input.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to load wallet: %w", err)
	}

	// Calculate portfolio metrics
	var totalCurrentValue float64
	var portfolioTokens []TokenMetrics

	for _, token := range tokens {
		if token.ConfirmedBalance == "0" {
			continue // Skip tokens with zero balance
		}

		metrics := TokenMetrics{
			Symbol:    token.Currency.Symbol,
			AssetPath: token.Currency.AssetPath,
		}

		// Calculate actual token balance
		balance, err := calculateTokenBalance(token.ConfirmedBalance, token.Currency.Decimals)
		if err != nil {
			log.Printf("warning: failed to calculate balance for %s: %v", token.Currency.Symbol, err)
			continue
		}
		metrics.Balance = balance

		// Get current token price
		currentPrice := 0.0
		if token.Currency.AssetPath == "solana/native/sol" {
			price, err := s.loaders.Solana.LoadCurrentSolanaPrice()
			if err != nil {
				log.Printf("warning: failed to fetch SOL price: %v", err)
				continue
			}
			currentPrice = price
		} else {
			price, _, err := s.loaders.Solana.LoadTokenPrice(token.Currency.AssetPath)
			if err != nil {
				log.Printf("warning: failed to fetch price for %s: %v", token.Currency.Symbol, err)
				continue
			}
			currentPrice = price
		}
		metrics.CurrentPrice = currentPrice

		// Calculate current value
		currentTokenValue := balance * currentPrice
		metrics.CurrentValue = currentTokenValue
		totalCurrentValue += currentTokenValue

		// Load historical transactions
		transactions, err := s.loaders.Solana.LoadTransactions(input.Address, token.Currency.AssetPath)
		if err != nil {
			log.Printf("warning: failed to load transactions for %s: %v", token.Currency.Symbol, err)
			continue
		}

		// Calculate invested value from transactions
		var totalCost float64
		var totalAmount float64

		for _, tx := range transactions {
			// Get historical price at transaction time
			historicalPrice, err := s.loaders.Solana.LoadHistoricalPrice(token.Currency.AssetPath, tx.Date)
			if err != nil {
				log.Printf("warning: failed to fetch historical price: %v", err)
				continue
			}

			switch tx.Type {
			case "RECEIVE":
				totalAmount += tx.Amount
				totalCost += tx.Amount * historicalPrice
				log.Printf("Buy: Amount: %.8f, Price: $%.2f, Cost: $%.2f",
					tx.Amount, historicalPrice, tx.Amount*historicalPrice)
			case "SEND":
				// For sells, remove proportional amount of cost basis
				if totalAmount > 0 {
					costBasisPerToken := totalCost / totalAmount
					removedCost := tx.Amount * costBasisPerToken
					totalCost -= removedCost
					totalAmount -= tx.Amount
					log.Printf("Sell: Amount: %.8f, Cost Basis: $%.2f, Removed Cost: $%.2f",
						tx.Amount, costBasisPerToken, removedCost)
				}
			}
		}

		// Calculate average cost basis if we have any remaining tokens
		var avgCostBasis float64
		if totalAmount > 0 {
			avgCostBasis = totalCost / totalAmount
		}

		metrics.InvestedValue = balance * avgCostBasis
		metrics.PnL = currentTokenValue - metrics.InvestedValue
		if metrics.InvestedValue > 0 {
			metrics.PnLPercentage = (metrics.PnL / metrics.InvestedValue) * 100
		}

		portfolioTokens = append(portfolioTokens, metrics)

		// Log token metrics
		log.Printf("\nToken: %s", metrics.Symbol)
		log.Printf("  Balance: %.8f", metrics.Balance)
		log.Printf("  Current Price: $%.2f", metrics.CurrentPrice)
		log.Printf("  Current Value: $%.2f", metrics.CurrentValue)
		log.Printf("  Avg Cost Basis: $%.2f", avgCostBasis)
		log.Printf("  Invested Value: $%.2f", metrics.InvestedValue)
		log.Printf("  PnL: $%.2f (%.2f%%)", metrics.PnL, metrics.PnLPercentage)
	}

	// Calculate portfolio totals
	var totalInvestedValue float64
	for _, token := range portfolioTokens {
		totalInvestedValue += token.InvestedValue
	}

	var pnlPercentage float64
	if totalInvestedValue > 0 {
		pnlPercentage = ((totalCurrentValue - totalInvestedValue) / totalInvestedValue) * 100
	}

	// Log portfolio summary
	log.Printf("\nPortfolio Summary:")
	log.Printf("Total Current Value: $%.2f", totalCurrentValue)
	log.Printf("Total Invested Value: $%.2f", totalInvestedValue)
	log.Printf("Total PnL: $%.2f (%.2f%%)",
		totalCurrentValue-totalInvestedValue,
		pnlPercentage)

	subwallet.CurrentValue = totalCurrentValue
	// Save subwallet to database
	/* TODO: reenable
	if err := s.db.Create(subwallet).Error; err != nil {
		return nil, fmt.Errorf("failed to create subwallet: %w", err)
	}
	*/
	/*
		// Store token balances and metrics
		for _, metrics := range portfolioTokens {
			tokenBalance := &models.TokenBalance{
				SubwalletID:     subwallet.ID,
				TokenAddress:    metrics.AssetPath,
				Symbol:          metrics.Symbol,
				Balance:         fmt.Sprintf("%.8f", metrics.Balance),
				CurrentPrice:    fmt.Sprintf("%.8f", metrics.CurrentPrice),
				CurrentValue:    fmt.Sprintf("%.2f", metrics.CurrentValue),
				InvestedValue:   fmt.Sprintf("%.2f", metrics.InvestedValue),
				PnL:             fmt.Sprintf("%.2f", metrics.PnL),
				PnLPercentage:   fmt.Sprintf("%.2f", metrics.PnLPercentage),
				LastUpdateBlock: time.Now().Unix(),
			}
			if err := s.db.Create(tokenBalance).Error; err != nil {
				log.Printf("warning: failed to store token balance for %s: %v", metrics.Symbol, err)
			}
		}
	*/
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

func calculateTokenBalance(balance string, decimals int) (float64, error) {
	rawBalance, err := parseFloat(balance)
	if err != nil {
		return 0, fmt.Errorf("error converting balance to float: %v", err)
	}
	return rawBalance / math.Pow10(decimals), nil
}

func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	f, err := json.Number(s).Float64()
	if err != nil {
		return 0, fmt.Errorf("error parsing float: %v", err)
	}
	return f, nil
}
