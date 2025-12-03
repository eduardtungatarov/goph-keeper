package main

import (
	"context"
	"database/sql"
	"os/signal"
	"syscall"

	"github.com/eduardtungatarov/goph-keeper/internal/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	_, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Логер.
	log, err := logger.MakeLogger()
	if err != nil {
		panic(err)
	}

	// Подключаемся к базе данных.
	db, err := sql.Open("pgx", "DSN")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()
}
