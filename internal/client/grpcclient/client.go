package grpcclient

import (
	"github.com/eduardtungatarov/goph-keeper/internal/client/config"
	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	cfg  *config.Config
	conn *grpc.ClientConn
	c    contracts.KeeperServiceClient
}

func New(cfg *config.Config) (*Client, error) {
	client := &Client{
		cfg: cfg,
	}

	conn, err := grpc.NewClient(cfg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client.c = contracts.NewKeeperServiceClient(conn)

	return client, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
