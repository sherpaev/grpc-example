// cmd/client/main.go
package main

import (
	"grpc-example/internal/client"
	"grpc-example/internal/client/examples"
	"log"
)

func main() {
	c, err := client.NewClient("localhost:50051")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	log.Println("Running unary example...")
	examples.RunUnaryExample(c)

	log.Println("Running server streaming example...")
	examples.RunServerStreamExample(c)

	log.Println("Running client streaming example...")
	examples.RunClientStreamExample(c)

	log.Println("Running bidirectional streaming example...")
	examples.RunBiDirectionalExample(c)
}
