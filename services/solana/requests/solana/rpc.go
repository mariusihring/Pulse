package solana_requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
)

const solanaRPC = "https://api.mainnet-beta.solana.com"

func queryRPC(method string, params []interface{}) (string, error) {
	requestPayload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	}

	requestBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Error("Error occured", "Stack", err)
	}
	resp, err := http.Post(solanaRPC, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		log.Errorf("Error indenting JSON:", err)
		return "", err
	}

	return prettyJSON.String(), nil
}
func QueryRPCWithRetry(method string, params []interface{}) (string, error) {
	requestPayload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	}

	requestBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Error("Error marshalling request", "Stack", err)
		return "", err
	}

	resp, err := http.Post(solanaRPC, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		log.Error("HTTP Post error", "Stack", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check for rate limiting or other errors via status code.
	if resp.StatusCode != http.StatusOK {
		retryAfterHeader := resp.Header.Get("Retry-After")
		if retryAfterHeader != "" {
			if delaySec, err := strconv.Atoi(retryAfterHeader); err == nil {
				return "", fmt.Errorf("retry after %d seconds", delaySec)
			}
		}
		return "", fmt.Errorf("received non-200 status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body", "Stack", err)
		return "", err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		log.Errorf("Error indenting JSON", err)
		return "", err
	}

	return prettyJSON.String(), nil
}
