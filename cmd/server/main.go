package main

import (
	"context"
	"database/sql"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	_, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Подключаемся к базе данных.
	db, err := sql.Open("pgx", "DSN")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()
}
