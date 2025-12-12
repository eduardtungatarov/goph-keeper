package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/eduardtungatarov/goph-keeper/internal/server/service/security"

	"github.com/eduardtungatarov/goph-keeper/internal/server/interceptor"

	"golang.org/x/sync/errgroup"

	"github.com/eduardtungatarov/goph-keeper/internal/config"
	grpcServer "github.com/eduardtungatarov/goph-keeper/internal/server/grpc"
	"github.com/eduardtungatarov/goph-keeper/internal/server/handler"
	dataRepository "github.com/eduardtungatarov/goph-keeper/internal/server/repository/data"
	userRepository "github.com/eduardtungatarov/goph-keeper/internal/server/repository/user"
	authService "github.com/eduardtungatarov/goph-keeper/internal/server/service/auth"
	dataService "github.com/eduardtungatarov/goph-keeper/internal/server/service/data"

	"github.com/eduardtungatarov/goph-keeper/internal/logger"
	"github.com/pressly/goose"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Логер.
	log, err := logger.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v", err)
		os.Exit(1)
	}

	// Инициализируем конфиг.
	cfg := config.Load()

	// Получаем экземпляр БД.
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// Применяем миграции.
	err = goose.Up(db, "migrations")
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	grp, ctx := errgroup.WithContext(ctx)

	// Собираем зависимости.
	authRepo := userRepository.New(db)
	authSrv := authService.New(cfg.JWTSecretKey, authRepo)
	secSrv := security.New([]byte(cfg.DataEncryptKey))
	dataRepo := dataRepository.New(db)
	dataSrv := dataService.New(dataRepo, secSrv)

	// Инициализируем и запускаем grpc сервер.
	grp.Go(func() error {
		h := handler.New(
			authSrv,
			dataSrv,
		)
		i := interceptor.New(authSrv)
		s := grpcServer.New(log, h, cfg, i)
		err := s.Run(ctx)
		if err != nil {
			return fmt.Errorf("gRPC server Run error: %w", err)
		}
		return nil
	})

	if err := grp.Wait(); err != nil {
		log.Errorf("failed service reason: %v", err)
	}
	log.Info("service has been stopped")
}
