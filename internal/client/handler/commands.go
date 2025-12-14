package handler

import "github.com/eduardtungatarov/goph-keeper/internal/client/grpcclient"

type Handler struct {
	Client *grpcclient.Client
}

func New(client *grpcclient.Client) *Handler {
	return &Handler{
		Client: client,
	}
}
