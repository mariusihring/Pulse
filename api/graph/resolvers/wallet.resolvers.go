package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.63

import (
	"context"
	"encoding/json"
	graphql_gen "pulse/graph/generated"
	"pulse/graph/graphql_model"
	grpc "pulse/internal/proto/solana_service"
	"pulse/internal/proto/solana_service/generated"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

// StartWalletUpdate is the resolver for the startWalletUpdate field.
func (r *mutationResolver) StartWalletUpdate(ctx context.Context, walletAddress string) (*graphql_model.Job, error) {
	// Generate a unique job ID.
	jobID := uuid.New().String()

	job := &graphql_model.Job{
		ID:            jobID,
		WalletAddress: walletAddress,
	}

	// Start the background process.
	go func(jobID, walletAddress string) {
		bgCtx := context.Background()
		// Create the gRPC client.
		client, cleanup, err := grpc.NewSolanaGRPCClient()
		if err != nil {
			log.Printf("failed to create gRPC client: %v", err)
			return
		}
		defer cleanup()

		req := &generated.WalletRequest{
			WalletAddress: walletAddress,
		}

		stream, err := client.AddWallet(bgCtx, req)
		if err != nil {
			log.Infof("error calling AddWallet: %v", err)
			return
		}

		progress := 0

		for {
			grpcResp, err := stream.Recv()
			if err != nil {
				log.Infof("stream closed or error occurred: %v", err)
				break
			}

			walletUpdate := &graphql_model.WalletUpdate{
				JobID:    jobID,
				Progress: grpcResp.Progress,
				Wallet:   convertWalletResponse(grpcResp),
			}

			if pubErr := r.RedisPubSub.Publish(bgCtx, jobID, walletUpdate); pubErr != nil {
				log.Errorf("failed to publish update: %v", pubErr)
			}

			if progress < 100 {
				progress += 10
			}
		}
	}(jobID, walletAddress)

	return job, nil
}

// WalletUpdates is the resolver for the walletUpdates field.
func (r *subscriptionResolver) WalletUpdates(ctx context.Context, jobID string) (<-chan *graphql_model.WalletUpdate, error) {
	pubsub := r.RedisPubSub.Subscribe(ctx, jobID)
	updates := make(chan *graphql_model.WalletUpdate)

	go func() {
		defer func() {
			if err := pubsub.Close(); err != nil {
				log.Errorf("Failed to close Redis pubsub", "err", err)
			}
			close(updates)
		}()

		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return
				}

				log.Errorf("Error receiving message", "err", err)
				return
			}

			var update graphql_model.WalletUpdate

			if err := json.Unmarshal([]byte(msg.Payload), &update); err != nil {
				log.Errorf("Failed to unmarshal message", "err", err)
				continue
			}

			select {
			case updates <- &update:
			case <-ctx.Done():
				return
			}
		}
	}()
	return updates, nil
}

// Subscription returns graphql_gen.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graphql_gen.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }

func convertWalletResponse(gr *generated.WalletResponse) *graphql_model.Wallet {
	return &graphql_model.Wallet{
		Address:      gr.Address,
		Name:         "My Wallet",
		Description:  "Description",
		Network:      "solana",
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
			HistoryPrices: convertHistoryPrices(t.HistoryPrices),
		}
	}
	return tokens
}
func convertHistoryPrices(gPrices []*generated.PricePoint) []*graphql_model.PricePoint {
	prices := make([]*graphql_model.PricePoint, len(gPrices))
	for i, p := range gPrices {
		prices[i] = &graphql_model.PricePoint{
			Timestamp: p.Timestamp,
			Open:      p.Open,
			High:      p.High,
			Low:       p.Low,
			Close:     p.Close,
			Volume:    p.Volume,
		}
	}
	return prices
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
		BlockTime:   int32(gr.BlockTime), // int64 -> int
		Meta:        convertMeta(gr.Meta),
		Slot:        int32(gr.Slot), // uint64 -> int
		Transaction: convertTransactionData(gr.Transaction),
	}
}
func convertError(ge *generated.Error) *graphql_model.Error {
	if ge == nil {
		return nil
	}
	return &graphql_model.Error{
		Code:    int32(ge.Code), // int32 -> int
		Message: ge.Message,
	}
}
func convertMeta(gm *generated.Meta) *graphql_model.Meta {
	if gm == nil {
		return nil
	}
	return &graphql_model.Meta{
		ComputeUnitsConsumed: int32(gm.ComputeUnitsConsumed), // uint64 -> int
		Fee:                  int32(gm.Fee),                  // uint64 -> int
		InnerInstructions:    convertInnerInstructions(gm.InnerInstructions),
		LogMessages:          gm.LogMessages,
		PostBalances:         convertUint64SliceToInt32Slice(gm.PostBalances),
		PreBalances:          convertUint64SliceToInt32Slice(gm.PreBalances),
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
			Index:        int32(inner.Index), // uint32 -> int
			Instructions: convertInstructions(inner.Instructions),
		}
	}
	return inners
}
func convertInstructions(gInsts []*generated.Instruction) []*graphql_model.Instruction {
	insts := make([]*graphql_model.Instruction, len(gInsts))
	for i, inst := range gInsts {
		insts[i] = &graphql_model.Instruction{
			Accounts:       convertUint32SliceToInt32Slice(inst.Accounts),
			Data:           inst.Data,
			ProgramIDIndex: int32(inst.ProgramIdIndex), // uint32 -> int
			StackHeight:    int32(inst.StackHeight),    // uint32 -> int
		}
	}
	return insts
}
func convertTokenBalances(gTokenBalances []*generated.TokenBalance) []*graphql_model.TokenBalance {
	balances := make([]*graphql_model.TokenBalance, len(gTokenBalances))
	for i, tb := range gTokenBalances {
		balances[i] = &graphql_model.TokenBalance{
			AccountIndex:  int32(tb.AccountIndex), // uint32 -> int
			Mint:          tb.Mint,
			Owner:         tb.Owner,
			ProgramID:     tb.ProgramId,
			UITokenAmount: convertTokenAmount(tb.UiTokenAmount),
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
		Decimals:       int32(gAmount.Decimals), // uint32 -> int
		UIAmount:       gAmount.UiAmount,
		UIAmountString: gAmount.UiAmountString,
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
		Ok:           nil,
		ErrorMessage: nil,
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
			ReadonlyIndexes: convertUint32SliceToInt32Slice(l.ReadonlyIndexes),
			WritableIndexes: convertUint32SliceToInt32Slice(l.WritableIndexes),
		}
	}
	return lookups
}
func convertMessageHeader(gHeader *generated.MessageHeader) *graphql_model.MessageHeader {
	if gHeader == nil {
		return nil
	}
	return &graphql_model.MessageHeader{
		NumReadonlySignedAccounts:   int32(gHeader.NumReadonlySignedAccounts),
		NumReadonlyUnsignedAccounts: int32(gHeader.NumReadonlyUnsignedAccounts),
		NumRequiredSignatures:       int32(gHeader.NumRequiredSignatures),
	}
}
func convertUint64SliceToInt32Slice(in []uint64) []int32 {
	out := make([]int32, len(in))
	for i, v := range in {
		out[i] = int32(v)
	}
	return out
}
func convertUint32SliceToInt32Slice(in []uint32) []int32 {
	out := make([]int32, len(in))
	for i, v := range in {
		out[i] = int32(v)
	}
	return out
}
