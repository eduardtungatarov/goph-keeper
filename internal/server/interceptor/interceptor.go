package interceptor

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/eduardtungatarov/goph-keeper/internal/server/service/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthService interface {
	GetUserIDByToken(tokenStr string) (int, error)
}

type Interceptor struct {
	authService AuthService
}

func New(
	authService AuthService,
) *Interceptor {
	return &Interceptor{
		authService: authService,
	}
}

func (i *Interceptor) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	authHeaders := md.Get("token")
	if len(authHeaders) > 0 {
		userID, err := i.authService.GetUserIDByToken(authHeaders[0])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "token invalid")
		}
		ctx = context.WithValue(ctx, auth.UserIDKeyName, userID)
	}

	return handler(ctx, req)
}
