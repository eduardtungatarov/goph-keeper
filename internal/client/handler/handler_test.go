package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/eduardtungatarov/goph-keeper/internal/client/handler"
	"github.com/eduardtungatarov/goph-keeper/internal/client/handler/mocks"
	"github.com/eduardtungatarov/goph-keeper/internal/server/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_HandleLogin_Success_TokenSaveOK(t *testing.T) {
	client := mocks.NewClient(t)
	tokenStorage := mocks.NewTokenStorage(t)

	client.On("Login", mock.Anything, mock.MatchedBy(func(req *contracts.LoginRequest) bool {
		return req.Login == "testuser" && req.Password == "testpass"
	}), mock.Anything).Return(&contracts.LoginResponse{Token: "token123"}, nil)
	tokenStorage.On("Save", "token123").Return(nil)

	h := handler.New(client, tokenStorage)
	result := h.HandleLogin(context.Background(), "testuser", "testpass")

	output := result.Buffer.String()
	assert.Contains(t, output, "Login successful!")
	assert.Contains(t, output, "Token saved to ~/.keeper/token.json")
	assert.NotContains(t, output, "Login failed")

	client.AssertExpectations(t)
	tokenStorage.AssertExpectations(t)
}

func TestHandler_HandleLogin_Success_TokenSaveError(t *testing.T) {
	client := mocks.NewClient(t)
	tokenStorage := mocks.NewTokenStorage(t)

	client.On("Login", mock.Anything, mock.Anything, mock.Anything).
		Return(&contracts.LoginResponse{Token: "token123"}, nil)
	tokenStorage.On("Save", "token123").Return(errors.New("disk full"))

	h := handler.New(client, tokenStorage)
	result := h.HandleLogin(context.Background(), "testuser", "testpass")

	output := result.Buffer.String()
	assert.Contains(t, output, "Login successful!")
	assert.Contains(t, output, "Failed to save token: disk full")
	assert.NotContains(t, output, "Login failed")

	client.AssertExpectations(t)
	tokenStorage.AssertExpectations(t)
}

func TestHandler_HandleLogin_Failure(t *testing.T) {
	client := mocks.NewClient(t)
	tokenStorage := mocks.NewTokenStorage(t)

	client.On("Login", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, status.Error(codes.InvalidArgument, "invalid credentials"))

	h := handler.New(client, tokenStorage)
	result := h.HandleLogin(context.Background(), "testuser", "testpass")

	output := result.Buffer.String()
	assert.Contains(t, output, "Login failed after retries: rpc error: code = InvalidArgument")

	client.AssertExpectations(t)
	tokenStorage.AssertNotCalled(t, "Save", mock.Anything)
}

func TestHandler_HandleRegister_Success(t *testing.T) {
	client := mocks.NewClient(t)
	tokenStorage := mocks.NewTokenStorage(t)

	client.On("Register", mock.Anything, mock.MatchedBy(func(req *contracts.LoginRequest) bool {
		return req.Login == "newuser" && req.Password == "secret"
	}), mock.Anything).Return(&contracts.LoginResponse{Token: "token456"}, nil)
	tokenStorage.On("Save", "token456").Return(nil)

	h := handler.New(client, tokenStorage)
	result := h.HandleRegister(context.Background(), "newuser", "secret")

	output := result.Buffer.String()
	assert.Contains(t, output, "Register successful! Token: token456")
	assert.Contains(t, output, "Token saved to ~/.keeper/token.json")

	client.AssertExpectations(t)
	tokenStorage.AssertExpectations(t)
}

func TestHandler_HandleCreate_InvalidType(t *testing.T) {
	client := mocks.NewClient(t)
	tokenStorage := mocks.NewTokenStorage(t)
	h := handler.New(client, tokenStorage)

	result := h.HandleCreate(context.Background(), "unknown", "data", "title")

	output := result.Buffer.String()
	assert.Contains(t, output, "Invalid type: unknown. Use: pwd, card, binary")

	// TokenStorage не вызывается
	tokenStorage.AssertNotCalled(t, "Load", mock.Anything)
	client.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

func TestHandler_HandleCreate_Success(t *testing.T) {
	client := mocks.NewClient(t)
	tokenStorage := mocks.NewTokenStorage(t)

	tokenStorage.On("Load").Return("valid-token", nil)
	client.On("Create",
		mock.Anything,
		mock.MatchedBy(func(req *contracts.CreateRequest) bool {
			return req.Type == contracts.DataType_PWD &&
				string(req.Data) == "user|pass" &&
				req.Title == "github"
		}),
	).Return(&contracts.CreateResponse{Success: true}, nil)

	h := handler.New(client, tokenStorage)
	result := h.HandleCreate(context.Background(), "pwd", "user|pass", "github")

	output := result.Buffer.String()
	assert.Contains(t, output, "Item 'github' created successfully")

	client.AssertExpectations(t)
	tokenStorage.AssertExpectations(t)
}

func TestHandler_HandleCreate_NoToken(t *testing.T) {
	client := mocks.NewClient(t)
	tokenStorage := mocks.NewTokenStorage(t)

	tokenStorage.On("Load").Return("", errors.New("no token file"))

	h := handler.New(client, tokenStorage)
	result := h.HandleCreate(context.Background(), "pwd", "data", "title")

	output := result.Buffer.String()
	assert.Contains(t, output, "Create failed after retries: no token, login first: no token file")

	client.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
	tokenStorage.AssertExpectations(t)
}

func TestHandler_HandleRead_Success(t *testing.T) {
	client := mocks.NewClient(t)
	tokenStorage := mocks.NewTokenStorage(t)

	tokenStorage.On("Load").Return("token123", nil)
	client.On("Read",
		mock.Anything,
		mock.MatchedBy(func(req *contracts.ReadRequest) bool {
			return req.Id == 123
		}),
		mock.Anything,
	).Return(&contracts.ReadResponse{
		Type: contracts.DataType_PWD,
		Data: []byte("user|pass"),
	}, nil)

	h := handler.New(client, tokenStorage)
	result := h.HandleRead(context.Background(), 123)

	output := result.Buffer.String()
	assert.Contains(t, output, "Type: PWD")
	assert.Contains(t, output, "Data: user|pass")

	client.AssertExpectations(t)
	tokenStorage.AssertExpectations(t)
}

func TestHandler_HandleDelete_Success(t *testing.T) {
	client := mocks.NewClient(t)
	tokenStorage := mocks.NewTokenStorage(t)

	tokenStorage.On("Load").Return("token", nil)
	client.On("Delete",
		mock.Anything,
		mock.MatchedBy(func(req *contracts.DeleteRequest) bool {
			return req.Id == 456
		}),
		mock.Anything,
	).Return(&contracts.DeleteResponse{Success: true}, nil)

	h := handler.New(client, tokenStorage)
	result := h.HandleDelete(context.Background(), 456)

	output := result.Buffer.String()
	assert.Contains(t, output, "Item 456 deleted successfully")

	client.AssertExpectations(t)
	tokenStorage.AssertExpectations(t)
}

func TestHandler_HandleList_Empty(t *testing.T) {
	client := mocks.NewClient(t)
	tokenStorage := mocks.NewTokenStorage(t)

	tokenStorage.On("Load").Return("token", nil)
	client.On("List", mock.Anything, mock.Anything, mock.Anything).
		Return(&contracts.ListResponse{}, nil)

	h := handler.New(client, tokenStorage)
	result := h.HandleList(context.Background())

	output := result.Buffer.String()
	assert.Contains(t, output, "📋 Your items:")
	assert.Contains(t, output, "(empty)")

	client.AssertExpectations(t)
	tokenStorage.AssertExpectations(t)
}

func TestHandler_HandleList_WithItems(t *testing.T) {
	client := mocks.NewClient(t)
	tokenStorage := mocks.NewTokenStorage(t)

	tokenStorage.On("Load").Return("token", nil)
	listResp := &contracts.ListResponse{
		DataList: []*contracts.DataItem{
			{Id: 1, Type: contracts.DataType_PWD, Title: "github"},
			{Id: 2, Type: contracts.DataType_CARD, Title: "SberCard"},
		},
	}
	client.On("List", mock.Anything, mock.Anything, mock.Anything).Return(listResp, nil)

	h := handler.New(client, tokenStorage)
	result := h.HandleList(context.Background())

	output := result.Buffer.String()
	assert.Contains(t, output, "📋 Your items:")
	assert.Contains(t, output, "ID: 1 | PWD | github")
	assert.Contains(t, output, "ID: 2 | CARD | SberCard")

	client.AssertExpectations(t)
	tokenStorage.AssertExpectations(t)
}
