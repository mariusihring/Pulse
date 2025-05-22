package coingecko_requests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	coingecko_types "solana/types/coingecko"

	"github.com/charmbracelet/log"
)

func GetOHLCVS(address string, timeframe string, start int64, end int64) ([][]float64, error) {
	request_url := fmt.Sprintf("https://api.geckoterminal.com/api/v2/networks/solana/pools/%s/ohlcv/%s?currency=usd", address, timeframe)
	resp, err := http.Get(request_url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return nil, err
	}
	var response coingecko_types.OHLCVSResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		log.Error("Error occured", "Stack", err)
		return nil, err
	}
	return response.Data.Attributes.OHLCVList, nil
}
