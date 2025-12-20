package interceptor

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/eduardtungatarov/goph-keeper/internal/server/service/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthService сервис аутентификации.
//
//go:generate mockery --name=AuthService
type AuthService interface {
	GetUserIDByToken(tokenStr string) (int, error)
}

// Interceptor grpc сервера.
type Interceptor struct {
	authService AuthService
}

// New создание interceptor.
func New(
	authService AuthService,
) *Interceptor {
	return &Interceptor{
		authService: authService,
	}
}

// AuthInterceptor interceptor аутентификации.
func (i *Interceptor) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Endpoints без аутентификации.
	publicMethods := map[string]bool{
		"/server.KeeperService/Login":    true,
		"/server.KeeperService/Register": true,
	}
	if publicMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("token")) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization required")
	}

	authHeaders := md.Get("token")
	userID, err := i.authService.GetUserIDByToken(authHeaders[0])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "token invalid")
	}

	ctx = context.WithValue(ctx, auth.UserIDKeyName, userID)
	return handler(ctx, req)
}
