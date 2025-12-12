package data

import (
	"context"
	"fmt"

	"github.com/eduardtungatarov/goph-keeper/internal/server"

	"github.com/eduardtungatarov/goph-keeper/internal/server/repository/data/queries"

	"github.com/eduardtungatarov/goph-keeper/internal/server/service/data/dto"
)

type Repository interface {
	Save(ctx context.Context, data queries.Datum) (queries.Datum, error)
	GetByUserIDAndID(ctx context.Context, userID, id int) (queries.Datum, error)
	DeleteByUserIDAndID(ctx context.Context, userID, id int) error
	ListByUserID(ctx context.Context, userID int) ([]queries.Datum, error)
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, read dto.Create) error {
	const op = "data.Service.Create"

	userID, err := server.GetUserID(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	model := queries.Datum{
		UserID: int64(userID),
		Type:   read.Type,
		Title:  read.Title,
		Data:   read.Data,
	}
	_, err = s.repository.Save(ctx, model)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Read(ctx context.Context, read dto.Read) (dto.ReadResult, error) {
	const op = "data.Service.Read"

	userID, err := server.GetUserID(ctx)
	if err != nil {
		return dto.ReadResult{}, fmt.Errorf("%s: %w", op, err)
	}

	data, err := s.repository.GetByUserIDAndID(ctx, userID, read.ID)
	if err != nil {
		return dto.ReadResult{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.ReadResult{
		Type: data.Type,
		Data: data.Data,
	}, nil
}

func (s *Service) Delete(ctx context.Context, read dto.Delete) error {
	const op = "data.Service.Delete"

	userID, err := server.GetUserID(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.repository.DeleteByUserIDAndID(ctx, userID, read.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) List(ctx context.Context) ([]queries.Datum, error) {
	const op = "data.Service.List"

	userID, err := server.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.repository.ListByUserID(ctx, userID)
}
