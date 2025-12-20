package handler

import (
	"context"
	"errors"

	"github.com/eduardtungatarov/goph-keeper/internal/server/repository/data/queries"

	"github.com/eduardtungatarov/goph-keeper/internal/server/repository"

	"github.com/eduardtungatarov/goph-keeper/internal/server/service/data/dto"

	"github.com/eduardtungatarov/goph-keeper/internal/server/service/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"
	userRepository "github.com/eduardtungatarov/goph-keeper/internal/server/repository/user"
)

// AuthService сервис аутентификации.
type AuthService interface {
	Register(ctx context.Context, login, pwd string) (string, error)
	Login(ctx context.Context, login, pwd string) (string, error)
}

// DataService сервис по работе с данными пользователей.
type DataService interface {
	Create(ctx context.Context, create dto.Create) error
	Read(ctx context.Context, create dto.Read) (dto.ReadResult, error)
	Delete(ctx context.Context, read dto.Delete) error
	List(ctx context.Context) ([]queries.Datum, error)
}

// Handler обработчик методов контракта.
type Handler struct {
	contracts.UnimplementedKeeperServiceServer
	authService AuthService
	dataService DataService
}

// New конструктор обработчиков.
func New(
	authService AuthService,
	dataService DataService,
) *Handler {
	return &Handler{
		authService: authService,
		dataService: dataService,
	}
}

// Register регистрация пользователя.
func (h *Handler) Register(ctx context.Context, req *contracts.LoginRequest) (*contracts.LoginResponse, error) {
	if req.Login == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password are required")
	}

	token, err := h.authService.Register(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, userRepository.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "registration failed")
	}

	return &contracts.LoginResponse{
		Token: token,
	}, nil
}

// Login аутентификация пользователя.
func (h *Handler) Login(ctx context.Context, req *contracts.LoginRequest) (*contracts.LoginResponse, error) {
	if req.Login == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password are required")
	}

	token, err := h.authService.Login(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrLoginPwd) {
			return nil, status.Error(codes.Unauthenticated, "")
		}

		return nil, status.Error(codes.Internal, "login failed")
	}

	return &contracts.LoginResponse{
		Token: token,
	}, nil
}

// Create сохранение данных пользователя.
func (h *Handler) Create(ctx context.Context, req *contracts.CreateRequest) (*contracts.CreateResponse, error) {
	err := h.dataService.Create(ctx, dto.Create{
		Type:  req.GetType().String(),
		Title: req.GetTitle(),
		Data:  req.GetData(),
	})
	if err != nil {
		if errors.Is(err, repository.ErrDataTooBig) {
			return nil, status.Error(codes.InvalidArgument, "data too large: maximum size is 1MB")
		}
		return nil, err
	}

	return &contracts.CreateResponse{
		Success: true,
	}, nil
}

// Read чтение данных пользователя.
func (h *Handler) Read(ctx context.Context, req *contracts.ReadRequest) (*contracts.ReadResponse, error) {
	data, err := h.dataService.Read(ctx, dto.Read{
		ID: int(req.GetId()),
	})
	if err != nil {
		if errors.Is(err, repository.ErrNoModel) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		return nil, err
	}
	return &contracts.ReadResponse{
		Type: contracts.DataType(contracts.DataType_value[data.Type]),
		Data: data.Data,
	}, nil
}

// Delete удаление данных пользователя.
func (h *Handler) Delete(ctx context.Context, req *contracts.DeleteRequest) (*contracts.DeleteResponse, error) {
	err := h.dataService.Delete(ctx, dto.Delete{
		ID: int(req.GetId()),
	})
	if err != nil {
		if errors.Is(err, repository.ErrNoModel) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		return nil, err
	}
	return &contracts.DeleteResponse{
		Success: true,
	}, nil
}

// List список данных пользователя.
func (h *Handler) List(ctx context.Context, _ *contracts.ListRequest) (*contracts.ListResponse, error) {
	datas, err := h.dataService.List(ctx)
	if err != nil {
		return nil, err
	}

	var dataList []*contracts.DataItem

	for _, v := range datas {
		dataList = append(dataList, &contracts.DataItem{
			Id:    v.ID,
			Type:  contracts.DataType(contracts.DataType_value[v.Type]),
			Data:  v.Data,
			Title: v.Title,
		})
	}

	return &contracts.ListResponse{
		DataList: dataList,
	}, nil
}
