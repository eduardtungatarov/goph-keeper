package handler

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/metadata"

	"github.com/eduardtungatarov/goph-keeper/internal/client/token"

	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"
)

// Handler обработчик команд.
type Handler struct {
	Client contracts.KeeperServiceClient
}

// New конструктор обработчиков команд.
func New(client contracts.KeeperServiceClient) *Handler {
	return &Handler{
		Client: client,
	}
}

// HandleLogin обработчик входа.
func (h *Handler) HandleLogin(ctx context.Context, login, password string) {
	resp, err := h.Client.Login(ctx, &contracts.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		log.Printf("Login failed: %v", err)
		return
	}
	fmt.Println("Login successful!")
	if err := token.Save(resp.Token); err != nil {
		log.Printf("Failed to save token: %v", err)
	} else {
		fmt.Println("Token saved to ~/.keeper/token.json")
	}
}

// HandleRegister обработчик регистрации.
func (h *Handler) HandleRegister(ctx context.Context, login, password string) {
	resp, err := h.Client.Register(ctx, &contracts.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		log.Printf("Register failed: %v", err)
		return
	}
	fmt.Printf("Register successful! Token: %s\n", resp.Token)
	if err := token.Save(resp.Token); err != nil {
		log.Printf("Failed to save token: %v", err)
	} else {
		fmt.Println("Token saved to ~/.keeper/token.json")
	}
}

// HandleCreate обработчик создания данных.
func (h *Handler) HandleCreate(ctx context.Context, dataTypeStr, data, title string) {
	ctx = h.withAuth(ctx)

	dataType, ok := h.parseDataType(dataTypeStr)
	if !ok {
		log.Fatalf("Invalid type: %s. Use: pwd, card, binary", dataTypeStr)
	}

	resp, err := h.Client.Create(ctx, &contracts.CreateRequest{
		Type:  dataType,
		Data:  []byte(data),
		Title: title,
	})
	if err != nil {
		log.Printf("Create failed: %v", err)
		return
	}

	if resp.Success {
		fmt.Printf("Item '%s' created successfully\n", title)
	} else {
		fmt.Println("Create failed")
	}
}

// HandleRead обработчик чтения данных.
func (h *Handler) HandleRead(ctx context.Context, id int64) {
	ctx = h.withAuth(ctx)

	resp, err := h.Client.Read(ctx, &contracts.ReadRequest{
		Id: id,
	})
	if err != nil {
		log.Printf("Read failed: %v", err)
		return
	}

	fmt.Printf("Type: %s\n", resp.Type.String())
	fmt.Printf("Data: %s\n", string(resp.Data))
}

// HandleDelete обработчик удаления данных.
func (h *Handler) HandleDelete(ctx context.Context, id int64) {
	ctx = h.withAuth(ctx)

	resp, err := h.Client.Delete(ctx, &contracts.DeleteRequest{
		Id: id,
	})
	if err != nil {
		log.Printf("Delete failed: %v", err)
		return
	}

	if resp.Success {
		fmt.Printf("Item %d deleted successfully\n", id)
	} else {
		fmt.Printf("Delete item %d failed\n", id)
	}
}

// HandleList обработчик получения данных пользователя.
func (h *Handler) HandleList(ctx context.Context) {
	ctx = h.withAuth(ctx)

	resp, err := h.Client.List(ctx, &contracts.ListRequest{})
	if err != nil {
		log.Printf("List failed: %v", err)
		return
	}

	fmt.Println("📋 Your items:")
	if len(resp.DataList) == 0 {
		fmt.Println("   (empty)")
		return
	}

	for _, item := range resp.DataList {
		fmt.Printf(" ID: %d | %s | %s\n",
			item.Id, item.Type.String(), item.Title)
	}
}

func (h *Handler) withAuth(ctx context.Context) context.Context {
	t, err := token.Load()
	if err != nil {
		log.Printf("No token, login first")
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, "token", t)
}

func (h *Handler) parseDataType(s string) (contracts.DataType, bool) {
	switch s {
	case "pwd", "PWD":
		return contracts.DataType_PWD, true
	case "card", "CARD":
		return contracts.DataType_CARD, true
	case "binary", "BINARY":
		return contracts.DataType_BINARY, true
	default:
		return contracts.DataType_UNSPECIFIED, false
	}
}
