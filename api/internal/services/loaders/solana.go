package loaders

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"pulse/internal/config"
	"pulse/internal/services/loaders/types"
	"strings"
	"time"
)

type SolanaLoader struct {
	apiKeys ApiKeys
	client  *http.Client
}

type ApiKeys struct {
	blockdaemon string
	coingecko   string
}

func NewSolanaLoader(cfg *config.Config) *SolanaLoader {
	apiKeys := ApiKeys{blockdaemon: cfg.ApiKeys.Blockdaemon, coingecko: cfg.ApiKeys.Coingecko}
	client := &http.Client{Timeout: 15 * time.Second}
	return &SolanaLoader{apiKeys, client}
}

// LoadTransactions retrieves transactions for a given address

func (s *SolanaLoader) LoadTransactions(address string, tokenAddress string) ([]types.ProcessedTransaction, error) {
	url := fmt.Sprintf(
		"https://svc.blockdaemon.com/universal/v1/solana/mainnet/account/%s/txs?order=desc&page_size=100",
		address)

	if tokenAddress != "" {
		url += "&token=" + tokenAddress
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
	var processedTxs []types.ProcessedTransaction
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
			processed.Amount = -netAmount // Make positive for consistency
		} else {
			processed.Type = "UNKNOWN"
			processed.Amount = 0
		}

		log.Printf("Transaction %s:", processed.ID)
		log.Printf("  Date: %s", processed.Date)
		log.Printf("  Type: %s", processed.Type)
		log.Printf("  Net Amount: %.8f", processed.Amount)
		log.Printf("  Inputs (from address): %.8f", totalInput)
		log.Printf("  Outputs (to address): %.8f", totalOutput)
		log.Printf("  Confirmations: %d", processed.Confirmations)

		if processed.Amount > 0 { // Only include transactions with non-zero net value
			processedTxs = append(processedTxs, processed)
		}
	}

	// If there's a next page token, log it
	if txResponse.Meta.Paging.NextPageToken != "" {
		log.Printf("More transactions available. Next page token: %s",
			txResponse.Meta.Paging.NextPageToken)
	}

	return processedTxs, nil
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

func (s *SolanaLoader) LoadHistoricalPrice(tokenAddress string, timestamp time.Time) (float64, error) {
	if tokenAddress == "solana/native/sol" {
		unixTime := timestamp.Unix()
		url := fmt.Sprintf(
			"https://api.coingecko.com/api/v3/coins/solana/market_chart/range?vs_currency=usd&from=%d&to=%d",
			unixTime-3600, unixTime+3600, // +/- 1 hour range
		)

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

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return 0, fmt.Errorf(
				"CoinGecko (historical) returned status %d: %s",
				resp.StatusCode,
				string(bodyBytes),
			)
		}

		// TODO: Implement proper historical price parsing
		return 0, nil
	}

	return 1.0, nil
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
