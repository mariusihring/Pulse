package loaders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"pulse/internal/config"
	"pulse/internal/db/models"
	"pulse/internal/services/loaders/types"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DetailedTransaction is a helper type for processing each event from a blockchain transaction.
type DetailedTransaction struct {
	EventID         string
	BlockNumber     int64
	Date            time.Time
	TransactionType string // "SEND", "RECEIVE", or "UNKNOWN"
	Amount          float64
	AssetPath       string // e.g. "solana/native/sol" or "solana/token/<mint>"
}

// SolanaLoader uses Solana JSON RPC endpoints for wallet balances/token accounts,
// CoinGecko/DexScreener for price data, and also loads & stores transactions via GORM.
type SolanaLoader struct {
	apiKeys    ApiKeys
	client     *http.Client
	priceCache map[string]HistoricalPrices
	db         *gorm.DB
}

type ApiKeys struct {
	blockdaemon string
	coingecko   string
}

// PricePoint and HistoricalPrices are used for caching historical price data.
type PricePoint struct {
	Timestamp int64
	Price     float64
}

type HistoricalPrices struct {
	Prices []PricePoint
}

// NewSolanaLoader constructs a new loader using the provided config and a *gorm.DB connection.
func NewSolanaLoader(cfg *config.Config, db *gorm.DB) *SolanaLoader {
	apiKeys := ApiKeys{
		blockdaemon: cfg.ApiKeys.Blockdaemon,
		coingecko:   cfg.ApiKeys.Coingecko,
	}
	client := &http.Client{Timeout: 15 * time.Second}
	return &SolanaLoader{
		apiKeys:    apiKeys,
		client:     client,
		priceCache: make(map[string]HistoricalPrices),
		db:         db,
	}
}

// ----------------------------------------------------------------
// TRANSACTION LOADING (Aggregated and Detailed)
// ----------------------------------------------------------------

// LoadTransactions is your existing aggregated transaction loader (using Blockdaemon).
func (s *SolanaLoader) LoadTransactions(address string, tokenAddress string) ([]types.ProcessedTransaction, error) {
	var allProcessedTxs []types.ProcessedTransaction
	pageToken := "" // Start with empty token for first page

	log.Info("loading transactions", "address", address, "token", tokenAddress)

	for {
		url := fmt.Sprintf(
			"https://svc.blockdaemon.com/universal/v1/solana/mainnet/account/%s/txs?order=desc&page_size=100",
			address)
		if tokenAddress != "" {
			url += "&token=" + tokenAddress
		}
		if pageToken != "" {
			url += "&page_token=" + pageToken
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create GetTransactions request: %w", err)
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.apiKeys.blockdaemon))
		resp, err := s.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to execute GetTransactions request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("GetTransactions request returned status %d: %s", resp.StatusCode, string(bodyBytes))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read GetTransactions response body: %w", err)
		}

		var txResponse types.TransactionListResponse
		if err := json.Unmarshal(body, &txResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transaction response: %w", err)
		}

		// Process each transaction (aggregated)
		for _, tx := range txResponse.Data {
			processed := types.ProcessedTransaction{
				ID:            tx.ID,
				BlockNumber:   tx.BlockNumber,
				Date:          time.Unix(tx.Date, 0),
				Confirmations: tx.Confirmations,
				Type:          "", // Will be set based on analysis
				Amount:        0,  // Will be calculated
			}

			var totalInput, totalOutput float64
			for _, event := range tx.Events {
				amount := float64(event.Amount) / math.Pow10(event.Decimals)
				if event.Source == address {
					totalInput += amount
				}
				if event.Destination == address {
					totalOutput += amount
				}
			}

			netAmount := totalOutput - totalInput
			if netAmount > 0 {
				processed.Type = "RECEIVE"
				processed.Amount = netAmount
			} else if netAmount < 0 {
				processed.Type = "SEND"
				processed.Amount = -netAmount
			} else {
				processed.Type = "UNKNOWN"
				processed.Amount = 0
			}

			log.Debug("processed transaction",
				"id", processed.ID,
				"date", processed.Date,
				"type", processed.Type,
				"amount", processed.Amount,
				"inputs", totalInput,
				"outputs", totalOutput,
				"confirmations", processed.Confirmations)

			if processed.Amount > 0 {
				allProcessedTxs = append(allProcessedTxs, processed)
			}
		}

		if txResponse.Meta.Paging.NextPageToken == "" {
			break
		}
		pageToken = txResponse.Meta.Paging.NextPageToken
	}

	log.Info("finished loading transactions", "address", address, "token", tokenAddress, "count", len(allProcessedTxs))
	return allProcessedTxs, nil
}

// LoadAllTransactions loads detailed transactions from Blockdaemon,
// creating one DetailedTransaction per event.
func (s *SolanaLoader) LoadAllTransactions(address string) ([]DetailedTransaction, error) {
	var detailedTxs []DetailedTransaction
	pageToken := ""
	log.Info("loading all transactions", "address", address)

	for {
		url := fmt.Sprintf("https://svc.blockdaemon.com/universal/v1/solana/mainnet/account/%s/txs?order=desc&page_size=100", address)
		if pageToken != "" {
			url += "&page_token=" + pageToken
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create transactions request: %w", err)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.apiKeys.blockdaemon))
		resp, err := s.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to execute transactions request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("transactions request returned status %d: %s", resp.StatusCode, string(bodyBytes))
		}

		var txResponse types.TransactionListResponse
		if err := json.NewDecoder(resp.Body).Decode(&txResponse); err != nil {
			return nil, fmt.Errorf("failed to decode transactions response: %w", err)
		}

		// For each transaction, iterate over each event.
		for _, tx := range txResponse.Data {
			for _, event := range tx.Events {
				var dt DetailedTransaction
				dt.EventID = event.ID
				dt.BlockNumber = tx.BlockNumber
				dt.Date = time.Unix(event.Date, 0)
				// Determine transaction type based on whether the wallet is source or destination.
				if strings.EqualFold(event.Source, address) {
					dt.TransactionType = "SEND"
				} else if strings.EqualFold(event.Destination, address) {
					dt.TransactionType = "RECEIVE"
				} else {
					dt.TransactionType = "UNKNOWN"
				}
				dt.Amount = float64(event.Amount) / math.Pow10(event.Decimals)
				// If the event's Denomination is empty or "SOL", assume native SOL.
				if strings.EqualFold(event.Denomination, "SOL") || event.Denomination == "" {
					dt.AssetPath = "solana/native/sol"
				} else {
					dt.AssetPath = fmt.Sprintf("solana/token/%s", event.Denomination)
				}
				detailedTxs = append(detailedTxs, dt)
			}
		}

		if txResponse.Meta.Paging.NextPageToken == "" {
			break
		}
		pageToken = txResponse.Meta.Paging.NextPageToken
	}

	log.Info("loaded all transactions", "count", len(detailedTxs))
	return detailedTxs, nil
}

// ----------------------------------------------------------------
// DATABASE HELPERS USING GORM (using models from "myproject/models")
// ----------------------------------------------------------------

// lookupTokenID returns the token ID (as a UUID string) from the tokens table
// by matching on the token name. For native SOL we assume the name is "Solana";
// for tokens we assume the token name is the mint address.
func (s *SolanaLoader) lookupTokenID(assetPath string) (*uuid.UUID, error) {
	var tokenName string
	if assetPath == "solana/native/sol" {
		tokenName = "Solana"
	} else {
		parts := strings.Split(assetPath, "/")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid assetPath: %s", assetPath)
		}
		tokenName = parts[2]
	}

	var token models.Token
	err := s.db.Where("name = ?", tokenName).First(&token).Error
	if err != nil {
		return nil, fmt.Errorf("lookupTokenID error: %w", err)
	}
	return &token.ID, nil
}

// lookupCategoryID returns the transaction category ID from the transaction_categories table
// based on the transaction type (e.g. "SEND", "RECEIVE").
func (s *SolanaLoader) lookupCategoryID(transactionType string) (*uuid.UUID, error) {
	var category models.TransactionCategory
	err := s.db.Where("name = ?", transactionType).First(&category).Error
	if err != nil {
		return nil, fmt.Errorf("lookupCategoryID error: %w", err)
	}
	return &category.ID, nil
}

// storeDetailedTransaction inserts a single DetailedTransaction into the transactions table.
func (s *SolanaLoader) storeDetailedTransaction(dt DetailedTransaction) error {
	tokenID, err := s.lookupTokenID(dt.AssetPath)
	if err != nil {
		return err
	}
	categoryID, err := s.lookupCategoryID(dt.TransactionType)
	if err != nil {
		return err
	}
	txModel := models.Transaction{
		TokenID:         *tokenID,
		TransactionType: dt.TransactionType,
		Amount:          dt.Amount,
		TransactionDate: dt.Date,
		CategoryID:      *categoryID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := s.db.Create(&txModel).Error; err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}
	return nil
}

// StoreTransactions loads all detailed transactions for the given address and stores them in the database.
func (s *SolanaLoader) StoreTransactions(address string) error {
	detailedTxs, err := s.LoadAllTransactions(address)
	if err != nil {
		return fmt.Errorf("error loading detailed transactions: %w", err)
	}

	for _, dt := range detailedTxs {
		if err := s.storeDetailedTransaction(dt); err != nil {
			log.Error("failed to store transaction", "error", err, "transaction", dt)
			// Optionally, you can choose to continue processing or return the error.
		}
	}
	return nil
}

// ----------------------------------------------------------------
// WALLET & PRICE FUNCTIONS (Unchanged from previous implementation)
// ----------------------------------------------------------------

// LoadWallet retrieves the SOL balance and SPL token accounts for a given address,
// and converts them into types.TokenBalance.
func (s *SolanaLoader) LoadWallet(address string) ([]types.TokenBalance, error) {
	solBalance, solSlot, err := s.getBalance(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get SOL balance: %w", err)
	}

	tokenBalances, err := s.getTokenAccounts(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get token accounts: %w", err)
	}

	var balances []types.TokenBalance

	solToken := types.TokenBalance{
		Currency: types.Currency{
			AssetPath: "solana/native/sol",
			Symbol:    "SOL",
			Name:      "Solana",
			Decimals:  9,
			Type:      "native",
		},
		ConfirmedBalance: fmt.Sprintf("%f", solBalance),
		ConfirmedBlock:   solSlot,
	}
	balances = append(balances, solToken)
	balances = append(balances, tokenBalances...)

	return balances, nil
}

func (s *SolanaLoader) getBalance(address string) (float64, int64, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBalance",
		"params":  []interface{}{address},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to marshal getBalance payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.mainnet-beta.solana.com", bytes.NewReader(data))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create getBalance request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to perform getBalance request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return 0, 0, fmt.Errorf("getBalance request returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var rpcResp struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			Context struct {
				Slot int64 `json:"slot"`
			} `json:"context"`
			Value uint64 `json:"value"`
		} `json:"result"`
		Id int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return 0, 0, fmt.Errorf("failed to decode getBalance response: %w", err)
	}

	solBalance := float64(rpcResp.Result.Value) / 1e9
	return solBalance, rpcResp.Result.Context.Slot, nil
}

func (s *SolanaLoader) getTokenAccounts(address string) ([]types.TokenBalance, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTokenAccountsByOwner",
		"params": []interface{}{
			address,
			map[string]interface{}{
				"programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
			},
			map[string]interface{}{
				"encoding": "jsonParsed",
			},
		},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal getTokenAccounts payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.mainnet-beta.solana.com", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create getTokenAccounts request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform getTokenAccounts request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("getTokenAccounts request returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenResp struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			Context struct {
				Slot int64 `json:"slot"`
			} `json:"context"`
			Value []struct {
				Pubkey  string `json:"pubkey"`
				Account struct {
					Data struct {
						Parsed struct {
							Info struct {
								Mint        string `json:"mint"`
								Owner       string `json:"owner"`
								TokenAmount struct {
									Amount         string  `json:"amount"`
									Decimals       int     `json:"decimals"`
									UiAmount       float64 `json:"uiAmount"`
									UiAmountString string  `json:"uiAmountString"`
								} `json:"tokenAmount"`
							} `json:"info"`
							Type string `json:"type"`
						} `json:"parsed"`
						Program string `json:"program"`
						Space   int    `json:"space"`
					} `json:"data"`
					Executable bool   `json:"executable"`
					Lamports   int64  `json:"lamports"`
					Owner      string `json:"owner"`
					RentEpoch  uint64 `json:"rentEpoch"`
				} `json:"account"`
			} `json:"value"`
		} `json:"result"`
		Id int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode getTokenAccounts response: %w", err)
	}

	var tokens []types.TokenBalance
	for _, account := range tokenResp.Result.Value {
		info := account.Account.Data.Parsed.Info
		token := types.TokenBalance{
			Currency: types.Currency{
				AssetPath: fmt.Sprintf("solana/token/%s", info.Mint),
				Symbol:    info.Mint,
				Name:      info.Mint,
				Decimals:  info.TokenAmount.Decimals,
				Type:      "token",
			},
			ConfirmedBalance: fmt.Sprintf("%f", info.TokenAmount.UiAmount),
			ConfirmedBlock:   0, // not available from this endpoint
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (s *SolanaLoader) LoadCurrentSolanaPrice() (float64, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=solana&vs_currencies=usd"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("x-cg-demo-api-key", s.apiKeys.coingecko)
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to execute Coingecko data request: %w", err)
	}
	defer resp.Body.Close()

	var result types.PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return result.Solana.USD, nil
}

func (s *SolanaLoader) LoadTokenPrice(address string) (float64, *types.DexScreenerResponse, error) {
	if address == "solana/native/sol" {
		price, err := s.LoadCurrentSolanaPrice()
		return price, nil, err
	}

	parts := strings.Split(address, "/")
	if len(parts) != 3 {
		return 0, nil, fmt.Errorf("invalid token address format: %s", address)
	}
	tokenContractAddress := parts[2]
	lowerMint := strings.ToLower(tokenContractAddress)

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/token_price/solana?contract_addresses=%s&vs_currencies=usd", lowerMint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create request for CoinGecko token price: %w", err)
	}
	req.Header.Add("x-cg-demo-api-key", s.apiKeys.coingecko)
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to execute CoinGecko token price request: %w", err)
	}
	defer resp.Body.Close()

	var cgResp map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&cgResp); err != nil {
		return 0, nil, fmt.Errorf("failed to decode CoinGecko token price response: %w", err)
	}
	if data, exists := cgResp[lowerMint]; exists {
		if price, exists := data["usd"]; exists && price > 0 {
			return price, nil, nil
		}
	}

	return s.getTokenPriceDexscreener(tokenContractAddress)
}

func (s *SolanaLoader) getTokenPriceDexscreener(mint string) (float64, *types.DexScreenerResponse, error) {
	url := fmt.Sprintf("https://api.dexscreener.com/latest/dex/tokens/%s", strings.ToLower(mint))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create request for DexScreener: %w", err)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to execute DexScreener request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return 0, nil, fmt.Errorf("DexScreener returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var dsResp types.DexScreenerResponse
	if err := json.NewDecoder(resp.Body).Decode(&dsResp); err != nil {
		return 0, nil, fmt.Errorf("failed to decode DexScreener response: %w", err)
	}

	if len(dsResp.Pairs) > 0 {
		price, err := parseFloat(dsResp.Pairs[0].PriceUSD)
		return price, &dsResp, err
	}
	return 0, &dsResp, nil
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

func (s *SolanaLoader) LoadHistoricalSolanaPrice(timestamp time.Time) (float64, error) {
	cacheKey := timestamp.Format("2006-01-02")
	if prices, exists := s.priceCache[cacheKey]; exists {
		return findClosestPrice(prices.Prices, timestamp.UnixMilli())
	}

	startOfDay := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, timestamp.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/solana/market_chart/range?vs_currency=usd&from=%d&to=%d", startOfDay.Unix(), endOfDay.Unix())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("x-cg-demo-api-key", s.apiKeys.coingecko)
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to execute Coingecko request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("CoinGecko returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Prices [][2]float64 `json:"prices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	prices := HistoricalPrices{
		Prices: make([]PricePoint, len(result.Prices)),
	}
	for i, p := range result.Prices {
		prices.Prices[i] = PricePoint{
			Timestamp: int64(p[0]),
			Price:     p[1],
		}
	}
	s.priceCache[cacheKey] = prices

	return findClosestPrice(prices.Prices, timestamp.UnixMilli())
}

func (s *SolanaLoader) LoadHistoricalTokenPrice(assetPath string, timestamp time.Time) (float64, error) {
	cacheKey := fmt.Sprintf("%s_%s", assetPath, timestamp.Format("2006-01-02"))
	if prices, exists := s.priceCache[cacheKey]; exists {
		return findClosestPrice(prices.Prices, timestamp.Unix())
	}

	var tokenAddress string
	if parts := strings.Split(assetPath, "/"); len(parts) == 3 {
		tokenAddress = parts[2]
	} else {
		return 0, fmt.Errorf("invalid token address format: %s", assetPath)
	}

	dateStr := timestamp.Format("2006-01-02")
	url := fmt.Sprintf("https://api.dexscreener.com/latest/dex/tokens/%s/candles?from=%s&to=%s&resolution=1h", tokenAddress, dateStr, dateStr)
	resp, err := s.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to execute DexScreener request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("DexScreener returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Candles []struct {
			Timestamp int64   `json:"timestamp"`
			Close     float64 `json:"close"`
		} `json:"candles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	prices := HistoricalPrices{
		Prices: make([]PricePoint, len(result.Candles)),
	}
	for i, candle := range result.Candles {
		prices.Prices[i] = PricePoint{
			Timestamp: candle.Timestamp,
			Price:     candle.Close,
		}
	}
	s.priceCache[cacheKey] = prices

	return findClosestPrice(prices.Prices, timestamp.Unix())
}

func findClosestPrice(prices []PricePoint, targetTs int64) (float64, error) {
	if len(prices) == 0 {
		return 0, fmt.Errorf("no price data available")
	}

	closestPrice := prices[0].Price
	smallestDiff := math.Abs(float64(prices[0].Timestamp - targetTs))

	for _, pp := range prices {
		diff := math.Abs(float64(pp.Timestamp - targetTs))
		if diff < smallestDiff {
			smallestDiff = diff
			closestPrice = pp.Price
		}
	}

	return closestPrice, nil
}

func (s *SolanaLoader) LoadHistoricalSolanaPrices(startDate, endDate time.Time) (map[string]float64, error) {
	prices := make(map[string]float64)

	startDay := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/solana/market_chart/range?vs_currency=usd&from=%d&to=%d", startDay.Unix(), endDay.Unix())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("x-cg-demo-api-key", s.apiKeys.coingecko)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute Coingecko request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("CoinGecko returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Prices [][2]float64 `json:"prices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	for _, p := range result.Prices {
		timestamp := time.Unix(int64(p[0]/1000), 0)
		dateKey := timestamp.Format("2006-01-02")
		prices[dateKey] = p[1]
	}

	return prices, nil
}

func (s *SolanaLoader) LoadHistoricalTokenPrices(assetPath string, startDate, endDate time.Time) (map[string]float64, error) {
	prices := make(map[string]float64)

	var tokenAddress string
	if parts := strings.Split(assetPath, "/"); len(parts) == 3 {
		tokenAddress = parts[2]
	} else {
		return nil, fmt.Errorf("invalid token address format: %s", assetPath)
	}

	startDay := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	url := fmt.Sprintf("https://api.dexscreener.com/latest/dex/tokens/%s/candles?from=%s&to=%s&resolution=1d",
		tokenAddress,
		startDay.Format("2006-01-02"),
		endDay.Format("2006-01-02"),
	)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to execute DexScreener request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("DexScreener returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Candles []struct {
			Timestamp int64   `json:"timestamp"`
			Close     float64 `json:"close"`
		} `json:"candles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	for _, candle := range result.Candles {
		timestamp := time.Unix(candle.Timestamp, 0)
		dateKey := timestamp.Format("2006-01-02")
		prices[dateKey] = candle.Close
	}

	return prices, nil
}
