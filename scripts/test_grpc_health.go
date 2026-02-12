package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	addr := "localhost:50051"
	log.Printf("🔌 Connecting to gRPC server at %s...", addr)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer conn.Close()

	client := healthpb.NewHealthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		log.Fatalf("❌ Health check failed: %v", err)
	}

	status := resp.GetStatus()
	if status == healthpb.HealthCheckResponse_SERVING {
		fmt.Println("✅ Health Check: SERVING")
	} else {
		fmt.Printf("⚠️ Health Check: %s\n", status)
	}
}
