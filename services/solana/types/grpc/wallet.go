package grpc_types

import (
	solana_types "solana/types/solana_rpc"
	"time"
)

type Wallet struct {
	Address      string                             `json:"address"`
	SolBalance   float64                            `json:"solBalance"`
	SolValue     float64                            `json:"solValue"`
	Value        float64                            `json:"walletValue"`
	Tokens       []Token                            `json:"tokens"`
	Transactions []solana_types.TransactionResponse `json:"transactions"`
	LastUpdated  time.Time                          `json:"last_updated"`
}

type Token struct {
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	Pool           string    `json:"pool"`
	Description    string    `json:"description"`
	Image          string    `json:"image"`
	Amount         float64   `json:"amount"`
	Price          float64   `json:"price"`
	History_prices []float64 `json:"history_prices"`
	PnL            float64   `json:"pnl"`
	Invested       float64   `json:"invested"`
	Value          float64   `json:"value"`
}
