package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/eduardtungatarov/goph-keeper/internal/server/handler"

	"go.uber.org/zap"

	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	log *zap.SugaredLogger
	h   *handler.Handler
}

func New(log *zap.SugaredLogger, h *handler.Handler) *server {
	return &server{
		log: log,
		h:   h,
	}
}

func (s *server) Run(ctx context.Context) error {
	// Настраиваем.
	grpcServer := grpc.NewServer()
	contracts.RegisterKeeperServiceServer(grpcServer, nil)
	reflection.Register(grpcServer)

	// Открываем порт.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Запускаем.
	serverErr := make(chan error, 1)
	go func() {
		s.log.Info("gRPC server starting ...")
		if err := grpcServer.Serve(lis); err != nil {
			serverErr <- fmt.Errorf("serve failed: %w", err)
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		s.log.Info("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		// Ждем пока grpcServer.Serve() завершится
		<-serverErr
		return nil
	}
}
