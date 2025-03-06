package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/mr-tron/base58"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "solana/generated"
	coingecko_requests "solana/requests/coingecko"
	solana_requests "solana/requests/solana"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// validateSolanaAddress checks if the given address is a valid Solana address.
func validateSolanaAddress(address string) error {
	decodedAddr, err := base58.Decode(address)
	if err != nil {
		return errors.New("invalid base58 string")
	}
	if len(decodedAddr) != 32 {
		return errors.New("invalid address length - must be 32 bytes when decoded")
	}
	return nil
}

// getNearestOHLCVPrice returns the open price from the OHLCVS data point
// whose timestamp is closest to the target.
func getNearestOHLCVPrice(data []*pb.PricePoint, target int32) float64 {
	if len(data) == 0 {
		return 0
	}
	nearest := data[0]
	minDiff := absInt32(nearest.Timestamp - target)
	for _, point := range data {
		diff := absInt32(point.Timestamp - target)
		if diff < minDiff {
			minDiff = diff
			nearest = point
		}
	}
	return nearest.Open
}

// absInt32 returns the absolute value of an int32.
func absInt32(n int32) int32 {
	if n < 0 {
		return -n
	}
	return n
}

// isInternalTransfer examines a transaction’s instructions to decide whether
// it represents a transfer between two owned wallets. This implementation
// uses a simple heuristic: for each instruction in the transaction message,
// it retrieves the involved addresses (using the instruction’s account indexes
// and the transaction message’s account_keys). If the first two accounts in the
// instruction are both in the ownedWallets set, we assume it’s an internal transfer.
// (In your environment, you may wish to inspect the instruction data or the program id.)
func isInternalTransfer(tx *pb.Transaction, ownedWallets map[string]bool) bool {
	// If the transaction does not include a result, we cannot inspect its details.
	if tx.Result == nil || tx.Result.Transaction == nil || tx.Result.Transaction.Message == nil {
		return false
	}
	msg := tx.Result.Transaction.Message
	for _, instr := range msg.Instructions {
		// Ensure the instruction has at least two accounts.
		if len(instr.Accounts) < 2 {
			continue
		}
		sourceIdx := instr.Accounts[0]
		destIdx := instr.Accounts[1]
		// Make sure the indexes are valid.
		if int(sourceIdx) >= len(msg.AccountKeys) || int(destIdx) >= len(msg.AccountKeys) {
			continue
		}
		sourceAddr := msg.AccountKeys[sourceIdx]
		destAddr := msg.AccountKeys[destIdx]
		if ownedWallets[sourceAddr] && ownedWallets[destAddr] {
			return true
		}
	}
	return false
}

type server struct {
	pb.UnimplementedWalletServiceServer
}

// AddWallet is your original single-wallet method.
func (s *server) AddWallet(req *pb.WalletRequest, stream pb.WalletService_AddWalletServer) error {
	// Validate wallet address.
	if err := validateSolanaAddress(req.WalletAddress); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid wallet address: %v", err)
	}
	response := &pb.WalletResponse{
		Address: req.WalletAddress,
	}
	// --- Stage 1: Fetch base wallet info and Solana price (10% progress) ---
	wallet, err := solana_requests.RequestAccountInfo(req.WalletAddress)
	if err != nil {
		return status.Errorf(codes.NotFound, "failed to fetch wallet info: %v", err)
	}
	solanaPrice, err := coingecko_requests.GetSolanaPrice()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "failed to get solana price: %v", err)
	}
	response.SolBalance = wallet.SolAmount
	response.SolValue = wallet.SolAmount * solanaPrice
	response.WalletValue = wallet.SolAmount * solanaPrice
	response.LastUpdated = time.Now().UTC().Format(time.RFC3339)
	response.Progress = 10
	if err := stream.Send(response); err != nil {
		log.Error("error sending update", "error", err)
		return err
	}

	// --- Stage 2: Fetch token accounts (20% progress) ---
	accounts, err := solana_requests.RequestTokenAccounts(req.WalletAddress)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "failed to get token accounts: %v", err)
	}
	response.TokenAmount = int32(len(accounts.Result.Value))
	response.Progress = 20
	if err := stream.Send(response); err != nil {
		log.Error("error sending update", "error", err)
		return err
	}

	var addresses []string
	for _, token := range accounts.Result.Value {
		addresses = append(addresses, token.Account.Data.Parsed.Info.Mint)
	}
	currentTokenPrices, err := coingecko_requests.GetCoinGeckoTokenPrices(addresses)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "failed to get token prices: %v", err)
	}

	// --- Stage 3: Process tokens (progress 20% - 60%) ---
	var tokens []*pb.Token
	totalTokens := len(accounts.Result.Value)
	for i, account := range accounts.Result.Value {
		data, _ := solana_requests.GetTokenMetadata(account.Account.Data.Parsed.Info.Mint)
		pool, _ := coingecko_requests.GetTokenPools(account.Account.Data.Parsed.Info.Mint)
		prices, _ := coingecko_requests.GetOHLCVS(pool, "minute", 0, 0)
		var ohlcvsData []*pb.PricePoint
		for _, price := range prices {
			point := pb.PricePoint{
				Timestamp: int32(price[0]),
				Open:      price[1],
				High:      price[2],
				Low:       price[3],
				Close:     price[4],
				Volume:    price[5],
			}
			ohlcvsData = append(ohlcvsData, &point)
		}
		currentPrice, err := strconv.ParseFloat(currentTokenPrices[account.Account.Data.Parsed.Info.Mint], 64)
		if err != nil {
			currentPrice = 0
			log.Error("error parsing token price", "error", err)
			continue
		}
		var invested, pnl float64
		if len(ohlcvsData) > 0 {
			// For demo purposes, use the current time as a placeholder for the purchase timestamp.
			purchaseTimestamp := int32(time.Now().Unix())
			baselinePrice := getNearestOHLCVPrice(ohlcvsData, purchaseTimestamp)
			invested = account.Account.Data.Parsed.Info.TokenAmount.UIAmount * baselinePrice
			pnl = account.Account.Data.Parsed.Info.TokenAmount.UIAmount * (currentPrice - baselinePrice)
		}
		tokens = append(tokens, &pb.Token{
			Name:          data.Result.Content.Metadata.Name,
			Address:       account.Account.Data.Parsed.Info.Mint,
			Pool:          pool,
			Description:   data.Result.Content.Metadata.Description,
			Image:         data.Result.Content.Links.Image,
			Amount:        account.Account.Data.Parsed.Info.TokenAmount.UIAmount,
			Price:         currentPrice,
			Pnl:           pnl,
			Invested:      invested,
			Value:         account.Account.Data.Parsed.Info.TokenAmount.UIAmount * currentPrice,
			HistoryPrices: ohlcvsData,
		})
		response.Tokens = tokens
		response.WalletValue += account.Account.Data.Parsed.Info.TokenAmount.UIAmount * currentPrice
		response.LastUpdated = time.Now().UTC().Format(time.RFC3339)
		response.Progress = float64(20 + float32(i+1)*40/float32(totalTokens))
		if err := stream.Send(response); err != nil {
			log.Error("error sending token update", "error", err)
			return err
		}
	}

	// --- Stage 4: Fetch transaction hashes (70% progress) ---
	hashes, err := solana_requests.GetTransactionHashes(req.WalletAddress)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "failed to get transaction hashes: %v", err)
	}
	response.TransactionAmount = int32(len(hashes))
	response.Progress = 70
	if err := stream.Send(response); err != nil {
		log.Error("error sending transaction hash update", "error", err)
		return err
	}

	// --- Stage 5: Process transactions (70% - 100%) ---
	var transactions []*pb.Transaction
	totalTx := len(hashes)
	successfulTxCount := 0
	var queue []string
	for _, sig := range hashes {
		queue = append(queue, sig.Signature)
	}
	for len(queue) > 0 {
		signature := queue[0]
		queue = queue[1:]
		params := []interface{}{
			signature,
			map[string]interface{}{
				"encoding":                       "json",
				"maxSupportedTransactionVersion": 0,
			},
		}
		result, err := solana_requests.QueryRPCWithRetry("getTransaction", params)
		if err != nil {
			var delaySeconds int
			if n, _ := fmt.Sscanf(err.Error(), "retry after %d seconds", &delaySeconds); n == 1 {
				log.Info("rate limited; retrying", "delaySeconds", delaySeconds, "signature", signature)
				time.Sleep(time.Duration(delaySeconds) * time.Second)
			} else {
				log.Error("error fetching transaction", "signature", signature, "error", err)
				time.Sleep(1 * time.Second)
			}
			queue = append(queue, signature)
			continue
		}
		var tmp map[string]interface{}
		if err = json.Unmarshal([]byte(result), &tmp); err != nil {
			log.Error("error checking JSON for error field", "signature", signature, "error", err)
			queue = append(queue, signature)
			continue
		}
		if _, found := tmp["err"]; found {
			log.Error("RPC returned error", "signature", signature, "response", tmp["err"])
			queue = append(queue, signature)
			continue
		}
		var txResponse pb.Transaction
		opts := protojson.UnmarshalOptions{
			DiscardUnknown: true,
		}
		if err = opts.Unmarshal([]byte(result), &txResponse); err != nil {
			log.Error("error unmarshalling transaction", "signature", signature, "error", err)
			queue = append(queue, signature)
			continue
		}
		transactions = append(transactions, &txResponse)
		response.Transactions = transactions
		response.LastUpdated = time.Now().UTC().Format(time.RFC3339)
		successfulTxCount++
		response.Progress = float64(70 + float32(successfulTxCount)*30/float32(totalTx))
		if err := stream.Send(response); err != nil {
			log.Error("error sending transaction update", "error", err)
		}
	}
	response.Progress = 100
	if err := stream.Send(response); err != nil {
		log.Error("error sending final update", "error", err)
		return err
	}
	return nil
}

// AggregateWallets aggregates data from multiple wallets. It expects a new request message
// (for example, MultiWalletRequest with field wallet_addresses) and returns aggregated information
// while flagging transactions that represent internal transfers between the supplied wallets.
func (s *server) AggregateWallets(req *pb.MultiWalletRequest, stream pb.WalletService_AggregateWalletsServer) error {
	// Build a set of owned wallet addresses.
	ownedWallets := make(map[string]bool)
	for _, addr := range req.WalletAddresses {
		if err := validateSolanaAddress(addr); err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid wallet address %s: %v", addr, err)
		}
		ownedWallets[addr] = true
	}

	aggregated := &pb.WalletResponse{
		Address:      "aggregated",
		LastUpdated:  time.Now().UTC().Format(time.RFC3339),
		Progress:     0,
		Tokens:       []*pb.Token{},
		Transactions: []*pb.Transaction{},
	}

	// --- Stage 1: Aggregate base wallet info (10%) ---
	var totalSolBalance float64
	for _, addr := range req.WalletAddresses {
		wallet, err := solana_requests.RequestAccountInfo(addr)
		if err != nil {
			log.Error("error fetching wallet info", "wallet", addr, "error", err)
			continue
		}
		totalSolBalance += wallet.SolAmount
	}
	solanaPrice, err := coingecko_requests.GetSolanaPrice()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "failed to get solana price: %v", err)
	}
	aggregated.SolBalance = totalSolBalance
	aggregated.SolValue = totalSolBalance * solanaPrice
	aggregated.WalletValue = totalSolBalance * solanaPrice
	aggregated.Progress = 10
	if err := stream.Send(aggregated); err != nil {
		log.Error("error sending aggregated update", "error", err)
		return err
	}

	// --- Stage 2: Fetch token accounts from each wallet (10%) ---
	// Aggregate tokens by mint address.
	tokenMap := make(map[string]*pb.Token)
	for _, addr := range req.WalletAddresses {
		accounts, err := solana_requests.RequestTokenAccounts(addr)
		if err != nil {
			log.Error("error fetching token accounts", "wallet", addr, "error", err)
			continue
		}
		for _, account := range accounts.Result.Value {
			mint := account.Account.Data.Parsed.Info.Mint
			tokenAmount := account.Account.Data.Parsed.Info.TokenAmount.UIAmount
			pool, _ := coingecko_requests.GetTokenPools(mint)
			prices, _ := coingecko_requests.GetOHLCVS(pool, "minute", 0, 0)
			var ohlcvsData []*pb.PricePoint
			for _, price := range prices {
				point := pb.PricePoint{
					Timestamp: int32(price[0]),
					Open:      price[1],
					High:      price[2],
					Low:       price[3],
					Close:     price[4],
					Volume:    price[5],
				}
				ohlcvsData = append(ohlcvsData, &point)
			}
			var baselinePrice float64
			if len(ohlcvsData) > 0 {
				// Using current time as a placeholder for purchase timestamp.
				purchaseTimestamp := int32(time.Now().Unix())
				baselinePrice = getNearestOHLCVPrice(ohlcvsData, purchaseTimestamp)
			}
			invested := tokenAmount * baselinePrice

			if existing, ok := tokenMap[mint]; ok {
				// Merge amounts and preserve the lower cost basis.
				prevAmount := existing.Amount
				existing.Amount += tokenAmount
				prevCostPerToken := existing.Invested / prevAmount
				currentCostPerToken := baselinePrice
				if currentCostPerToken < prevCostPerToken {
					existing.Invested = currentCostPerToken * existing.Amount
				} else {
					existing.Invested = prevCostPerToken * existing.Amount
				}
			} else {
				tokenMap[mint] = &pb.Token{
					Address:       mint,
					Pool:          pool,
					Amount:        tokenAmount,
					Invested:      invested,
					HistoryPrices: ohlcvsData,
				}
			}
		}
	}
	aggregated.Progress = 20
	if err := stream.Send(aggregated); err != nil {
		log.Error("error sending aggregated token accounts update", "error", err)
		return err
	}

	// --- Stage 3: Process tokens (progress 20% - 60%) ---
	// Get current prices for all token mints.
	var tokenMints []string
	for mint := range tokenMap {
		tokenMints = append(tokenMints, mint)
	}
	currentTokenPrices, err := coingecko_requests.GetCoinGeckoTokenPrices(tokenMints)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "failed to get token prices: %v", err)
	}
	var aggregatedTokens []*pb.Token
	totalTokenTypes := len(tokenMap)
	i := 0
	for _, token := range tokenMap {
		// Fill in metadata.
		data, _ := solana_requests.GetTokenMetadata(token.Address)
		token.Name = data.Result.Content.Metadata.Name
		token.Description = data.Result.Content.Metadata.Description
		token.Image = data.Result.Content.Links.Image
		currentPrice, err := strconv.ParseFloat(currentTokenPrices[token.Address], 64)
		if err != nil {
			log.Error("error parsing token price", "token", token.Address, "error", err)
			continue
		}
		token.Price = currentPrice
		token.Value = token.Amount * currentPrice
		token.Pnl = token.Amount*currentPrice - token.Invested
		aggregatedTokens = append(aggregatedTokens, token)
		aggregated.WalletValue += token.Value
		aggregated.LastUpdated = time.Now().UTC().Format(time.RFC3339)
		i++
		aggregated.Progress = float64(20 + 40*float32(i)/float32(totalTokenTypes))
		if err := stream.Send(aggregated); err != nil {
			log.Error("error sending token processing update", "error", err)
			return err
		}
	}
	aggregated.Tokens = aggregatedTokens

	// --- Stage 4: Fetch transaction hashes from all wallets (70% progress) ---
	var allHashes []string
	for _, addr := range req.WalletAddresses {
		hashes, err := solana_requests.GetTransactionHashes(addr)
		if err != nil {
			log.Error("error fetching transaction hashes", "wallet", addr, "error", err)
			continue
		}
		for _, sig := range hashes {
			allHashes = append(allHashes, sig.Signature)
		}
	}
	aggregated.TransactionAmount = int32(len(allHashes))
	aggregated.Progress = 70
	if err := stream.Send(aggregated); err != nil {
		log.Error("error sending aggregated transaction hash update", "error", err)
		return err
	}

	// --- Stage 5: Process transactions (progress 70% - 100%) ---
	var aggregatedTransactions []*pb.Transaction
	totalTx := len(allHashes)
	successfulTxCount := 0
	queue := allHashes
	for len(queue) > 0 {
		signature := queue[0]
		queue = queue[1:]
		params := []interface{}{
			signature,
			map[string]interface{}{
				"encoding":                       "json",
				"maxSupportedTransactionVersion": 0,
			},
		}
		result, err := solana_requests.QueryRPCWithRetry("getTransaction", params)
		if err != nil {
			var delaySeconds int
			if n, _ := fmt.Sscanf(err.Error(), "retry after %d seconds", &delaySeconds); n == 1 {
				log.Info("rate limited; retrying", "delaySeconds", delaySeconds, "signature", signature)
				time.Sleep(time.Duration(delaySeconds) * time.Second)
			} else {
				log.Error("error fetching transaction", "signature", signature, "error", err)
				time.Sleep(1 * time.Second)
			}
			queue = append(queue, signature)
			continue
		}
		var tmp map[string]interface{}
		if err = json.Unmarshal([]byte(result), &tmp); err != nil {
			log.Error("error checking JSON for error field", "signature", signature, "error", err)
			queue = append(queue, signature)
			continue
		}
		if _, found := tmp["err"]; found {
			log.Error("RPC returned error", "signature", signature, "response", tmp["err"])
			queue = append(queue, signature)
			continue
		}
		var txResponse pb.Transaction
		opts := protojson.UnmarshalOptions{
			DiscardUnknown: true,
		}
		if err = opts.Unmarshal([]byte(result), &txResponse); err != nil {
			log.Error("error unmarshalling transaction", "signature", signature, "error", err)
			queue = append(queue, signature)
			continue
		}
		// Check if this transaction is an internal transfer.
		if isInternalTransfer(&txResponse, ownedWallets) {
			log.Info("detected internal transfer; preserving original cost basis", "signature", signature)
			// You might mark this transaction in the response if desired.
		}
		aggregatedTransactions = append(aggregatedTransactions, &txResponse)
		aggregated.Transactions = aggregatedTransactions
		aggregated.LastUpdated = time.Now().UTC().Format(time.RFC3339)
		successfulTxCount++
		aggregated.Progress = float64(70 + 30*float32(successfulTxCount)/float32(totalTx))
		if err := stream.Send(aggregated); err != nil {
			log.Error("error sending transaction update", "error", err)
		}
	}
	aggregated.Progress = 100
	if err := stream.Send(aggregated); err != nil {
		log.Error("error sending final aggregated update", "error", err)
		return err
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterWalletServiceServer(s, &server{})
	log.Info("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
