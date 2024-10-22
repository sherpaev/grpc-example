// cmd/server/main.go
package main

import (
	"grpc-example/internal/server"
	"log"
)

func main() {
	srv := server.NewServer(":50051")
	log.Println("Starting gRPC server...")
	if err := srv.Start(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
