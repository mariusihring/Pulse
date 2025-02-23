package resolvers

import (
	"context"

	"github.com/charmbracelet/log"
	graphql_gen "pulse/graph/generated"
	"pulse/graph/graphql_model"
	grpc "pulse/internal/proto/solana_service"
	"pulse/internal/proto/solana_service/generated"
)

// WalletUpdates is the resolver for the walletUpdates field.
func (r *subscriptionResolver) WalletUpdates(ctx context.Context, walletAddress string) (<-chan *graphql_model.Wallet, error) {
	updates := make(chan *graphql_model.Wallet)

	go func() {
		defer close(updates)

		// Create the gRPC client.
		client, cleanup, err := grpc.NewSolanaGRPCClient()
		if err != nil {
			log.Printf("failed to create gRPC client: %v", err)
			return
		}
		defer cleanup()

		// Build the gRPC request with the wallet address.
		req := &generated.WalletRequest{
			WalletAddress: walletAddress,
		}

		// Call the streaming RPC. This method returns a stream of WalletResponse.
		stream, err := client.AddWallet(ctx, req)
		if err != nil {
			log.Printf("error calling AddWallet: %v", err)
			return
		}

		// Loop and forward each update from the gRPC stream to the subscription channel.
		for {
			grpcResp, err := stream.Recv()
			if err != nil {
				log.Printf("stream closed or error occurred: %v", err)
				return
			}

			// Convert the gRPC response to the GraphQL type.
			resp := convertWalletResponse(grpcResp)

			select {
			case updates <- resp:
			case <-ctx.Done():
				return
			}
		}
	}()

	return updates, nil
}

// convertWalletResponse maps a generated.WalletResponse (proto) to graphql_model.Wallet.
func convertWalletResponse(gr *generated.WalletResponse) *graphql_model.Wallet {
	return &graphql_model.Wallet{
		Address:      gr.Address,
		SolBalance:   gr.SolBalance,
		SolValue:     gr.SolValue,
		WalletValue:  gr.WalletValue,
		LastUpdated:  gr.LastUpdated,
		Tokens:       convertTokens(gr.Tokens),
		Transactions: convertTransactions(gr.Transactions),
	}
}

func convertTokens(gTokens []*generated.Token) []*graphql_model.Token {
	tokens := make([]*graphql_model.Token, len(gTokens))
	for i, t := range gTokens {
		tokens[i] = &graphql_model.Token{
			Name:          t.Name,
			Address:       t.Address,
			Pool:          t.Pool,
			Description:   t.Description,
			Image:         t.Image,
			Amount:        t.Amount,
			Price:         t.Price,
			Pnl:           t.Pnl,
			Invested:      t.Invested,
			Value:         t.Value,
			HistoryPrices: t.HistoryPrices,
		}
	}
	return tokens
}

func convertTransactions(gTrans []*generated.Transaction) []*graphql_model.Transaction {
	transactions := make([]*graphql_model.Transaction, len(gTrans))
	for i, t := range gTrans {
		transactions[i] = &graphql_model.Transaction{
			Jsonrpc: t.Jsonrpc,
			ID:      t.Id,
		}
		if t.Result != nil {
			transactions[i].Result = convertTransactionResult(t.Result)
		}
		if t.Err != nil {
			transactions[i].Err = convertError(t.Err)
		}
	}
	return transactions
}

func convertTransactionResult(gr *generated.TransactionResult) *graphql_model.TransactionResult {
	if gr == nil {
		return nil
	}
	return &graphql_model.TransactionResult{
		BlockTime:   int32(int(gr.BlockTime)), // int64 -> int
		Meta:        convertMeta(gr.Meta),
		Slot:        int32(int(gr.Slot)), // uint64 -> int
		Transaction: convertTransactionData(gr.Transaction),
	}
}

func convertError(ge *generated.Error) *graphql_model.Error {
	if ge == nil {
		return nil
	}
	return &graphql_model.Error{
		Code:    int32(int(ge.Code)), // int32 -> int
		Message: ge.Message,
	}
}

func convertMeta(gm *generated.Meta) *graphql_model.Meta {
	if gm == nil {
		return nil
	}
	return &graphql_model.Meta{
		ComputeUnitsConsumed: int32(int(gm.ComputeUnitsConsumed)), // uint64 -> int
		Fee:                  int32(int(gm.Fee)),                  // uint64 -> int
		InnerInstructions:    convertInnerInstructions(gm.InnerInstructions),
		LogMessages:          gm.LogMessages,
		PostBalances:         convertUint64SliceToIntSlice(gm.PostBalances),
		PreBalances:          convertUint64SliceToIntSlice(gm.PreBalances),
		PostTokenBalances:    convertTokenBalances(gm.PostTokenBalances),
		PreTokenBalances:     convertTokenBalances(gm.PreTokenBalances),
		Rewards:              convertRewards(gm.Rewards),
		Status:               convertStatus(gm.Status),
	}
}

func convertInnerInstructions(gInners []*generated.InnerInstruction) []*graphql_model.InnerInstruction {
	inners := make([]*graphql_model.InnerInstruction, len(gInners))
	for i, inner := range gInners {
		inners[i] = &graphql_model.InnerInstruction{
			Index:        int32(int(inner.Index)), // uint32 -> int
			Instructions: convertInstructions(inner.Instructions),
		}
	}
	return inners
}

func convertInstructions(gInsts []*generated.Instruction) []*graphql_model.Instruction {
	insts := make([]*graphql_model.Instruction, len(gInsts))
	for i, inst := range gInsts {
		insts[i] = &graphql_model.Instruction{
			Accounts:       convertUint32SliceToIntSlice(inst.Accounts),
			Data:           inst.Data,
			ProgramIDIndex: int32(int(inst.ProgramIdIndex)), // uint32 -> int
			StackHeight:    int32(int(inst.StackHeight)),    // uint32 -> int
		}
	}
	return insts
}

func convertTokenBalances(gTokenBalances []*generated.TokenBalance) []*graphql_model.TokenBalance {
	balances := make([]*graphql_model.TokenBalance, len(gTokenBalances))
	for i, tb := range gTokenBalances {
		balances[i] = &graphql_model.TokenBalance{
			AccountIndex:  int32(int(tb.AccountIndex)), // uint32 -> int
			Mint:          tb.Mint,
			Owner:         tb.Owner,
			ProgramID:     tb.ProgramId,
			UiTokenAmount: convertTokenAmount(tb.UiTokenAmount),
		}
	}
	return balances
}

func convertTokenAmount(gAmount *generated.TokenAmount) *graphql_model.TokenAmount {
	if gAmount == nil {
		return nil
	}
	return &graphql_model.TokenAmount{
		Amount:         gAmount.Amount,
		Decimals:       int32(int(gAmount.Decimals)), // uint32 -> int
		UiAmount:       gAmount.UiAmount,
		UiAmountString: gAmount.UiAmountString,
	}
}

func convertRewards(gRewards []*generated.Reward) []*graphql_model.Reward {
	rewards := make([]*graphql_model.Reward, len(gRewards))
	for i, r := range gRewards {
		rewards[i] = &graphql_model.Reward{
			Info: r.Info,
		}
	}
	return rewards
}

func convertStatus(gStatus *generated.Status) *graphql_model.Status {
	if gStatus == nil {
		return nil
	}
	return &graphql_model.Status{
		Ok:           gStatus.Ok,
		ErrorMessage: gStatus.ErrorMessage,
	}
}

func convertTransactionData(gData *generated.TransactionData) *graphql_model.TransactionData {
	if gData == nil {
		return nil
	}
	return &graphql_model.TransactionData{
		Message:    convertTransactionMessage(gData.Message),
		Signatures: gData.Signatures,
	}
}

func convertTransactionMessage(gMsg *generated.TransactionMessage) *graphql_model.TransactionMessage {
	if gMsg == nil {
		return nil
	}
	return &graphql_model.TransactionMessage{
		AccountKeys:         gMsg.AccountKeys,
		AddressTableLookups: convertAddressTableLookups(gMsg.AddressTableLookups),
		Header:              convertMessageHeader(gMsg.Header),
		Instructions:        convertInstructions(gMsg.Instructions),
		RecentBlockhash:     gMsg.RecentBlockhash,
	}
}

func convertAddressTableLookups(gLookups []*generated.AddressTableLookup) []*graphql_model.AddressTableLookup {
	lookups := make([]*graphql_model.AddressTableLookup, len(gLookups))
	for i, l := range gLookups {
		lookups[i] = &graphql_model.AddressTableLookup{
			AccountKey:      l.AccountKey,
			ReadonlyIndexes: convertUint32SliceToIntSlice(l.ReadonlyIndexes),
			WritableIndexes: convertUint32SliceToIntSlice(l.WritableIndexes),
		}
	}
	return lookups
}

func convertMessageHeader(gHeader *generated.MessageHeader) *graphql_model.MessageHeader {
	if gHeader == nil {
		return nil
	}
	return &graphql_model.MessageHeader{
		NumReadonlySignedAccounts:   int32(int(gHeader.NumReadonlySignedAccounts)),
		NumReadonlyUnsignedAccounts: int32(int(gHeader.NumReadonlyUnsignedAccounts)),
		NumRequiredSignatures:       int32(int(gHeader.NumRequiredSignatures)),
	}
}

// Helper functions to convert slices of unsigned integers to slices of int.
func convertUint64SliceToIntSlice(in []uint64) []int {
	out := make([]int, len(in))
	for i, v := range in {
		out[i] = int(v)
	}
	return out
}

func convertUint32SliceToIntSlice(in []uint32) []int {
	out := make([]int, len(in))
	for i, v := range in {
		out[i] = int(v)
	}
	return out
}

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graphql_gen.SubscriptionResolver {
	return &subscriptionResolver{r}
}

type subscriptionResolver struct{ *Resolver }
