package examples

import (
	"context"
	"grpc-example/internal/client"
	pb "grpc-example/proto/service"
	"log"
	"time"
)

func RunUnaryExample(c *client.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.GetUser(ctx, &pb.UserRequest{UserId: 123})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("User: %v", resp)
}
