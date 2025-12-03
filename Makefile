contract-generate:
	protoc \
      --go_out=./ \
      --go-grpc_out=./ \
      contracts/server/server.proto