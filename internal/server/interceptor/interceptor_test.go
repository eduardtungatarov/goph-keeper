package interceptor_test

import (
	"context"
	"testing"

	"github.com/eduardtungatarov/goph-keeper/internal/server/service/auth"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/require"

	"github.com/eduardtungatarov/goph-keeper/internal/server/interceptor/mocks"

	"github.com/eduardtungatarov/goph-keeper/internal/server/interceptor"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// TestInterceptor_AuthInterceptor_PublicMethods проверяем что паблик методы выполняются без токена.
func TestInterceptor_AuthInterceptor_PublicMethods(t *testing.T) {
	authSvc := &mocks.AuthService{}
	i := interceptor.New(authSvc)

	// Public methods не требуют токена
	publicMethods := []string{
		"/server.KeeperService/Login",
		"/server.KeeperService/Register",
	}

	for _, method := range publicMethods {
		t.Run(method, func(t *testing.T) {
			info := &grpc.UnaryServerInfo{FullMethod: method}
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "success", nil
			}

			resp, err := i.AuthInterceptor(context.Background(), nil, info, grpc.UnaryHandler(handler))

			require.NoError(t, err)
			assert.Equal(t, "success", resp)
			authSvc.AssertNotCalled(t, "GetUserIDByToken")
		})
	}
}

// TestInterceptor_AuthInterceptor_NoToken проверяем что приват методы требуют токена.
func TestInterceptor_AuthInterceptor_NoToken(t *testing.T) {
	authSvc := &mocks.AuthService{}
	i := interceptor.New(authSvc)

	// private methods требуют токена
	privateMethods := []string{
		"/server.KeeperService/Create",
		"/server.KeeperService/Read",
		"/server.KeeperService/Delete",
		"/server.KeeperService/List",
	}

	for _, method := range privateMethods {
		t.Run(method, func(t *testing.T) {
			info := &grpc.UnaryServerInfo{FullMethod: method}
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "success", nil
			}

			resp, err := i.AuthInterceptor(context.Background(), nil, info, grpc.UnaryHandler(handler))

			assert.Nil(t, resp)
			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, codes.Unauthenticated, st.Code())
			assert.Equal(t, "authorization required", st.Message())
			authSvc.AssertNotCalled(t, "GetUserIDByToken")
		})
	}
}

// TestInterceptor_AuthInterceptor_ValidToken проверяем сценарий с переданным токеном.
func TestInterceptor_AuthInterceptor_ValidToken(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"token": "valid-token",
	}))
	authSvc := &mocks.AuthService{}
	i := interceptor.New(authSvc)

	authSvc.On("GetUserIDByToken", "valid-token").Return(123, nil)

	var capturedUserID int
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		// Проверяем, что userID добавлен в контекст
		userID, ok := ctx.Value(auth.UserIDKeyName).(int)
		require.True(t, ok)
		assert.Equal(t, 123, userID)
		capturedUserID = userID
		return "success", nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/server.KeeperService/Create"}

	resp, err := i.AuthInterceptor(ctx, nil, info, grpc.UnaryHandler(handler))

	assert.NoError(t, err)
	assert.Equal(t, "success", resp)
	assert.Equal(t, 123, capturedUserID)
	authSvc.AssertExpectations(t)
}
