package main

import (
	"context"
	"github.com/charmbracelet/log"
	"net"

	"google.golang.org/grpc"

	pb "solana/generated"
)

type server struct {
	pb.UnimplementedWalletServiceServer
}

// GetWalletInfo is the main gRPC method that performs all operations.
func (s *server) GetWalletInfo(ctx context.Context, req *pb.WalletRequest) (*pb.WalletResponse, error) {
	walletAddr := req.WalletAddress
	var tokens []*pb.Token
	totalValue := 0.0
	log.Infof("GetWalletInfo: %v", walletAddr)
	solBalance, err := GetBalance(walletAddr)
	if err != nil {
		return nil, err
	}
	toks, err := GetTokenAccounts(walletAddr)
	if err != nil {
		return nil, err
	}
	solPrice, err := GetSolPrice()
	if err != nil {
		log.Errorf("Error fetching sol price: %v", err)
		solPrice = 0
	}
	solBalanceUsd := solBalance * solPrice
	totalValue += solBalanceUsd
	for _, token := range toks {
		tokenPrice, err := GetTokenPrice(token.Mint)
		if err != nil {
			log.Printf("Error fetching price for token mint %s: %v", token.Mint, err)
			tokenPrice = 0
		}
		totalValue += tokenPrice * token.Amount
		t := &pb.Token{
			TokenSymbol:   token.TokenAccount,
			TokenAddress:  token.Mint,
			TokenBalance:  token.Amount,
			UsdBalance:    token.Amount * tokenPrice,
			CurrentPrice:  tokenPrice,
			TotalEntry:    0,
			Pnl:           0,
			Transactions:  nil,
			HistoryPrices: nil,
		}
		tokens = append(tokens, t)

	}

	solB := &pb.Token{
		TokenSymbol:   "SOL",
		TokenAddress:  "So11111111111111111111111111111111111111112",
		TokenBalance:  solBalance,
		UsdBalance:    solBalanceUsd,
		CurrentPrice:  solPrice,
		TotalEntry:    0,
		Pnl:           0,
		Transactions:  nil,
		HistoryPrices: nil,
	}
	tokens = append(tokens, solB)
	return &pb.WalletResponse{
		Tokens:             tokens,
		WalletTotalBalance: float32(totalValue),
		WalletTotalPnl:     0,
	}, nil
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
