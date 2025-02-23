package solana_types

import (
	"encoding/json"
	"fmt"
)

// TransactionResponse is the top-level response type.
type TransactionResponse struct {
	JsonRPC string             `json:"jsonrpc"`
	Result  *TransactionResult `json:"result"`
	Error   *SolanaError       `json:"error"`
	ID      int                `json:"id"`
}

type SolanaError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Version is a custom type that can handle both string and numeric JSON values
type Version string

// UnmarshalJSON implements json.Unmarshaler interface
func (v *Version) UnmarshalJSON(data []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*v = Version(s)
		return nil
	}

	// Try number
	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*v = Version(n.String())
		return nil
	}

	return fmt.Errorf("version must be either string or number")
}

// Result contains the main result fields.
type TransactionResult struct {
	BlockTime   int64       `json:"blockTime"`
	Meta        Meta        `json:"meta"`
	Slot        int         `json:"slot"`
	Transaction Transaction `json:"transaction"`
	Version     Version     `json:"version"`
}

// Meta holds metadata about the transaction.
type Meta struct {
	ComputeUnitsConsumed int                `json:"computeUnitsConsumed"`
	Err                  interface{}        `json:"err"` // Can be further refined if error details are known.
	Fee                  int                `json:"fee"`
	InnerInstructions    []InnerInstruction `json:"innerInstructions"`
	LoadedAddresses      LoadedAddresses    `json:"loadedAddresses"`
	LogMessages          []string           `json:"logMessages"`
	PostBalances         []int64            `json:"postBalances"`
	PostTokenBalances    []TokenBalance     `json:"postTokenBalances"`
	PreBalances          []int64            `json:"preBalances"`
	PreTokenBalances     []TokenBalance     `json:"preTokenBalances"`
	Rewards              []interface{}      `json:"rewards"` // Adjust type if rewards have a defined structure.
	Status               Status             `json:"status"`
}

// InnerInstruction groups a set of instructions with an index.
type InnerInstruction struct {
	Index        int           `json:"index"`
	Instructions []Instruction `json:"instructions"`
}

// Instruction represents a single instruction.
type Instruction struct {
	Accounts       []int  `json:"accounts"`
	Data           string `json:"data"`
	ProgramIdIndex int    `json:"programIdIndex"`
	// stackHeight can be null so we use a pointer.
	StackHeight *int `json:"stackHeight"`
}

// LoadedAddresses lists read-only and writable addresses.
type LoadedAddresses struct {
	Readonly []string `json:"readonly"`
	Writable []string `json:"writable"`
}

// TokenBalance represents the balance details for a token account.
type TokenBalance struct {
	AccountIndex  int                    `json:"accountIndex"`
	Mint          string                 `json:"mint"`
	Owner         string                 `json:"owner"`
	ProgramId     string                 `json:"programId"`
	UiTokenAmount TransactionTokenAmount `json:"uiTokenAmount"`
}

// TokenAmount holds amount details in different formats.
type TransactionTokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       int     `json:"decimals"`
	UiAmount       float64 `json:"uiAmount"`
	UiAmountString string  `json:"uiAmountString"`
}

// Status represents the transaction status.
type Status struct {
	Ok interface{} `json:"Ok"`
}

// Transaction contains the signed transaction data.
type Transaction struct {
	Message    Message  `json:"message"`
	Signatures []string `json:"signatures"`
}

// Message contains details of the transaction message.
type Message struct {
	AccountKeys         []string             `json:"accountKeys"`
	AddressTableLookups []AddressTableLookup `json:"addressTableLookups"`
	Header              Header               `json:"header"`
	Instructions        []MessageInstruction `json:"instructions"`
	RecentBlockhash     string               `json:"recentBlockhash"`
}

// AddressTableLookup represents an address table lookup.
type AddressTableLookup struct {
	AccountKey      string `json:"accountKey"`
	ReadonlyIndexes []int  `json:"readonlyIndexes"`
	WritableIndexes []int  `json:"writableIndexes"`
}

// Header defines the header information for the transaction message.
type Header struct {
	NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
	NumRequiredSignatures       int `json:"numRequiredSignatures"`
}

// MessageInstruction represents an instruction within the transaction message.
type MessageInstruction struct {
	Accounts       []int  `json:"accounts"`
	Data           string `json:"data"`
	ProgramIdIndex int    `json:"programIdIndex"`
	// stackHeight can be null so we use a pointer.
	StackHeight *int `json:"stackHeight"`
}
