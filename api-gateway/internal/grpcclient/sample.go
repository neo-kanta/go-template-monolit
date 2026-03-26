package grpcclient

import (
	"context"
	"log"
	"time"

	samplev1 "go-transfer-agent/common/gen/sample/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SampleClient wraps the generated gRPC client for the Sample service.
type SampleClient struct {
	conn   *grpc.ClientConn
	Client samplev1.SampleServiceClient
}

// NewSampleClient creates a connection to the Sample gRPC service.
func NewSampleClient(addr string) *SampleClient {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("failed to connect to Sample service at %s: %v", addr, err)
	}

	log.Printf("✅ Connected to Sample gRPC service at %s", addr)
	return &SampleClient{
		conn:   conn,
		Client: samplev1.NewSampleServiceClient(conn),
	}
}

// Close releases the gRPC connection.
func (c *SampleClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
