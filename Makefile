up:
	docker-compose up -d

server:
	go run cmd/server/main.go

client:
	go run cmd/client/main.go

contract-generate:
	protoc \
      --go_out=./ \
      --go-grpc_out=./ \
      contracts/server/server.proto

make gen-sql:
	sqlc generate

migrate-create:
	goose -dir migrations create ${name} sql

vet:
	go vet ./... && goimports -w . && gofumpt -w -extra . && staticcheck ./...

test:
	go test ./...

dep:
	go mod tidy

build:
	go build -o keeper ./cmd/client

build-all:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=v1.0.0 -X main.buildDate=`date -u +%Y-%m-%dT%H:%M:%SZ`" -o bin/keeper-linux-amd64 ./cmd/client
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=v1.0.0 -X main.buildDate=`date -u +%Y-%m-%dT%H:%M:%SZ`" -o bin/keeper-windows-amd64.exe ./cmd/client
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=v1.0.0 -X main.buildDate=`date -u +%Y-%m-%dT%H:%M:%SZ`" -o bin/keeper-darwin-amd64 ./cmd/client
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=v1.0.0 -X main.buildDate=`date -u +%Y-%m-%dT%H:%M:%SZ`" -o bin/keeper-darwin-arm64 ./cmd/client