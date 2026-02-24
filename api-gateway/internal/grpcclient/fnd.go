package grpcclient

import (
	"context"
	"log"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// FNDClient wraps the generated gRPC client for the FND service.
type FNDClient struct {
	conn   *grpc.ClientConn
	Client fndv1.FNDServiceClient
}

// NewFNDClient creates a connection to the FND gRPC service.
func NewFNDClient(addr string) *FNDClient {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("failed to connect to FND service at %s: %v", addr, err)
	}

	log.Printf("✅ Connected to FND gRPC service at %s", addr)
	return &FNDClient{
		conn:   conn,
		Client: fndv1.NewFNDServiceClient(conn),
	}
}

// Close releases the gRPC connection.
func (c *FNDClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
