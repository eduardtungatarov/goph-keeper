package client

import (
	"fmt"
	"os"

	"github.com/eduardtungatarov/goph-keeper/internal/client/config"
	"github.com/eduardtungatarov/goph-keeper/internal/client/grpcclient"
	"github.com/eduardtungatarov/goph-keeper/internal/client/handler"
	"github.com/eduardtungatarov/goph-keeper/internal/logger"
)

func main() {
	// Логер.
	log, err := logger.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v", err)
		os.Exit(1)
	}

	// Инициализируем конфиг.
	cfg := config.Load()

	// Инициализируем grpc клиент.
	c, err := grpcclient.New(cfg)
	if err != nil {
		log.Fatalf("Failed to init grpc client: %v", err)
	}
	defer c.Close()

	// Обработчики команд.
	_ = handler.New(c)
}
