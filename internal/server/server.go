package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-example/internal/middleware"
	pb "grpc-example/proto/service"
	"net"
	"time"
)

type Server struct {
	pb.UnimplementedExampleServiceServer
	port string
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	// Создаем zap логгер
	zapLogger, err := middleware.NewZapLogger()
	if err != nil {
		return err
	}

	// Создаем интерцептор с zap логгером
	loggingInterceptor := middleware.NewLoggingInterceptor(zapLogger)

	// Создаем сервер с middleware
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor.UnaryServerInterceptor()),
		grpc.StreamInterceptor(loggingInterceptor.StreamServerInterceptor()),
	)

	pb.RegisterExampleServiceServer(grpcServer, s)

	zapLogger.Info("Starting gRPC server on ", s.port)
	return grpcServer.Serve(lis)
}

// Унарный вызов
func (s *Server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		UserId: req.UserId,
		Name:   fmt.Sprintf("User %d", req.UserId),
		Email:  fmt.Sprintf("user%d@example.com", req.UserId),
	}, nil
}

// Серверный стриминг
func (s *Server) GetPriceUpdates(req *pb.PriceRequest, stream pb.ExampleService_GetPriceUpdatesServer) error {
	for i := 0; i < 5; i++ {
		price := &pb.PriceResponse{
			Symbol:    req.Symbol,
			Price:     100.0 + float64(i),
			Timestamp: time.Now().String(),
		}
		if err := stream.Send(price); err != nil {
			return fmt.Errorf("failed to send price update: %v", err)
		}
		time.Sleep(time.Second)
	}
	return nil
}

// Клиентский стриминг
func (s *Server) UploadFile(stream pb.ExampleService_UploadFileServer) error {
	var totalSize int
	for {
		chunk, err := stream.Recv()
		if err == nil {
			totalSize += len(chunk.Content)
			continue
		}
		if err.Error() == "EOF" {
			return stream.SendAndClose(&pb.UploadStatus{
				Success: true,
				Message: fmt.Sprintf("Upload complete. Total size: %d bytes", totalSize),
			})
		}
		return fmt.Errorf("failed to receive chunk: %v", err)
	}
}

// Двунаправленный стриминг
func (s *Server) Chat(stream pb.ExampleService_ChatServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				return nil
			}
			return fmt.Errorf("failed to receive message: %v", err)
		}

		response := &pb.ChatMessage{
			UserId:    "server",
			Content:   fmt.Sprintf("Server received: %s", in.Content),
			Timestamp: time.Now().String(),
		}

		if err := stream.Send(response); err != nil {
			return fmt.Errorf("failed to send response: %v", err)
		}
	}
}
