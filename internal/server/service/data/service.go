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

type SecurityService interface {
	GetEncrypted(ctx context.Context, data []byte) ([]byte, error)
	GetDecrypted(ctx context.Context, data []byte) ([]byte, error)
}

type Service struct {
	repository      Repository
	securityService SecurityService
}

func New(repository Repository, securityService SecurityService) *Service {
	return &Service{
		repository:      repository,
		securityService: securityService,
	}
}

func (s *Service) Create(ctx context.Context, read dto.Create) error {
	const op = "data.Service.Create"

	userID, err := server.GetUserID(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	encryptedData, err := s.securityService.GetEncrypted(ctx, read.Data)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	model := queries.Datum{
		UserID: int64(userID),
		Type:   read.Type,
		Title:  read.Title,
		Data:   encryptedData,
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

	decryptedData, err := s.securityService.GetDecrypted(ctx, data.Data)
	if err != nil {
		return dto.ReadResult{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.ReadResult{
		Type: data.Type,
		Data: decryptedData,
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

	datas, err := s.repository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for i := range datas {
		decryptedData, err := s.securityService.GetDecrypted(ctx, datas[i].Data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		datas[i].Data = decryptedData
	}

	return datas, nil
}
