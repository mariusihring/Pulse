package solana_requests

import (
	"encoding/json"
	solana_types "solana/types/solana_rpc"

	"github.com/charmbracelet/log"
)

func RequestTokenAccounts(address string) (solana_types.TokenAccountsByOwnerResponse, error) {
	data, err := queryRPC("getTokenAccountsByOwner", []interface{}{
		address,
		map[string]interface{}{
			"programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
		},
		map[string]interface{}{
			"encoding": "jsonParsed",
		},
	})

	if err != nil {
		return solana_types.TokenAccountsByOwnerResponse{}, err
	}
	var response solana_types.TokenAccountsByOwnerResponse
	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		log.Error("Error occured", "Stack", err)
		return solana_types.TokenAccountsByOwnerResponse{}, err
	}
	return response, nil
}
