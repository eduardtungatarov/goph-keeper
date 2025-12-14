package grpcclient

import (
	"github.com/eduardtungatarov/goph-keeper/internal/client/config"
	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Cfg  *config.Config
	Conn *grpc.ClientConn
	C    contracts.KeeperServiceClient
}

func New(cfg *config.Config) (*Client, error) {
	client := &Client{
		Cfg: cfg,
	}

	conn, err := grpc.NewClient(cfg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client.Conn = conn
	client.C = contracts.NewKeeperServiceClient(conn)

	return client, nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
