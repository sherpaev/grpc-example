package examples

import (
	"context"
	"fmt"
	"grpc-example/internal/client"
	pb "grpc-example/proto/service"
	"log"
	"time"
)

func RunClientStreamExample(c *client.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := c.UploadFile(ctx)
	if err != nil {
		log.Fatalf("could not upload file: %v", err)
	}

	// Simulate sending file chunks
	for i := 0; i < 3; i++ {
		chunk := &pb.FileChunk{
			Content: []byte(fmt.Sprintf("chunk %d", i)),
		}
		if err := stream.Send(chunk); err != nil {
			log.Fatalf("failed to send chunk: %v", err)
		}
	}

	status, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to receive status: %v", err)
	}
	log.Printf("Upload status: %v", status)
}
