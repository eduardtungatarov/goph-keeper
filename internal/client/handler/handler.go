package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/eduardtungatarov/goph-keeper/internal/client/token"
	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"
)

// Handler обработчик команд.
type Handler struct {
	Client         contracts.KeeperServiceClient // Сгенерированный grpc клиент
	MaxRetries     int                           // Максимальное количество попыток
	RetryDelay     time.Duration                 // Задержка между попытками
	BackoffFactor  float64                       // Коэффициент экспоненциальной задержки
	RetryableCodes []codes.Code                  // Коды ошибок, при которых нужно повторять
}

// New конструктор обработчиков команд.
func New(client contracts.KeeperServiceClient) *Handler {
	return &Handler{
		Client:        client,
		MaxRetries:    3,               // По умолчанию 3 попытки
		RetryDelay:    1 * time.Second, // Начальная задержка 1 секунда
		BackoffFactor: 2.0,             // Удваиваем задержку каждый раз
		RetryableCodes: []codes.Code{
			codes.Unavailable,       // Сервер недоступен
			codes.DeadlineExceeded,  // Превышено время ожидания
			codes.ResourceExhausted, // Ресурсы исчерпаны
			codes.Internal,          // Внутренняя ошибка сервера
		},
	}
}

// HandleLogin обработчик входа.
func (h *Handler) HandleLogin(ctx context.Context, login, password string) {
	err := h.withRetry(ctx, func(ctx context.Context) error {
		resp, err := h.Client.Login(ctx, &contracts.LoginRequest{
			Login:    login,
			Password: password,
		})
		if err != nil {
			return err
		}

		fmt.Println("Login successful!")
		if err := token.Save(resp.Token); err != nil {
			log.Printf("Failed to save token: %v", err)
		} else {
			fmt.Println("Token saved to ~/.keeper/token.json")
		}
		return nil
	})
	if err != nil {
		log.Printf("Login failed after retries: %v", err)
	}
}

// HandleRegister обработчик регистрации.
func (h *Handler) HandleRegister(ctx context.Context, login, password string) {
	err := h.withRetry(ctx, func(ctx context.Context) error {
		resp, err := h.Client.Register(ctx, &contracts.LoginRequest{
			Login:    login,
			Password: password,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Register successful! Token: %s\n", resp.Token)
		if err := token.Save(resp.Token); err != nil {
			log.Printf("Failed to save token: %v", err)
		} else {
			fmt.Println("Token saved to ~/.keeper/token.json")
		}
		return nil
	})
	if err != nil {
		log.Printf("Register failed after retries: %v", err)
	}
}

// HandleCreate обработчик создания данных.
func (h *Handler) HandleCreate(ctx context.Context, dataTypeStr, data, title string) {
	dataType, ok := h.parseDataType(dataTypeStr)
	if !ok {
		log.Fatalf("Invalid type: %s. Use: pwd, card, binary", dataTypeStr)
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
			fmt.Printf("Item '%s' created successfully\n", title)
		} else {
			fmt.Println("Create failed")
		}
		return nil
	})
	if err != nil {
		log.Printf("Create failed after retries: %v", err)
	}
}

// HandleRead обработчик чтения данных.
func (h *Handler) HandleRead(ctx context.Context, id int64) {
	err := h.withRetryAuth(ctx, func(ctx context.Context) error {
		resp, err := h.Client.Read(ctx, &contracts.ReadRequest{
			Id: id,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Type: %s\n", resp.Type.String())
		fmt.Printf("Data: %s\n", string(resp.Data))
		return nil
	})
	if err != nil {
		log.Printf("Read failed after retries: %v", err)
	}
}

// HandleDelete обработчик удаления данных.
func (h *Handler) HandleDelete(ctx context.Context, id int64) {
	err := h.withRetryAuth(ctx, func(ctx context.Context) error {
		resp, err := h.Client.Delete(ctx, &contracts.DeleteRequest{
			Id: id,
		})
		if err != nil {
			return err
		}

		if resp.Success {
			fmt.Printf("Item %d deleted successfully\n", id)
		} else {
			fmt.Printf("Delete item %d failed\n", id)
		}
		return nil
	})
	if err != nil {
		log.Printf("Delete failed after retries: %v", err)
	}
}

// HandleList обработчик получения данных пользователя.
func (h *Handler) HandleList(ctx context.Context) {
	err := h.withRetryAuth(ctx, func(ctx context.Context) error {
		resp, err := h.Client.List(ctx, &contracts.ListRequest{})
		if err != nil {
			return err
		}

		fmt.Println("📋 Your items:")
		if len(resp.DataList) == 0 {
			fmt.Println("   (empty)")
			return nil
		}

		for _, item := range resp.DataList {
			fmt.Printf(" ID: %d | %s | %s\n",
				item.Id, item.Type.String(), item.Title)
		}
		return nil
	})
	if err != nil {
		log.Printf("List failed after retries: %v", err)
	}
}

// withRetryAuth выполняет операцию с аутентификацией и ретраями.
func (h *Handler) withRetryAuth(ctx context.Context, operation func(ctx context.Context) error) error {
	t, err := token.Load()
	if err != nil {
		log.Printf("No token, login first")
		return err
	}

	return h.withRetry(ctx, func(ctx context.Context) error {
		authCtx := metadata.AppendToOutgoingContext(ctx, "token", t)
		return operation(authCtx)
	})
}

// withRetry выполняет операцию с ретраями.
func (h *Handler) withRetry(ctx context.Context, operation func(ctx context.Context) error) error {
	var lastErr error
	delay := h.RetryDelay

	for attempt := 0; attempt <= h.MaxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("Retry attempt %d/%d (delay: %v)",
				attempt, h.MaxRetries, delay)

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
