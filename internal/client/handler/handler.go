package handler

import (
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
	//
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
