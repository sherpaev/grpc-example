package client

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc-example/proto/service"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.ExampleServiceClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewExampleServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// Унарный вызов
func (c *Client) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return c.client.GetUser(ctx, req)
}

// Серверный стриминг
func (c *Client) GetPriceUpdates(ctx context.Context, req *pb.PriceRequest) (pb.ExampleService_GetPriceUpdatesClient, error) {
	return c.client.GetPriceUpdates(ctx, req)
}

// Клиентский стриминг
func (c *Client) UploadFile(ctx context.Context) (pb.ExampleService_UploadFileClient, error) {
	return c.client.UploadFile(ctx)
}

// Двунаправленный стриминг
func (c *Client) Chat(ctx context.Context) (pb.ExampleService_ChatClient, error) {
	return c.client.Chat(ctx)
}
