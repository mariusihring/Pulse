package coingecko_requests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	coingecko_types "solana/types/coingecko"
	"strings"

	"github.com/charmbracelet/log"
)

func GetSolanaPrice() (float64, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=solana&vs_currencies=usd"
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to get SOL price: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read SOL price response: %w", err)
	}

	var priceResp map[string]map[string]float64
	if err := json.Unmarshal(body, &priceResp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal SOL price response: %w", err)
	}

	// Expected response: {"solana": {"usd": <price>}}
	price, ok := priceResp["solana"]["usd"]
	if !ok {
		return 0, fmt.Errorf("SOL price not found in response")
	}

	return price, nil
}

func GetCoinGeckoTokenPrices(addresses []string) (map[string]string, error) {
	result := make(map[string]string)
	const batchSize = 30

	for i := 0; i < len(addresses); i += batchSize {
		end := i + batchSize
		if end > len(addresses) {
			end = len(addresses)
		}
		batch := addresses[i:end]
		tokens := strings.Join(batch, ",")
		requestURL := fmt.Sprintf("https://api.geckoterminal.com/api/v2/simple/networks/solana/token_price/%s", tokens)
		log.Info(requestURL)

		resp, err := http.Get(requestURL)
		if err != nil {
			log.Error("Error occurred", "Stack", err)
			return nil, fmt.Errorf("failed to get token prices: %w", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Error("Error occurred", "Stack", err)
			return nil, fmt.Errorf("failed to read token prices response: %w", err)
		}

		var response coingecko_types.PriceResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Error("Error occurred", "Stack", err)
			return nil, fmt.Errorf("failed to unmarshal token prices response: %w", err)
		}

		// For each address in the current batch, add its price or default to "0" if missing.
		for _, address := range batch {
			if price, ok := response.Data.Attributes.TokenPrices[address]; ok {
				result[address] = price
			} else {
				result[address] = "0"
			}
		}
	}

	return result, nil
}
