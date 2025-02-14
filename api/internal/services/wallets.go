package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"pulse/graph/graphql_model"
	"pulse/internal/auth"
	"pulse/internal/db/models"
	"pulse/internal/services/loaders"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	log.Info("creating subwallet",
		"wallet_id", input.WalletID,
		"chain_id", input.ChainID,
		"address", input.Address)

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
		log.Error("failed to load wallet",
			"address", input.Address,
			"error", err)
		return nil, fmt.Errorf("failed to load wallet: %w", err)
	}

	var totalCurrentValue float64
	var portfolioTokens []TokenMetrics

	for _, token := range tokens {
		if token.ConfirmedBalance == "0" {
			continue
		}

		logger := log.With(
			"symbol", token.Currency.Symbol,
			"asset_path", token.Currency.AssetPath)

		metrics := TokenMetrics{
			Symbol:    token.Currency.Symbol,
			AssetPath: token.Currency.AssetPath,
		}

		// Calculate actual token balance
		balance, err := calculateTokenBalance(token.ConfirmedBalance, token.Currency.Decimals)
		if err != nil {
			logger.Warn("failed to calculate balance", "error", err)
			continue
		}
		metrics.Balance = balance

		// Get current price first
		currentPrice := 0.0
		if token.Currency.AssetPath == "solana/native/sol" {
			currentPrice, err = s.loaders.Solana.LoadCurrentSolanaPrice()
		} else {
			currentPrice, _, err = s.loaders.Solana.LoadTokenPrice(token.Currency.AssetPath)
		}
		if err != nil {
			logger.Warn("failed to fetch current price", "error", err)
			continue
		}

		metrics.CurrentPrice = currentPrice
		metrics.CurrentValue = balance * currentPrice
		totalCurrentValue += metrics.CurrentValue

		// Now load historical data for cost basis calculation
		transactions, err := s.loaders.Solana.LoadTransactions(input.Address, token.Currency.AssetPath)
		if err != nil {
			logger.Warn("failed to load transactions", "error", err)
			continue
		}

		if len(transactions) == 0 {
			logger.Debug("no transactions found")
			continue
		}

		// Find transaction date range
		firstTx := transactions[len(transactions)-1]
		lastTx := transactions[0]

		logger.Debug("loading historical prices",
			"from", firstTx.Date,
			"to", lastTx.Date)

		// Load historical prices
		var historicalPrices map[string]float64
		if token.Currency.AssetPath == "solana/native/sol" {
			historicalPrices, err = s.loaders.Solana.LoadHistoricalSolanaPrices(firstTx.Date, lastTx.Date)
		} else {
			historicalPrices, err = s.loaders.Solana.LoadHistoricalTokenPrices(token.Currency.AssetPath, firstTx.Date, lastTx.Date)
		}
		if err != nil {
			logger.Warn("failed to fetch historical prices", "error", err)
			continue
		}

		// Calculate invested value
		var totalCost float64
		var totalAmount float64

		for _, tx := range transactions {
			dateKey := tx.Date.Format("2006-01-02")
			historicalPrice, exists := historicalPrices[dateKey]
			if !exists {
				logger.Warn("no historical price found", "date", dateKey)
				continue
			}

			switch tx.Type {
			case "RECEIVE":
				totalAmount += tx.Amount
				totalCost += tx.Amount * historicalPrice
				logger.Debug("buy transaction",
					"amount", tx.Amount,
					"price", historicalPrice,
					"cost", tx.Amount*historicalPrice)
			case "SEND":
				if totalAmount > 0 {
					costBasisPerToken := totalCost / totalAmount
					removedCost := tx.Amount * costBasisPerToken
					totalCost -= removedCost
					totalAmount -= tx.Amount
					logger.Debug("sell transaction",
						"amount", tx.Amount,
						"cost_basis", costBasisPerToken,
						"removed_cost", removedCost)
				}
			}
		}

		// Calculate metrics
		var avgCostBasis float64
		if totalAmount > 0 {
			avgCostBasis = totalCost / totalAmount
		}

		metrics.InvestedValue = balance * avgCostBasis
		metrics.PnL = metrics.CurrentValue - metrics.InvestedValue
		if metrics.InvestedValue > 0 {
			metrics.PnLPercentage = (metrics.PnL / metrics.InvestedValue) * 100
		}

		portfolioTokens = append(portfolioTokens, metrics)

		logger.Info("token metrics calculated",
			"balance", metrics.Balance,
			"current_price", metrics.CurrentPrice,
			"current_value", metrics.CurrentValue,
			"avg_cost_basis", avgCostBasis,
			"invested_value", metrics.InvestedValue,
			"pnl", metrics.PnL,
			"pnl_percentage", metrics.PnLPercentage)
	}

	// Calculate and log portfolio totals
	var totalInvestedValue float64
	for _, token := range portfolioTokens {
		totalInvestedValue += token.InvestedValue
		totalCurrentValue += token.CurrentValue
	}

	var pnlPercentage float64
	if totalInvestedValue > 0 {
		pnlPercentage = ((totalCurrentValue - totalInvestedValue) / totalInvestedValue) * 100
	}

	log.Info("portfolio summary",
		"current_value", totalCurrentValue,
		"invested_value", totalInvestedValue,
		"pnl", totalCurrentValue-totalInvestedValue,
		"pnl_percentage", pnlPercentage)

	subwallet.CurrentValue = totalCurrentValue
	// Save subwallet to database
	if err := s.db.Create(subwallet).Error; err != nil {
		return nil, fmt.Errorf("failed to create subwallet: %w", err)
	}
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
	if err := s.db.Preload(clause.Associations).First(subwallet, "id = ?", subwallet.ID).Error; err != nil {
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
