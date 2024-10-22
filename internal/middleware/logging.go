package middleware

import (
	"context"
	"fmt"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor содержит опции для логирования
type LoggingInterceptor struct {
	logger Logger
}

// Logger интерфейс для логирования
type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

// NewLoggingInterceptor создает новый инстанс интерцептора
func NewLoggingInterceptor(logger Logger) *LoggingInterceptor {
	return &LoggingInterceptor{
		logger: logger,
	}
}

// extractMetadata извлекает метаданные из контекста
func extractMetadata(ctx context.Context) map[string]string {
	md := make(map[string]string)

	if mtd, ok := metadata.FromIncomingContext(ctx); ok {
		for k, v := range mtd {
			if len(v) > 0 {
				md[k] = v[0]
			}
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		md["peer_address"] = p.Addr.String()
	}

	return md
}

// UnaryServerInterceptor создает унарный серверный перехватчик
func (l *LoggingInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()

		// Извлекаем метаданные
		md := extractMetadata(ctx)

		// Логируем начало запроса
		l.logger.Infof("Starting unary call: %s, metadata: %v",
			path.Base(info.FullMethod),
			md,
		)

		// Выполняем запрос
		resp, err := handler(ctx, req)

		// Вычисляем длительность
		duration := time.Since(startTime)

		// Логируем результат
		if err != nil {
			st, _ := status.FromError(err)
			l.logger.Errorf("Unary call failed: %s, code: %s, error: %v, duration: %v",
				path.Base(info.FullMethod),
				st.Code(),
				err,
				duration,
			)
		} else {
			l.logger.Infof("Unary call successful: %s, duration: %v",
				path.Base(info.FullMethod),
				duration,
			)
		}

		return resp, err
	}
}

// StreamServerInterceptor создает потоковый серверный перехватчик
func (l *LoggingInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		startTime := time.Now()

		// Извлекаем метаданные
		md := extractMetadata(stream.Context())

		// Оборачиваем стрим для подсчета сообщений
		wrapped := &wrappedStream{
			ServerStream: stream,
			receivedMsgs: 0,
			sentMsgs:     0,
		}

		// Логируем начало стрима
		l.logger.Infof("Starting stream: %s, metadata: %v",
			path.Base(info.FullMethod),
			md,
		)

		// Выполняем обработку стрима
		err := handler(srv, wrapped)

		// Вычисляем длительность
		duration := time.Since(startTime)

		// Логируем результат
		if err != nil {
			st, _ := status.FromError(err)
			l.logger.Errorf("Stream failed: %s, code: %s, error: %v, duration: %v, received messages: %d, sent messages: %d",
				path.Base(info.FullMethod),
				st.Code(),
				err,
				duration,
				wrapped.receivedMsgs,
				wrapped.sentMsgs,
			)
		} else {
			l.logger.Infof("Stream successful: %s, duration: %v, received messages: %d, sent messages: %d",
				path.Base(info.FullMethod),
				duration,
				wrapped.receivedMsgs,
				wrapped.sentMsgs,
			)
		}

		return err
	}
}

// wrappedStream оборачивает grpc.ServerStream для подсчета сообщений
type wrappedStream struct {
	grpc.ServerStream
	receivedMsgs int
	sentMsgs     int
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	err := w.ServerStream.RecvMsg(m)
	if err == nil {
		w.receivedMsgs++
	}
	return err
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	err := w.ServerStream.SendMsg(m)
	if err == nil {
		w.sentMsgs++
	}
	return err
}

// DefaultLogger реализация логгера по умолчанию
type DefaultLogger struct{}

func (l *DefaultLogger) Info(args ...interface{}) {
	fmt.Printf("INFO: %v\n", args...)
}

func (l *DefaultLogger) Error(args ...interface{}) {
	fmt.Printf("ERROR: %v\n", args...)
}

func (l *DefaultLogger) Infof(template string, args ...interface{}) {
	fmt.Printf("INFO: "+template+"\n", args...)
}

func (l *DefaultLogger) Errorf(template string, args ...interface{}) {
	fmt.Printf("ERROR: "+template+"\n", args...)
}
