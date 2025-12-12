package data

import (
	"context"
	"fmt"

	"github.com/eduardtungatarov/goph-keeper/internal/server/repository/data/queries"

	"github.com/eduardtungatarov/goph-keeper/internal/server/service/data/dto"
)

type Repository interface {
	Save(ctx context.Context, data queries.Datum) (queries.Datum, error)
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, create dto.Create) error {
	const op = "data.Service.Create"

	model := queries.Datum{
		Type:  create.Type,
		Title: create.Title,
		Data:  create.Data,
	}
	_, err := s.repository.Save(ctx, model)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
