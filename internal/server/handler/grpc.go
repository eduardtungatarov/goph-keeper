package handler

import (
	"context"

	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"
)

type Handler struct {
	contracts.UnimplementedKeeperServiceServer
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Create(ctx context.Context, req *contracts.CreateRequest) (*contracts.CreateResponse, error) {
	return &contracts.CreateResponse{}, nil
}
