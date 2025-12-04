up:
	docker-compose up -d

server:
	go run cmd/server/main.go

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