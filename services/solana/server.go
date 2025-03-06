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

	"solana/generated"
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
// whose timestamp is closest to the target timestamp.
func getNearestOHLCVPrice(data []*generated.PricePoint, target int32) float64 {
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

type server struct {
	generated.UnimplementedWalletServiceServer
}

// AddWallet is the main gRPC method that performs all operations.
func (s *server) AddWallet(req *generated.WalletRequest, stream generated.WalletService_AddWalletServer) error {
	// Validate the Solana address.
	if err := validateSolanaAddress(req.WalletAddress); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid wallet address: %v", err)
	}
	response := &generated.WalletResponse{
		Address: req.WalletAddress,
	}
	// Stage 1: Get the base wallet and current Solana price (10% progress).
	wallet, err := solana_requests.RequestAccountInfo(req.WalletAddress)
	if err != nil {
		return status.Errorf(codes.NotFound, "invalid wallet address: %v", err)
	}
	solanaPrice, err := coingecko_requests.GetSolanaPrice()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not get solana price: %v", err)
	}
	response.SolBalance = wallet.SolAmount
	response.SolValue = wallet.SolAmount * solanaPrice
	response.WalletValue = wallet.SolAmount * solanaPrice
	response.LastUpdated = time.Now().UTC().Format(time.RFC3339)
	response.Progress = 10 // 10% complete.
	if err := stream.Send(response); err != nil {
		log.Printf("error sending update: %v", err)
		return err
	}

	// Stage 2: Get the token accounts (next 10%).
	accounts, err := solana_requests.RequestTokenAccounts(req.WalletAddress)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not get tokens: %v", err)
	}
	response.TokenAmount = int32(len(accounts.Result.Value))
	response.Progress = 20 // 20% complete.
	if err := stream.Send(response); err != nil {
		log.Printf("error sending update: %v", err)
		return err
	}

	var addresses []string
	for _, token := range accounts.Result.Value {
		addresses = append(addresses, token.Account.Data.Parsed.Info.Mint)
	}
	// Get current token prices.
	currentTokenPrices, err := coingecko_requests.GetCoinGeckoTokenPrices(addresses)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not get token prices: %v", err)
	}

	// Stage 3: Process tokens (40% total weight).
	var tokens []*generated.Token
	totalTokens := len(accounts.Result.Value)
	for i, account := range accounts.Result.Value {
		// Fetch token metadata and pool data.
		data, _ := solana_requests.GetTokenMetadata(account.Account.Data.Parsed.Info.Mint)
		pool, _ := coingecko_requests.GetTokenPools(account.Account.Data.Parsed.Info.Mint)

		// Get historical OHLCVS data.
		prices, _ := coingecko_requests.GetOHLCVS(pool, "minute", 0, 0)
		var ohlcvsData []*generated.PricePoint
		for _, price := range prices {
			point := generated.PricePoint{
				Timestamp: int32(price[0]),
				Open:      price[1],
				High:      price[2],
				Low:       price[3],
				Close:     price[4],
				Volume:    price[5],
			}
			ohlcvsData = append(ohlcvsData, &point)
		}

		// Get the current token price.
		currentPrice, err := strconv.ParseFloat(currentTokenPrices[account.Account.Data.Parsed.Info.Mint], 64)
		if err != nil {
			log.Error("Error parsing current token price", "error", err)
			continue
		}

		// Calculate invested amount and PnL.
		var invested, pnl float64
		if len(ohlcvsData) > 0 {
			// For demonstration, we use the current time as a placeholder purchase timestamp.
			purchaseTimestamp := int32(time.Now().Unix())
			baselinePrice := getNearestOHLCVPrice(ohlcvsData, purchaseTimestamp)
			invested = account.Account.Data.Parsed.Info.TokenAmount.UIAmount * baselinePrice
			pnl = account.Account.Data.Parsed.Info.TokenAmount.UIAmount * (currentPrice - baselinePrice)
		}

		// Append the token with updated details.
		tokens = append(tokens, &generated.Token{
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
		// Update progress for tokens: 20% + portion of 40%.
		response.Progress = float64(20 + float32(i+1)*40/float32(totalTokens))
		if err := stream.Send(response); err != nil {
			log.Printf("error sending token update: %v", err)
			return err
		}
	}

	// Stage 4: Transaction hashes fetch (10%).
	// If there were no tokens, ensure we set progress accordingly.
	if totalTokens == 0 {
		response.Progress = 20 + 40 // 60%
	}
	hashes, err := solana_requests.GetTransactionHashes(req.WalletAddress)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not get transactions: %v", err)
	}
	response.TransactionAmount = int32(len(hashes))
	response.Progress = 70 // 70% complete.
	if err := stream.Send(response); err != nil {
		log.Printf("error sending transaction hash update: %v", err)
		return err
	}

	// Stage 5: Process transactions (30%).
	var transactions []*generated.Transaction
	totalTx := len(hashes)
	successfulTxCount := 0

	// Create a queue of signatures.
	var queue []string
	for _, sig := range hashes {
		queue = append(queue, sig.Signature)
	}

	for len(queue) > 0 {
		// Dequeue the first signature.
		signature := queue[0]
		queue = queue[1:]

		params := []interface{}{
			signature,
			map[string]interface{}{
				"encoding":                       "json",
				"maxSupportedTransactionVersion": 0,
			},
		}

		// Fetch transaction data with retry.
		result, err := solana_requests.QueryRPCWithRetry("getTransaction", params)
		if err != nil {
			var delaySeconds int
			if n, _ := fmt.Sscanf(err.Error(), "retry after %d seconds", &delaySeconds); n == 1 {
				log.Info("Rate limited. Retrying after delay", "delaySeconds", delaySeconds, "signature", signature)
				time.Sleep(time.Duration(delaySeconds) * time.Second)
			} else {
				log.Error("Error fetching transaction", "signature", signature, "error", err)
				time.Sleep(1 * time.Second)
			}
			queue = append(queue, signature)
			continue
		}

		// Check JSON response for an error field.
		var tmp map[string]interface{}
		if err = json.Unmarshal([]byte(result), &tmp); err != nil {
			log.Error("Error checking JSON for error field", "signature", signature, "error", err)
			queue = append(queue, signature)
			continue
		}
		if _, found := tmp["err"]; found {
			log.Error("RPC returned an error", "signature", signature, "response", tmp["err"])
			queue = append(queue, signature)
			continue
		}

		var txResponse generated.Transaction
		opts := protojson.UnmarshalOptions{
			DiscardUnknown: true,
		}
		if err = opts.Unmarshal([]byte(result), &txResponse); err != nil {
			log.Error("Error unmarshalling transaction", "signature", signature, "error", err)
			queue = append(queue, signature)
			continue
		}

		transactions = append(transactions, &txResponse)
		response.Transactions = transactions
		response.LastUpdated = time.Now().UTC().Format(time.RFC3339)
		successfulTxCount++
		// Update progress: 70% + portion of the final 30% based on successful transactions.
		response.Progress = float64(70 + float32(successfulTxCount)*30/float32(totalTx))
		if err := stream.Send(response); err != nil {
			log.Printf("error sending transaction update: %v", err)
		}
	}

	// Final update: 100% complete.
	response.Progress = 100
	log.Info("Scan done for", "wallet", req.WalletAddress)
	if err := stream.Send(response); err != nil {
		log.Printf("error sending final update: %v", err)
		return err
	}
	return nil
}

func main() {
	// Listen on port 50051.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	generated.RegisterWalletServiceServer(s, &server{})
	log.Info("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
