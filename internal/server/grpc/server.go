package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/eduardtungatarov/goph-keeper/internal/server/interceptor"

	"github.com/eduardtungatarov/goph-keeper/internal/config"

	"github.com/eduardtungatarov/goph-keeper/internal/server/handler"

	"go.uber.org/zap"

	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	log *zap.SugaredLogger
	h   *handler.Handler
	cfg *config.Config
	i   *interceptor.Interceptor
}

func New(log *zap.SugaredLogger, h *handler.Handler, cfg *config.Config, i *interceptor.Interceptor) *server {
	return &server{
		log: log,
		h:   h,
		cfg: cfg,
		i:   i,
	}
}

func (s *server) Run(ctx context.Context) error {
	// Настраиваем.
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(s.i.AuthInterceptor))
	contracts.RegisterKeeperServiceServer(grpcServer, s.h)
	reflection.Register(grpcServer)

	// Открываем порт.
	lis, err := net.Listen("tcp", ":"+s.cfg.GRPCPort)
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
