package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/eduardtungatarov/goph-keeper/internal/server/repository"
	dataQueries "github.com/eduardtungatarov/goph-keeper/internal/server/repository/data/queries"

	"github.com/stretchr/testify/require"

	userRepository "github.com/eduardtungatarov/goph-keeper/internal/server/repository/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/mock"

	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"
	"github.com/eduardtungatarov/goph-keeper/internal/server/handler"
	"github.com/eduardtungatarov/goph-keeper/internal/server/repository/user/queries"
	"github.com/eduardtungatarov/goph-keeper/internal/server/service/auth"
	"github.com/eduardtungatarov/goph-keeper/internal/server/service/auth/mocks"
	"github.com/eduardtungatarov/goph-keeper/internal/server/service/data"
	authDataMocks "github.com/eduardtungatarov/goph-keeper/internal/server/service/data/mocks"
	"github.com/eduardtungatarov/goph-keeper/internal/server/service/security"
	"github.com/stretchr/testify/assert"
)

func makeHandler(t *testing.T) (*handler.Handler, *mocks.UserRepository, *authDataMocks.Repository, *security.Service) {
	userRepo := mocks.NewUserRepository(t)
	dataRepo := authDataMocks.NewRepository(t)

	authSvc := auth.New("jwt-secret-key", userRepo)
	secSvc := security.New([]byte("cc3a7e2c1b8d4a6f"))
	dataSvc := data.New(dataRepo, secSvc)

	h := handler.New(authSvc, dataSvc)

	return h, userRepo, dataRepo, secSvc
}

// TestHandler_Register_Success регистрация пользователя - успешный сценарий.
func TestHandler_Register_Success(t *testing.T) {
	h, userRepo, _, _ := makeHandler(t)

	userRepo.On("SaveUser", mock.Anything, mock.MatchedBy(func(user queries.User) bool {
		return user.Login == "user"
	})).Return(queries.User{ID: 1, Login: "user"}, nil)

	req := &contracts.LoginRequest{Login: "user", Password: "pwd"}
	resp, err := h.Register(context.Background(), req)

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
}

// TestHandler_Register_InvalidInput регистрация пользователя - неверный ввод.
func TestHandler_Register_InvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		password string
	}{
		{"empty login", "", "pwd"},
		{"empty password", "user", ""},
		{"both empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, _, _, _ := makeHandler(t)

			req := &contracts.LoginRequest{Login: tt.login, Password: tt.password}
			_, err := h.Register(context.Background(), req)

			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, codes.InvalidArgument, st.Code())
		})
	}
}

// TestHandler_Register_UserAlreadyExists регистрация пользователя - пользователь уже был зарегистрирован.
func TestHandler_Register_UserAlreadyExists(t *testing.T) {
	h, userRepo, _, _ := makeHandler(t)

	userRepo.On("SaveUser", mock.Anything, mock.Anything).Return(queries.User{}, userRepository.ErrUserAlreadyExists)

	req := &contracts.LoginRequest{Login: "user", Password: "pwd"}
	_, err := h.Register(context.Background(), req)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())
}

// TestHandler_NoUserIDInContext_Register вызов метода без user ID в контексте.
func TestHandler_NoUserIDInContext_Register(t *testing.T) {
	ctx := context.Background() // Без userID
	h, _, _, _ := makeHandler(t)

	req := &contracts.LoginRequest{}
	_, err := h.Register(ctx, req)

	assert.Error(t, err)
}

// TestHandler_Login_Success вход пользователя - успешный сценарий.
func TestHandler_Login_Success(t *testing.T) {
	h, userRepo, _, _ := makeHandler(t)

	userRepo.On("FindUserByLogin", mock.Anything, "user").Return(queries.User{
		ID:       1,
		Login:    "user",
		Password: "$2a$14$pavsQ63/1bTfrCO8NcGGNeWDG9a20pa6uBzeNVH3aw/0QbsXruOT.", // хеш от пароля pwd
	}, nil)

	req := &contracts.LoginRequest{Login: "user", Password: "pwd"}
	resp, err := h.Login(context.Background(), req)

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
}

// TestHandler_Login_InvalidInput регистрация пользователя - неверный ввод.
func TestHandler_Login_InvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		password string
	}{
		{"empty login", "", "pwd"},
		{"empty password", "user", ""},
		{"both empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, _, _, _ := makeHandler(t)

			req := &contracts.LoginRequest{Login: tt.login, Password: tt.password}
			_, err := h.Login(context.Background(), req)

			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, codes.InvalidArgument, st.Code())
		})
	}
}

// TestHandler_Login_InvalidCredentials_NotFound вход - пользователь не найден.
func TestHandler_Login_InvalidCredentials_NotFound(t *testing.T) {
	h, userRepo, _, _ := makeHandler(t)

	userRepo.On("FindUserByLogin", mock.Anything, "user").Return(queries.User{}, repository.ErrNoModel)

	req := &contracts.LoginRequest{Login: "user", Password: "wrong_pwd"}
	_, err := h.Login(context.Background(), req)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

// TestHandler_NoUserIDInContext_Login вызов метода без user ID в контексте.
func TestHandler_NoUserIDInContext_Login(t *testing.T) {
	ctx := context.Background() // Без userID
	h, _, _, _ := makeHandler(t)

	req := &contracts.LoginRequest{}
	_, err := h.Login(ctx, req)

	assert.Error(t, err)
}

// TestHandler_Create_Success сохраняем данные пользователя - успешный сценарий.
func TestHandler_Create_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.UserIDKeyName, 1)
	h, _, dataRepo, _ := makeHandler(t)

	dataRepo.On("Save", mock.Anything, mock.MatchedBy(func(datum dataQueries.Datum) bool {
		return datum.UserID == 1 && datum.Title == "github" && datum.Type == "PWD"
	})).Return(dataQueries.Datum{ID: 1}, nil)

	req := &contracts.CreateRequest{
		Type:  contracts.DataType_PWD,
		Title: "github",
		Data:  []byte("github_login|github_pwd"),
	}
	resp, err := h.Create(ctx, req)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

// TestHandler_Create_DataTooBig сохраняем данные пользователя - данные слишком большие.
func TestHandler_Create_DataTooBig(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.UserIDKeyName, 1)
	h, _, dataRepo, _ := makeHandler(t)

	dataRepo.On("Save", mock.Anything, mock.Anything).Return(dataQueries.Datum{}, repository.ErrDataTooBig)

	req := &contracts.CreateRequest{
		Type:  contracts.DataType_PWD,
		Title: "github",
		Data:  []byte("github_login|github_pwd"),
	}
	_, err := h.Create(ctx, req)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// TestHandler_Read_Success чтение данных пользователя - успешный сценарий.
func TestHandler_Read_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.UserIDKeyName, 1)
	h, _, dataRepo, secSrv := makeHandler(t)
	encrypted, err := secSrv.GetEncrypted(context.Background(), []byte("github_login|github_pwd"))
	require.NoError(t, err)

	dataRepo.On("GetByUserIDAndID", mock.Anything, 1, 123).Return(dataQueries.Datum{
		ID:   123,
		Type: "PWD",
		Data: encrypted,
	}, nil)

	req := &contracts.ReadRequest{Id: 123}
	resp, err := h.Read(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, contracts.DataType_PWD, resp.Type)
	assert.Equal(t, "github_login|github_pwd", string(resp.Data))
}

// TestHandler_Read_NotFound чтение данных пользователя - не найдено данных.
func TestHandler_Read_NotFound(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.UserIDKeyName, 1)
	h, _, dataRepo, _ := makeHandler(t)

	dataRepo.On("GetByUserIDAndID", mock.Anything, 1, 123).Return(dataQueries.Datum{}, repository.ErrNoModel)

	req := &contracts.ReadRequest{Id: 123}
	_, err := h.Read(ctx, req)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

// TestHandler_NoUserIDInContext_Read вызов метода без user ID в контексте.
func TestHandler_NoUserIDInContext_Read(t *testing.T) {
	ctx := context.Background() // Без userID
	h, _, _, _ := makeHandler(t)

	req := &contracts.ReadRequest{}
	_, err := h.Read(ctx, req)

	assert.Error(t, err)
}

// TestHandler_Delete_Success удаление данных - успешный сценарий.
func TestHandler_Delete_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.UserIDKeyName, 1)
	h, _, dataRepo, _ := makeHandler(t)

	dataRepo.On("DeleteByUserIDAndID", mock.Anything, 1, 123).Return(nil)

	req := &contracts.DeleteRequest{Id: 123}
	resp, err := h.Delete(ctx, req)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

// TestHandler_Delete_NotFound удаление данных - не найдено данных.
func TestHandler_Delete_NotFound(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.UserIDKeyName, 1)
	h, _, dataRepo, _ := makeHandler(t)

	dataRepo.On("DeleteByUserIDAndID", mock.Anything, 1, 123).Return(repository.ErrNoModel)

	req := &contracts.DeleteRequest{Id: 123}
	_, err := h.Delete(ctx, req)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

// TestHandler_NoUserIDInContext_Delete вызов метода без user ID в контексте.
func TestHandler_NoUserIDInContext_Delete(t *testing.T) {
	ctx := context.Background() // Без userID
	h, _, _, _ := makeHandler(t)

	req := &contracts.DeleteRequest{}
	_, err := h.Delete(ctx, req)

	assert.Error(t, err)
}

// TestHandler_List_Success список данных - успешный сценарий.
func TestHandler_List_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.UserIDKeyName, 1)
	h, _, dataRepo, _ := makeHandler(t)

	dataRepo.On("ListByUserID", mock.Anything, 1).Return([]dataQueries.Datum{
		{ID: 1, Type: "PWD", Title: "github"},
		{ID: 2, Type: "CARD", Title: "SBER"},
		{ID: 3, Type: "BINARY", Title: "test.txt"},
	}, nil)

	req := &contracts.ListRequest{}
	resp, err := h.List(ctx, req)

	require.NoError(t, err)
	assert.Len(t, resp.DataList, 3)

	assert.Equal(t, int64(1), resp.DataList[0].Id)
	assert.Equal(t, contracts.DataType_PWD, resp.DataList[0].Type)
	assert.Equal(t, "github", resp.DataList[0].Title)

	assert.Equal(t, int64(2), resp.DataList[1].Id)
	assert.Equal(t, contracts.DataType_CARD, resp.DataList[1].Type)
	assert.Equal(t, "SBER", resp.DataList[1].Title)

	assert.Equal(t, int64(3), resp.DataList[2].Id)
	assert.Equal(t, contracts.DataType_BINARY, resp.DataList[2].Type)
	assert.Equal(t, "test.txt", resp.DataList[2].Title)
}

// TestHandler_List_Error список данных - ошибка бд.
func TestHandler_List_Error(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.UserIDKeyName, 1)
	h, _, dataRepo, _ := makeHandler(t)

	dataRepo.On("ListByUserID", mock.Anything, 1).Return([]dataQueries.Datum{}, errors.New("db error"))

	req := &contracts.ListRequest{}
	_, err := h.List(ctx, req)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}
