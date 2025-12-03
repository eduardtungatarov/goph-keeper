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

func (h *Handler) Example(ctx context.Context, req *contracts.ExampleRequest) (*contracts.ExampleResponse, error) {
	return &contracts.ExampleResponse{}, nil
}
