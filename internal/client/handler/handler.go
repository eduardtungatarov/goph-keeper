package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"

	"github.com/eduardtungatarov/goph-keeper/internal/client/grpcclient"
)

type Handler struct {
	Client *grpcclient.Client
}

func New(client *grpcclient.Client) *Handler {
	return &Handler{
		Client: client,
	}
}

func (h *Handler) HandleLogin(login, password string) {
	resp, err := h.Client.C.Login(context.Background(), &contracts.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Printf("✅ Login successful. Token: %s\n", resp.Token)
}

func (h *Handler) HandleRegister(login, password string) {
	//
}

func (h *Handler) HandleCreate(dataTypeStr, data, title string) {
	//
}

func (h *Handler) HandleRead(id int64) {
	//
}

func (h *Handler) HandleDelete(id int64) {
	//
}

func (h *Handler) HandleList() {
	//
}
