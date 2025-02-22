package main

import (
	"errors"
	"math/rand"
	"net"
	"time"

	"solana/requests"

	"github.com/charmbracelet/log"
	"github.com/mr-tron/base58"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "solana/generated"
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

	// Get the base wallet
	wallet, err := requests.RequestAccountInfo(req.WalletAddress)
	if err != nil {
		return status.Errorf(codes.NotFound, "invalid wallet address: %v", err)
	}

	//return first response

	// if err := stream.Send(update); err != nil {
	// 	log.Printf("error sending update: %v", err)
	// 	return err
	// }
	// Get the current solana price
	solana_price, err := requests.GetSolanaPrice()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not get solana price: %v", err)
	}

	//return the updated wallet with solana value and amount

	// get the tokens in the wallet
	accounts, err := requests.RequestTokenAccounts(req.WalletAddress)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not get tokens: %v", err)
	}

	var addresses []string
	for _, token := range accounts {
		addresses = append(addresses, token.Account.Data.Parsed.Info.Mint)
	}
	// get the current token prices
	current_token_prices, err := requests.GetCoinGeckoTokenPrices(addresses)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not get token prices: %v", err)
	}

	//return the updated wallet with all the tokens
	var tokens []pb.Token
	for _, account := range accounts.Result.Value {
		data := requests.GetTokenMetadata(account.Account.Data.Parsed.Info.Mint)
		pool, _ := requests.GetTokenPools(account.Account.Data.Parsed.Info.Mint)
		tokens = append(tokens, pb.Token{
			Name:        data.Result.Content.Metadata.Name,
			Address:     account.Account.Data.Parsed.Info.Mint,
			Pool:        pool,
			Description: data.Result.Content.Metadata.Description,
			Image:       data.Result.Content.Metadata.Image,
			Amount:      account.Account.Data.Parsed.Info.TokenAmount.UIAmount,
			Price:       current_token_prices[account.Account.Data.Parsed.Info.Mint],
			Pnl:         0,
			Invested:    0,
			Value:       0,
			HistoryPrices: []float64{
				0,
			},
		})

		// send updated response with the new tokens
	}

	// get the transactions in the wallet
	transactions, err := requests.RequestTransactions(req.WalletAddress)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not get transactions: %v", err)
	}

	//loop over transaction hashes get the transaction info and send updated wallet with the transaction attached

	// get history price data for each token that is in the transactions

	// send updated wallet with history price data

	// calculate pnl for each token

	// send for each calculation a update

	//IMPORTANT: for each update that we do update the wallet value and total pnl along the way

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
