# gRPC Example Service

Пример реализации gRPC сервиса с поддержкой всех типов взаимодействия (unary, server streaming, client streaming, bidirectional streaming) и middleware для логирования.

## Структура проекта

```
grpc-example/
├── Makefile                  # Команды для сборки и запуска
├── README.md                 # Документация проекта
├── go.mod                    # Определение модуля и зависимостей
├── proto/                    # Директория с proto файлами
│   └── service/
│       └── service.proto     # Определения сервисов и сообщений
├── pkg/                      # Общий код для переиспользования
│   └── models/
│       └── message.go
├── internal/                 # Внутренняя логика приложения
│   ├── server/              # Серверная часть
│   │   └── server.go        # Реализация сервера
│   ├── client/              # Клиентская часть
│   │   ├── client.go
│   │   └── examples/        # Примеры использования
│   │       ├── unary.go
│   │       ├── server_stream.go
│   │       ├── client_stream.go
│   │       └── bidirectional.go
│   └── middleware/          # Middleware компоненты
│       ├── logging.go       # Логирование запросов
│       └── zap_logger.go    # Реализация с помощью zap
└── cmd/                     # Точки входа в приложение
    ├── server/
    │   └── main.go         # Запуск сервера
    └── client/
        └── main.go         # Запуск клиента
```

## Предварительные требования

1. Go 1.21 или выше
2. Protocol Buffers компилятор
3. Go плагины для protoc

### Установка зависимостей

1. Установка protoc:
```bash
# macOS
brew install protobuf

# Ubuntu/Debian
sudo apt install -y protobuf-compiler

# Windows (через chocolatey)
choco install protoc
```

2. Установка Go плагинов для protoc:
```bash
# Плагин для генерации Go кода
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Плагин для генерации gRPC кода
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

3. Добавление GOPATH в PATH:
```bash
# Добавьте в ~/.bashrc или ~/.zshrc
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Установка и запуск

1. Клонирование репозитория:
```bash
git clone https://github.com/your-username/grpc-example.git
cd grpc-example
```

2. Установка зависимостей:
```bash
go mod tidy
```

3. Генерация кода из proto файлов:
```bash
make proto
```

4. Запуск сервера и клиента:
```bash
# Терминал 1: Запуск сервера
make run-server

# Терминал 2: Запуск клиента
make run-client
```

## Примеры использования

### Унарный вызов
```go
// Клиент
resp, err := client.GetUser(ctx, &pb.UserRequest{UserId: 123})
if err != nil {
    log.Fatalf("Error: %v", err)
}
log.Printf("User: %v", resp)

// Сервер
func (s *Server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
    return &pb.UserResponse{
        UserId: req.UserId,
        Name:   fmt.Sprintf("User %d", req.UserId),
        Email:  fmt.Sprintf("user%d@example.com", req.UserId),
    }, nil
}
```

### Серверный стриминг
```go
// Клиент
stream, err := client.GetPriceUpdates(ctx, &pb.PriceRequest{Symbol: "BTC"})
if err != nil {
    log.Fatalf("Error: %v", err)
}
for {
    price, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    log.Printf("Price update: %v", price)
}

// Сервер
func (s *Server) GetPriceUpdates(req *pb.PriceRequest, stream pb.ExampleService_GetPriceUpdatesServer) error {
    for i := 0; i < 5; i++ {
        if err := stream.Send(&pb.PriceResponse{
            Symbol:    req.Symbol,
            Price:     100.0 + float64(i),
            Timestamp: time.Now().String(),
        }); err != nil {
            return err
        }
        time.Sleep(time.Second)
    }
    return nil
}
```

### Клиентский стриминг
```go
// Клиент
stream, err := client.UploadFile(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}
for i := 0; i < 3; i++ {
    if err := stream.Send(&pb.FileChunk{
        Content: []byte(fmt.Sprintf("chunk %d", i)),
    }); err != nil {
        log.Fatalf("Error: %v", err)
    }
}
status, err := stream.CloseAndRecv()
if err != nil {
    log.Fatalf("Error: %v", err)
}
log.Printf("Upload status: %v", status)

// Сервер
func (s *Server) UploadFile(stream pb.ExampleService_UploadFileServer) error {
    var totalSize int
    for {
        chunk, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&pb.UploadStatus{
                Success: true,
                Message: fmt.Sprintf("Upload complete. Size: %d bytes", totalSize),
            })
        }
        if err != nil {
            return err
        }
        totalSize += len(chunk.Content)
    }
}
```

### Двунаправленный стриминг
```go
// Клиент
stream, err := client.Chat(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

waitc := make(chan struct{})
go func() {
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            close(waitc)
            return
        }
        if err != nil {
            log.Fatalf("Error: %v", err)
        }
        log.Printf("Received: %v", msg)
    }
}()

for i := 0; i < 3; i++ {
    if err := stream.Send(&pb.ChatMessage{
        UserId:    "client-1",
        Content:   fmt.Sprintf("Message %d", i),
        Timestamp: time.Now().String(),
    }); err != nil {
        log.Fatalf("Error: %v", err)
    }
    time.Sleep(time.Second)
}

stream.CloseSend()
<-waitc

// Сервер
func (s *Server) Chat(stream pb.ExampleService_ChatServer) error {
    for {
        in, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }

        if err := stream.Send(&pb.ChatMessage{
            UserId:    "server",
            Content:   fmt.Sprintf("Received: %s", in.Content),
            Timestamp: time.Now().String(),
        }); err != nil {
            return err
        }
    }
}
```

## Middleware для логирования

### Базовое использование
```go
// Создание сервера с middleware
logger := &middleware.DefaultLogger{}
loggingInterceptor := middleware.NewLoggingInterceptor(logger)

grpcServer := grpc.NewServer(
    grpc.UnaryInterceptor(loggingInterceptor.UnaryServerInterceptor()),
    grpc.StreamInterceptor(loggingInterceptor.StreamServerInterceptor()),
)
```

### Использование с Zap логгером
```go
// Создание zap логгера
zapLogger, err := middleware.NewZapLogger()
if err != nil {
    log.Fatal(err)
}

// Создание интерцептора с zap логгером
loggingInterceptor := middleware.NewLoggingInterceptor(zapLogger)

grpcServer := grpc.NewServer(
    grpc.UnaryInterceptor(loggingInterceptor.UnaryServerInterceptor()),
    grpc.StreamInterceptor(loggingInterceptor.StreamServerInterceptor()),
)
```

## Команды Makefile

```bash
# Генерация кода из proto файлов
make proto

# Сборка сервера и клиента
make build

# Запуск сервера
make run-server

# Запуск клиента
make run-client

# Очистка сгенерированных файлов
make clean

# Обновление зависимостей
make tidy
```

## Логирование

Middleware для логирования предоставляет следующую информацию:
- Время начала и длительность запроса
- Метод gRPC
- Метаданные запроса
- IP адрес клиента
- Статус выполнения
- Количество отправленных и полученных сообщений (для стримов)
- Ошибки, если они возникли

### Пример вывода логов

Базовый логгер:
```
INFO: Starting unary call: GetUser, metadata: map[peer_address:127.0.0.1:52431]
INFO: Unary call successful: GetUser, duration: 1.234ms
```

Zap логгер (JSON формат):
```json
{
    "level": "info",
    "ts": 1635789654.789,
    "caller": "middleware/logging.go:123",
    "msg": "Starting unary call",
    "method": "GetUser",
    "metadata": {
        "peer_address": "127.0.0.1:52431"
    }
}
```

## Лицензия

MIT
