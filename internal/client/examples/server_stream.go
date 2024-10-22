package examples

import (
	"context"
	"grpc-example/internal/client"
	pb "grpc-example/proto/service"
	"io"
	"log"
	"time"
)

func RunServerStreamExample(c *client.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := c.GetPriceUpdates(ctx, &pb.PriceRequest{Symbol: "BTC"})
	if err != nil {
		log.Fatalf("could not get price updates: %v", err)
	}

	for {
		price, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive price: %v", err)
		}
		log.Printf("Price update: %v", price)
	}
}
