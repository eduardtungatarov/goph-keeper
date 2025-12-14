package client

import (
	"fmt"
	"os"

	"github.com/eduardtungatarov/goph-keeper/internal/client/handler"

	"github.com/eduardtungatarov/goph-keeper/internal/client/command"

	"github.com/eduardtungatarov/goph-keeper/internal/client/config"
	"github.com/eduardtungatarov/goph-keeper/internal/client/grpcclient"
	"github.com/eduardtungatarov/goph-keeper/internal/logger"
	"github.com/spf13/cobra"
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
	client, err := grpcclient.New(cfg)
	if err != nil {
		log.Fatalf("Failed to init grpc client: %v", err)
	}
	defer client.Close()

	// Обработчики команд.
	h := handler.New(client)

	// Команды.
	c := command.New(h)

	// Создаем root команду
	rootCmd := &cobra.Command{
		Use:   "keeper",
		Short: "Goph Keeper CLI client",
	}

	// Добавляем команды
	rootCmd.AddCommand(
		c.LoginCmd(),
		c.RegisterCmd(),
		c.CreateCmd(),
		c.ReadCmd(),
		c.DeleteCmd(),
		c.ListCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}
