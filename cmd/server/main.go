package main

import (
	"context"
	"database/sql"
	"os/signal"
	"syscall"

	"github.com/eduardtungatarov/goph-keeper/internal/logger"
	"github.com/pressly/goose"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	_, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Логер.
	log, err := logger.New()
	if err != nil {
		panic(err)
	}

	// Получаем экземпляр БД.
	db, err := sql.Open("pgx", "DSN")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// Применяем миграции.
	err = goose.Up(db, "migrations")
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
}
