package requests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"solana/types"

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
	tokens := strings.Join(addresses, ",")
	request_url := fmt.Sprintf("https://api.geckoterminal.com/api/v2/simple/networks/solana/token_price/%s", tokens)
	resp, err := http.Get(request_url)
	if err != nil {
		log.Error("Error occured", "Stack", err)
		return nil, fmt.Errorf("failed to get token prices: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error occured", "Stack", err)
		return nil, fmt.Errorf("failed to read token prices response: %w", err)
	}
	var response types.CoinGeckoPriceResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		log.Error("Error occured", "Stack", err)
		return nil, fmt.Errorf("failed to unmarshal token prices response: %w", err)
	}
	return response.Data.Attributes.TokenPrices, nil
}
