package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/eduardtungatarov/goph-keeper/internal/server/repository"

	"github.com/eduardtungatarov/goph-keeper/internal/server/repository/data/queries"
)

type Repository struct {
	db      queries.DBTX
	querier queries.Querier
}

func New(db queries.DBTX) *Repository {
	return &Repository{
		db:      db,
		querier: queries.New(),
	}
}

func (r *Repository) Save(ctx context.Context, data queries.Datum) (queries.Datum, error) {
	const op = "data.Repository.Save"

	model, err := r.querier.SaveData(ctx, r.db, queries.SaveDataParams{
		UserID: data.UserID,
		Title:  data.Title,
		Type:   data.Type,
		Data:   data.Data,
	})
	if err != nil {
		return queries.Datum{}, fmt.Errorf("%s: %w", op, err)
	}

	return model, err
}

func (r *Repository) GetByUserIDAndID(ctx context.Context, userID, id int) (queries.Datum, error) {
	const op = "data.Repository.GetByUserIDAndID"

	data, err := r.querier.FindDataByUserIDAndID(ctx, r.db, queries.FindDataByUserIDAndIDParams{
		UserID: int64(userID),
		ID:     int64(id),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return queries.Datum{}, repository.ErrNoModel
		}
		return queries.Datum{}, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}
