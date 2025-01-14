package loaders

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"pulse/internal/config"
	"pulse/internal/services/loaders/types"
	"time"
)

type SolanaLoader struct {
	apiKeys ApiKeys
	client  *http.Client
}

type ApiKeys struct {
	blockdaemon string
	coinbase    string
}

func NewSolanaLoader(cfg *config.Config) *SolanaLoader {
	apiKeys := ApiKeys{blockdaemon: cfg.ApiKeys.Blockdaemon, coinbase: cfg.ApiKeys.Coinbase}
	client := &http.Client{Timeout: 15 * time.Second}
	return &SolanaLoader{apiKeys, client}
}

func (s *SolanaLoader) LoadTransactions(address string) ([]types.TokenBalance, error) {
	return nil, nil
}

func (s *SolanaLoader) LoadWallet(address string) ([]types.TokenBalance, error) {
	// TODO: load the solana address and return the tokens and their balances. This should only be called when we create a new Subwallet which is of type Solana
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

	// Check status code
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
	fmt.Println(balances)
	return balances, nil
}

func (s *SolanaLoader) LoadCurrentSolanaPrice() (float64, error) {
	//TODO: load from coinbase and parse into struct
	//
	//
	//TODO: calculateUsdPrice here
	price, err := calculateUsdPrice(0, 0)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (s *SolanaLoader) LoadTokenPrice(address string) (float64, error) {
	//TODO: load from blockdaemon and parse into struct
	//
	//
	//TODO: calculateUsdPrice here
	price, err := calculateUsdPrice(0, 0)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func calculateUsdPrice(price float64, decimals int) (float64, error) {
	return 0, nil
}
