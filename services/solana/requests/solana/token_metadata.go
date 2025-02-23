package solana_requests

import (
	"encoding/json"
	solana_types "solana/types/solana_rpc"

	"github.com/charmbracelet/log"
)

func GetTokenMetadata(address string) (solana_types.GetTokenMetaDataResponse, error) {
	data, err := queryRPC("getAsset", []interface{}{address})
	if err != nil {
		return solana_types.GetTokenMetaDataResponse{}, err
	}
	var response solana_types.GetTokenMetaDataResponse
	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		log.Error("Error occured", "Stack", err)
	}
	return response, nil
}
