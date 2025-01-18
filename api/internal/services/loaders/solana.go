package loaders

import (
	"encoding/json"
	"fmt"
	"io"

	"math"
	"net/http"
	"pulse/internal/config"
	"pulse/internal/services/loaders/types"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type SolanaLoader struct {
	apiKeys    ApiKeys
	client     *http.Client
	priceCache map[string]HistoricalPrices
}

type ApiKeys struct {
	blockdaemon string
	coingecko   string
}

// Add new types for historical price data
type PricePoint struct {
	Timestamp int64
	Price     float64
}

type HistoricalPrices struct {
	Prices []PricePoint
}

func NewSolanaLoader(cfg *config.Config) *SolanaLoader {
	apiKeys := ApiKeys{blockdaemon: cfg.ApiKeys.Blockdaemon, coingecko: cfg.ApiKeys.Coingecko}
	client := &http.Client{Timeout: 15 * time.Second}
	return &SolanaLoader{
		apiKeys:    apiKeys,
		client:     client,
		priceCache: make(map[string]HistoricalPrices),
	}
}

// LoadTransactions retrieves transactions for a given address

func (s *SolanaLoader) LoadTransactions(address string, tokenAddress string) ([]types.ProcessedTransaction, error) {
	var allProcessedTxs []types.ProcessedTransaction
	pageToken := "" // Start with empty token for first page

	log.Info("loading transactions",
		"address", address,
		"token", tokenAddress)

	for {
		// Build URL with page token if present
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
			return nil, fmt.Errorf(
				"GetTransactions request returned status %d: %s",
				resp.StatusCode,
				string(bodyBytes),
			)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read GetTransactions response body: %w", err)
		}

		var txResponse types.TransactionListResponse
		if err := json.Unmarshal(body, &txResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transaction response: %w", err)
		}

		// Process each transaction
		for _, tx := range txResponse.Data {
			processed := types.ProcessedTransaction{
				ID:            tx.ID,
				BlockNumber:   tx.BlockNumber,
				Date:          time.Unix(tx.Date, 0),
				Confirmations: tx.Confirmations,
				Type:          "", // Will be set based on analysis
				Amount:        0,  // Will be calculated
			}

			// Track total input and output for this address
			var totalInput, totalOutput float64

			for _, event := range tx.Events {
				// Convert amount to proper decimal value
				amount := float64(event.Amount) / math.Pow10(event.Decimals)

				if event.Source == address {
					totalInput += amount
				}
				if event.Destination == address {
					totalOutput += amount
				}
			}

			// Calculate net amount and determine transaction type
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

		// Check if there are more pages
		if txResponse.Meta.Paging.NextPageToken == "" {
			break
		}

		pageToken = txResponse.Meta.Paging.NextPageToken
	}

	log.Info("finished loading transactions",
		"address", address,
		"token", tokenAddress,
		"count", len(allProcessedTxs))

	return allProcessedTxs, nil
}

func (s *SolanaLoader) LoadWallet(address string) ([]types.TokenBalance, error) {
	url := fmt.Sprintf("https://svc.blockdaemon.com/universal/v1/solana/mainnet/account/%s", address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.apiKeys.blockdaemon))
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetWalletBalances request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"GetWalletBalances request returned status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GetWalletBalances response body: %w", err)
	}

	log.Printf("GetWalletBalances response body: %s", string(body))

	var balances []types.TokenBalance
	if err := json.Unmarshal(body, &balances); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GetWalletBalances JSON: %w", err)
	}

	return balances, nil
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

	var tokenContractAddress string
	if parts := strings.Split(address, "/"); len(parts) == 3 {
		tokenContractAddress = parts[2]
	} else {
		return 0, nil, fmt.Errorf("invalid token address format: %s", address)
	}

	url := fmt.Sprintf("https://api.dexscreener.com/latest/dex/tokens/%s", tokenContractAddress)
	resp, err := s.client.Get(url)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return 0, nil, fmt.Errorf(
			"DexScreener returned status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var dexRes types.DexScreenerResponse
	if err := json.NewDecoder(resp.Body).Decode(&dexRes); err != nil {
		return 0, nil, fmt.Errorf("failed to decode DexScreener response: %w", err)
	}

	if len(dexRes.Pairs) > 0 {
		price, err := parseFloat(dexRes.Pairs[0].PriceUSD)
		return price, &dexRes, err
	}
	return 0, &dexRes, nil
}

func (s *SolanaLoader) LoadHistoricalSolanaPrice(timestamp time.Time) (float64, error) {
	// Cache key could be based on the day to avoid too many API calls
	cacheKey := timestamp.Format("2006-01-02")

	// Check if we already have the data in memory
	if prices, exists := s.priceCache[cacheKey]; exists {
		return findClosestPrice(prices.Prices, timestamp.UnixMilli())
	}

	// If not in cache, load the full day's data
	startOfDay := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, timestamp.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/coins/solana/market_chart/range?vs_currency=usd&from=%d&to=%d",
		startOfDay.Unix(),
		endOfDay.Unix(),
	)

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

	// Convert to our format and cache
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
	url := fmt.Sprintf("https://api.dexscreener.com/latest/dex/tokens/%s/candles?from=%s&to=%s&resolution=1h",
		tokenAddress,
		dateStr,
		dateStr,
	)

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

	// Convert to our format and cache
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

// Helper function to find the closest price
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

func calculateUsdPrice(rawBalance float64, decimals int) (float64, error) {
	return rawBalance / math.Pow10(decimals), nil
}

// Helper function to parse float values from string
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

func (s *SolanaLoader) LoadHistoricalSolanaPrices(startDate, endDate time.Time) (map[string]float64, error) {
	prices := make(map[string]float64)

	// Round dates to start/end of days
	startDay := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/coins/solana/market_chart/range?vs_currency=usd&from=%d&to=%d",
		startDay.Unix(),
		endDay.Unix(),
	)

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

	// Group prices by day
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

	// Round dates to start/end of days
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

	// Group prices by day
	for _, candle := range result.Candles {
		timestamp := time.Unix(candle.Timestamp, 0)
		dateKey := timestamp.Format("2006-01-02")
		prices[dateKey] = candle.Close
	}

	return prices, nil
}
