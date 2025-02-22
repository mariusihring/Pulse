package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const solanaRpcEndpoint = "https://api.mainnet-beta.solana.com"

// GetBalance fetches the SOL Balance (in SOL) for the given address
func GetBalance(address string) (float64, error) {
	reqPayload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBalance",
		"params":  []interface{}{address},
	}

	reqBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal balance request: %w", err)
	}

	resp, err := http.Post(solanaRpcEndpoint, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return 0, fmt.Errorf("failed to post balance request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("non-OK HTTP status for balance: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read balance response: %w", err)
	}

	var rpcResp GetBalanceResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal balance response: %w", err)
	}

	// Convert lamports to SOL (1 SOL = 1e9 lamports)
	solBalance := float64(rpcResp.Result.Value) / 1e9
	return solBalance, nil
}

// GetTokenAccounts retrieves all token accounts with balances for the given wallet.
func GetTokenAccounts(address string) ([]TokenInfo, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTokenAccountsByOwner",
		"params": []interface{}{
			address,
			// Filter for SPL tokens using the SPL Token Program ID.
			map[string]interface{}{
				"programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
			},
			// Request parsed JSON so that we can easily extract token details.
			map[string]interface{}{
				"encoding": "jsonParsed",
			},
		},
	}

	reqBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal token accounts request: %w", err)
	}

	resp, err := http.Post(solanaRpcEndpoint, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to post token accounts request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK HTTP status for token accounts: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token accounts response: %w", err)
	}

	var tokenResp GetTokenAccountsResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token accounts response: %w", err)
	}

	var tokens []TokenInfo
	for _, account := range tokenResp.Result.Value {
		info := account.Account.Data.Parsed.Info
		token := TokenInfo{
			Mint:         info.Mint,
			Amount:       info.TokenAmount.UiAmount,
			Decimals:     info.TokenAmount.Decimals,
			TokenAccount: account.Pubkey,
			UsdValue:     0, // will be updated later
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// GetSolPrice retrieves the current USD price for SOL from CoinGecko.
func GetSolPrice() (float64, error) {
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

// GetTokenPrice retrieves the current USD price for an SPL token given its mint address.
func GetTokenPrice(mint string) (float64, error) {
	url := fmt.Sprintf("https://public-api.birdeye.so/defi/price?address=%s", mint)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-chain", "solana")
	req.Header.Add("X-API-KEY", "e18b4bd710814fbb9580596543952f5b")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to get token price: %w", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read token price: %w", err)
	}

	var priceResp *BirdeyeTokenPriceResponse
	if err := json.Unmarshal(body, &priceResp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal token price: %w", err)
	}

	return priceResp.Data.Value, nil
}
