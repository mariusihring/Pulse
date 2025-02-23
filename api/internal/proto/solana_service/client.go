package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"pulse/internal/proto/solana_service/generated"
)

func NewSolanaGRPCClient() (generated.WalletServiceClient, func(), error) {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	// Create a cleanup function to close the connection.
	cleanup := func() {
		conn.Close()
	}

	// Create the client from the generated code.
	client := generated.NewWalletServiceClient(conn)
	return client, cleanup, nil
}
