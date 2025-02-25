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

// validateSolanaAddress checks if the given address is a valid Solana address
func validateSolanaAddress(address string) error {
	// Check if the address is a valid base58 string
	decodedAddr, err := base58.Decode(address)
	if err != nil {
		return errors.New("invalid base58 string")
	}

	// Solana addresses are 32 bytes long
	if len(decodedAddr) != 32 {
		return errors.New("invalid address length - must be 32 bytes when decoded")
	}

	return nil
}

type server struct {
	pb.UnimplementedWalletServiceServer
}

// GetWalletInfo is the main gRPC method that performs all operations.
func (s *server) AddWallet(req *pb.WalletRequest, stream pb.WalletService_AddWalletServer) error {
	// Validate the Solana address
	if err := validateSolanaAddress(req.WalletAddress); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid wallet address: %v", err)
	}
	response := &pb.WalletResponse{
		Address: req.WalletAddress,
	}
	// Get the base wallet
	wallet, err := solana_requests.RequestAccountInfo(req.WalletAddress)
	if err != nil {
		return status.Errorf(codes.NotFound, "invalid wallet address: %v", err)
	}
	// Get the current solana price
	solana_price, err := coingecko_requests.GetSolanaPrice()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not get solana price: %v", err)
	}

	// Return the updated wallet with solana value and amount
	response.SolBalance = wallet.SolAmount
	response.SolValue = wallet.SolAmount * solana_price
	response.WalletValue = wallet.SolAmount * solana_price
	response.LastUpdated = time.Now().UTC().Format(time.RFC3339)

	if err := stream.Send(response); err != nil {
		log.Printf("error sending update: %v", err)
		return err
	}

	// Get the tokens in the wallet
	accounts, err := solana_requests.RequestTokenAccounts(req.WalletAddress)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not get tokens: %v", err)
	}

	var addresses []string
	for _, token := range accounts.Result.Value {
		addresses = append(addresses, token.Account.Data.Parsed.Info.Mint)
	}
	// Get the current token prices
	current_token_prices, err := coingecko_requests.GetCoinGeckoTokenPrices(addresses)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not get token prices: %v", err)
	}

	// Return the updated wallet with all the tokens
	var tokens []*pb.Token
	priceHistories := make(map[string][][]float64)
	for _, account := range accounts.Result.Value {
		data, _ := solana_requests.GetTokenMetadata(account.Account.Data.Parsed.Info.Mint)
		pool, _ := coingecko_requests.GetTokenPools(account.Account.Data.Parsed.Info.Mint)

		prices, _ := coingecko_requests.GetOHLCVS(pool, "minute", 0, 0)
		var ohlcvs_data []*pb.PricePoint
		for _, price := range prices {
			point := pb.PricePoint{
				Timestamp: int32(price[0]),
				Open:      price[1],
				High:      price[2],
				Low:       price[3],
				Close:     price[4],
				Volume:    price[5],
			}
			ohlcvs_data = append(ohlcvs_data, &point)
		}
		priceHistories[account.Account.Data.Parsed.Info.Mint] = prices
		f, err := strconv.ParseFloat(current_token_prices[account.Account.Data.Parsed.Info.Mint], 64)
		if err != nil {
			log.Error("Error occurred", "Stack", err)
			continue
		}
		tokens = append(tokens, &pb.Token{
			Name:          data.Result.Content.Metadata.Name,
			Address:       account.Account.Data.Parsed.Info.Mint,
			Pool:          pool,
			Description:   data.Result.Content.Metadata.Description,
			Image:         data.Result.Content.Links.Image,
			Amount:        account.Account.Data.Parsed.Info.TokenAmount.UIAmount,
			Price:         f,
			Pnl:           0,
			Invested:      0,
			Value:         account.Account.Data.Parsed.Info.TokenAmount.UIAmount * f,
			HistoryPrices: ohlcvs_data,
		})
		response.Tokens = tokens
		response.WalletValue = response.WalletValue + (account.Account.Data.Parsed.Info.TokenAmount.UIAmount * f)
		response.LastUpdated = time.Now().UTC().Format(time.RFC3339)

		// Send updated response with the new tokens
		if err := stream.Send(response); err != nil {
			log.Printf("error sending update: %v", err)
			return err
		}
	}

	// Get hashes for transactions
	hashes, err := solana_requests.GetTransactionHashes(req.WalletAddress)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not get transactions: %v", err)
	}

	// Create a queue of signatures.
	var queue []string
	for _, sig := range hashes {
		queue = append(queue, sig.Signature)
	}

	var transactions []*pb.Transaction

	// Process the queue.
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

		// Attempt to fetch the transaction data.
		result, err := solana_requests.QueryRPCWithRetry("getTransaction", params)
		if err != nil {

			// If the error contains a retry delay, parse it.
			var delaySeconds int
			if n, _ := fmt.Sscanf(err.Error(), "retry after %d seconds", &delaySeconds); n == 1 {
				log.Info("Rate limited. Retrying after delay", "delaySeconds", delaySeconds, "signature", signature)
				time.Sleep(time.Duration(delaySeconds) * time.Second)
			} else {
				// For other errors, log and briefly wait before requeuing.
				log.Error("Error fetching transaction", "signature", signature, "error", err)
				time.Sleep(1 * time.Second)
			}
			// Requeue the signature for a retry.
			queue = append(queue, signature)
			continue
		}

		// Pre-check the JSON response for an error field.
		var tmp map[string]interface{}
		if err = json.Unmarshal([]byte(result), &tmp); err != nil {
			log.Error("Error checking JSON for error field", "signature", signature, "error", err)
			queue = append(queue, signature)
			continue
		}
		if _, found := tmp["err"]; found {
			log.Error("RPC returned an error in response", "signature", signature, "response", tmp["err"])
			queue = append(queue, signature)
			continue
		}

		var txResponse pb.Transaction
		opts := protojson.UnmarshalOptions{
			DiscardUnknown: true,
		}
		if err = opts.Unmarshal([]byte(result), &txResponse); err != nil {
			log.Error("Error unmarshalling transaction", "signature", signature, "Stack", err)
			queue = append(queue, signature)
			continue
		}

		transactions = append(transactions, &txResponse)
		response.Transactions = transactions
		response.LastUpdated = time.Now().UTC().Format(time.RFC3339)

		// Send updated response with the new transactions
		if err := stream.Send(response); err != nil {
			log.Printf("error sending update: %v", err)
		}
	}

	// Additional steps (history price data, pnl calculations, etc.) can be added here.
	// Send final update
	if err := stream.Send(response); err != nil {
		log.Printf("error sending update: %v", err)
		return err
	}
	return nil
}

func main() {
	// Listen on port 50051
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
