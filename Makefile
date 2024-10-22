.PHONY: proto build run-server run-client

proto:
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		proto/service/service.proto

build:
	go build -o bin/server cmd/server/main.go
	go build -o bin/client cmd/client/main.go

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

# Дополнительные полезные команды
clean:
	rm -rf bin/
	rm -f proto/service/*.pb.go

tidy:
	go mod tidy

# Создание необходимых директорий
init:
	mkdir -p bin
	mkdir -p proto/service
	mkdir -p internal/server
	mkdir -p internal/client
	mkdir -p cmd/server
	mkdir -p cmd/client