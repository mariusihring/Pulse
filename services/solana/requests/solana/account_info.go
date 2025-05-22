package solana_requests

import (
	"encoding/json"
	"math"
	solana_types "solana/types/solana_rpc"

	"github.com/charmbracelet/log"
)

type Wallet struct {
	AccountInfo solana_types.AccountInfoResponse
	SolAmount   float64
}

func RequestAccountInfo(address string) (Wallet, error) {
	data, err := queryRPC("getAccountInfo", []interface{}{address})
	if err != nil {
		return Wallet{}, err
	}
	var response solana_types.AccountInfoResponse
	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		log.Error(" get account Error occured", "Stack", err)
	}
	balance, err := queryRPC("getBalance", []interface{}{address})
	if err != nil {
		return Wallet{}, err
	}
	var walletresponse solana_types.WalletResponse
	err = json.Unmarshal([]byte(balance), &walletresponse)
	if err != nil {
		log.Error("Balance unmarshaling Error occured", "Stack", err)
	}
	divisor := math.Pow10(9)
	floatValue := float64(walletresponse.Result.Value) / divisor
	return Wallet{response, floatValue}, nil
}
