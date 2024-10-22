package examples

import (
	"context"
	"fmt"
	"grpc-example/internal/client"
	pb "grpc-example/proto/service"
	"io"
	"log"
	"time"
)

func RunBiDirectionalExample(c *client.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := c.Chat(ctx)
	if err != nil {
		log.Fatalf("could not start chat: %v", err)
	}

	waitc := make(chan struct{})

	// Горутина для получения сообщений
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("failed to receive message: %v", err)
			}
			log.Printf("Received: %v", in)
		}
	}()

	// Отправка сообщений
	for i := 0; i < 3; i++ {
		msg := &pb.ChatMessage{
			UserId:    "client-1",
			Content:   fmt.Sprintf("Message %d", i),
			Timestamp: time.Now().String(),
		}
		if err := stream.Send(msg); err != nil {
			log.Fatalf("failed to send message: %v", err)
		}
		time.Sleep(time.Second)
	}

	stream.CloseSend()
	<-waitc
}
