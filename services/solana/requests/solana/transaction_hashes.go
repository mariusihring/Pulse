package solana_requests

import (
	"encoding/json"
	solana_types "solana/types/solana_rpc"

	"github.com/charmbracelet/log"
)

func GetTransactionHashes(address string) ([]solana_types.WalletTransactionHashResponse, error) {

	data, err := queryRPC("getSignaturesForAddress", []interface{}{address})
	if err != nil {

	}
	var sigResponse solana_types.SignaturesForAddressResponse
	if err := json.Unmarshal([]byte(data), &sigResponse); err != nil {
		log.Error("Error unmarshalling getSignaturesForAddress response", "Stack", err)
		return []solana_types.WalletTransactionHashResponse{}, err
	}

	return sigResponse.Result, nil

}
