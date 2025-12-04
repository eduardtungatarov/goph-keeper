package handler

import (
	"context"
	"errors"

	"github.com/eduardtungatarov/goph-keeper/internal/server/service/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"
	userRepository "github.com/eduardtungatarov/goph-keeper/internal/server/repository/user"
)

type AuthService interface {
	Register(ctx context.Context, login, pwd string) (string, error)
	Login(ctx context.Context, login, pwd string) (string, error)
}

type Handler struct {
	contracts.UnimplementedKeeperServiceServer
	authService AuthService
}

func New(
	authService AuthService,
) *Handler {
	return &Handler{
		authService: authService,
	}
}

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
