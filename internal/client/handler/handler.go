package handler

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"
)

// Result результат выполнения хендлера
type Result struct {
	Buffer *bytes.Buffer // накопитель всех сообщений
}

// NewBuffer создает новый буфер сообщений
func NewBuffer() *bytes.Buffer {
	return bytes.NewBufferString("")
}

// Client клиент к grpc серверу.
//
//go:generate mockery --name=Client
type Client interface {
	Login(ctx context.Context, in *contracts.LoginRequest, opts ...grpc.CallOption) (*contracts.LoginResponse, error)
	Register(ctx context.Context, in *contracts.LoginRequest, opts ...grpc.CallOption) (*contracts.LoginResponse, error)
	Create(ctx context.Context, in *contracts.CreateRequest, opts ...grpc.CallOption) (*contracts.CreateResponse, error)
	Read(ctx context.Context, in *contracts.ReadRequest, opts ...grpc.CallOption) (*contracts.ReadResponse, error)
	Delete(ctx context.Context, in *contracts.DeleteRequest, opts ...grpc.CallOption) (*contracts.DeleteResponse, error)
	List(ctx context.Context, in *contracts.ListRequest, opts ...grpc.CallOption) (*contracts.ListResponse, error)
}

// TokenStorage хранилище токена пользователя.
//
//go:generate mockery --name=TokenStorage
type TokenStorage interface {
	Save(token string) error
	Load() (string, error)
}

// Handler обработчик команд.
type Handler struct {
	Client         contracts.KeeperServiceClient
	MaxRetries     int
	RetryDelay     time.Duration
	BackoffFactor  float64
	RetryableCodes []codes.Code
	TokenStorage   TokenStorage
}

// New конструктор обработчиков команд.
func New(client contracts.KeeperServiceClient, tokenStorage TokenStorage) *Handler {
	return &Handler{
		Client:        client,
		MaxRetries:    3,
		RetryDelay:    time.Second,
		BackoffFactor: 2.0,
		RetryableCodes: []codes.Code{
			codes.Unavailable, codes.DeadlineExceeded,
			codes.ResourceExhausted, codes.Internal,
		},
		TokenStorage: tokenStorage,
	}
}

// HandleLogin обработчик входа.
func (h *Handler) HandleLogin(ctx context.Context, login, password string) *Result {
	buffer := NewBuffer()

	err := h.withRetry(ctx, func(ctx context.Context) error {
		resp, err := h.Client.Login(ctx, &contracts.LoginRequest{
			Login:    login,
			Password: password,
		})
		if err != nil {
			return err
		}

		buffer.WriteString("Login successful!\n")
		if err := h.TokenStorage.Save(resp.Token); err != nil {
			buffer.WriteString(fmt.Sprintf("Failed to save token: %v\n", err))
		} else {
			buffer.WriteString("Token saved to ~/.keeper/token.json\n")
		}
		return nil
	})
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Login failed after retries: %v\n", err))
	}

	return &Result{Buffer: buffer}
}

// HandleRegister обработчик регистрации.
func (h *Handler) HandleRegister(ctx context.Context, login, password string) *Result {
	buffer := NewBuffer()

	err := h.withRetry(ctx, func(ctx context.Context) error {
		resp, err := h.Client.Register(ctx, &contracts.LoginRequest{
			Login:    login,
			Password: password,
		})
		if err != nil {
			return err
		}

		buffer.WriteString(fmt.Sprintf("Register successful! Token: %s\n", resp.Token))
		if err := h.TokenStorage.Save(resp.Token); err != nil {
			buffer.WriteString(fmt.Sprintf("Failed to save token: %v\n", err))
		} else {
			buffer.WriteString("Token saved to ~/.keeper/token.json\n")
		}
		return nil
	})
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Register failed after retries: %v\n", err))
	}

	return &Result{Buffer: buffer}
}

// HandleCreate обработчик создания данных.
func (h *Handler) HandleCreate(ctx context.Context, dataTypeStr, data, title string) *Result {
	buffer := NewBuffer()

	dataType, ok := h.parseDataType(dataTypeStr)
	if !ok {
		buffer.WriteString(fmt.Sprintf("Invalid type: %s. Use: pwd, card, binary\n", dataTypeStr))
		return &Result{Buffer: buffer}
	}

	err := h.withRetryAuth(ctx, func(ctx context.Context) error {
		resp, err := h.Client.Create(ctx, &contracts.CreateRequest{
			Type:  dataType,
			Data:  []byte(data),
			Title: title,
		})
		if err != nil {
			return err
		}

		if resp.Success {
			buffer.WriteString(fmt.Sprintf("Item '%s' created successfully\n", title))
		} else {
			buffer.WriteString("Create failed\n")
		}
		return nil
	})
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Create failed after retries: %v\n", err))
	}

	return &Result{Buffer: buffer}
}

// HandleRead обработчик чтения данных.
func (h *Handler) HandleRead(ctx context.Context, id int64) *Result {
	buffer := NewBuffer()

	err := h.withRetryAuth(ctx, func(ctx context.Context) error {
		resp, err := h.Client.Read(ctx, &contracts.ReadRequest{Id: id})
		if err != nil {
			return err
		}

		buffer.WriteString(fmt.Sprintf("Type: %s\n", resp.Type.String()))
		buffer.WriteString(fmt.Sprintf("Data: %s\n", string(resp.Data)))
		return nil
	})
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Read failed after retries: %v\n", err))
	}

	return &Result{
		Buffer: buffer,
	}
}

// HandleDelete обработчик удаления данных.
func (h *Handler) HandleDelete(ctx context.Context, id int64) *Result {
	buffer := NewBuffer()

	err := h.withRetryAuth(ctx, func(ctx context.Context) error {
		resp, err := h.Client.Delete(ctx, &contracts.DeleteRequest{Id: id})
		if err != nil {
			return err
		}

		if resp.Success {
			buffer.WriteString(fmt.Sprintf("Item %d deleted successfully\n", id))
		} else {
			buffer.WriteString(fmt.Sprintf("Delete item %d failed\n", id))
		}
		return nil
	})
	if err != nil {
		buffer.WriteString(fmt.Sprintf("Delete failed after retries: %v\n", err))
	}

	return &Result{Buffer: buffer}
}

// HandleList обработчик получения данных пользователя.
func (h *Handler) HandleList(ctx context.Context) *Result {
	buffer := NewBuffer()

	err := h.withRetryAuth(ctx, func(ctx context.Context) error {
		resp, err := h.Client.List(ctx, &contracts.ListRequest{})
		if err != nil {
			return err
		}

		buffer.WriteString("📋 Your items:\n")
		if len(resp.DataList) == 0 {
			buffer.WriteString("   (empty)\n")
			return nil
		}

		for _, item := range resp.DataList {
			buffer.WriteString(fmt.Sprintf(" ID: %d | %s | %s\n",
				item.Id, item.Type.String(), item.Title))
		}
		return nil
	})
	if err != nil {
		buffer.WriteString(fmt.Sprintf("List failed after retries: %v\n", err))
	}

	return &Result{
		Buffer: buffer,
	}
}

// withRetryAuth выполняет операцию с аутентификацией и ретраями.
func (h *Handler) withRetryAuth(ctx context.Context, operation func(context.Context) error) error {
	t, err := h.TokenStorage.Load()
	if err != nil {
		return fmt.Errorf("no token, login first: %w", err)
	}

	return h.withRetry(ctx, func(ctx context.Context) error {
		authCtx := metadata.AppendToOutgoingContext(ctx, "token", t)
		return operation(authCtx)
	})
}

// withRetry выполняет операцию с ретраями.
func (h *Handler) withRetry(ctx context.Context, operation func(context.Context) error) error {
	var lastErr error
	delay := h.RetryDelay

	for attempt := 0; attempt <= h.MaxRetries; attempt++ {
		if attempt > 0 {

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}

			delay = time.Duration(float64(delay) * h.BackoffFactor)
		}

		err := operation(ctx)
		if err == nil {
			return nil
		}

		lastErr = err
		if !h.shouldRetry(err) {
			return err
		}

		if attempt == h.MaxRetries {
			break
		}
	}

	return fmt.Errorf("after %d attempts: %w", h.MaxRetries, lastErr)
}

// shouldRetry проверяет, нужно ли повторять операцию при данной ошибке.
func (h *Handler) shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	for _, code := range h.RetryableCodes {
		if st.Code() == code {
			return true
		}
	}
	return false
}

// parseDataType преобразует строку в DataType.
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
