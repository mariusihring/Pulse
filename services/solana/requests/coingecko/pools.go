package coingecko_requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	coingecko_types "solana/types/coingecko"
)

func GetTokenPools(address string) (string, error) {
	request_url := fmt.Sprintf("https://api.geckoterminal.com/api/v2/networks/solana/tokens/%s/pools?page=1", address)
	resp, err := http.Get(request_url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var response coingecko_types.PoolResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return "", err
	}
	if len(response.Data) == 0 {
		return "", errors.New("no pools")
	}
	return response.Data[0].Attributes.Address, nil
}
