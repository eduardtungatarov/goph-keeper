package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"

	"google.golang.org/grpc/credentials"

	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"

	"google.golang.org/grpc"

	"github.com/eduardtungatarov/goph-keeper/internal/client/handler"

	"github.com/eduardtungatarov/goph-keeper/internal/client/command"

	"github.com/eduardtungatarov/goph-keeper/internal/client/config"
	"github.com/spf13/cobra"

	_ "github.com/joho/godotenv/autoload"
)

var (
	version   = "dev"
	buildDate = "unknown"
)

func main() {
	// Инициализируем конфиг.
	cfg := config.Load()

	// Инициализируем grpc клиент.
	certBytes, _ := base64.StdEncoding.DecodeString(cfg.ServerCert)
	serverCert, _ := tls.X509KeyPair(certBytes, []byte{})
	conn, err := grpc.NewClient(
		cfg.ServerAddress,
		grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{
				Certificates:       []tls.Certificate{serverCert},
				InsecureSkipVerify: true,
			}),
		),
	)
	if err != nil {
		log.Fatalf("Failed to init grpc client: %v", err)
	}
	defer conn.Close()
	client := contracts.NewKeeperServiceClient(conn)

	// Обработчики команд.
	h := handler.New(client)

	// Команды.
	c := command.New(h)

	// Создаем root команду
	rootCmd := &cobra.Command{
		Use:   "keeper",
		Short: "Goph Keeper CLI client",
	}

	rootCmd.Version = fmt.Sprintf("%s (built %s)", version, buildDate)
	rootCmd.SetVersionTemplate("{{.Version}}\n")

	// Добавляем команды
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "version",
			Short: "Show version info",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Printf("Keeper %s (built %s)\n", version, buildDate)
			},
		},
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
