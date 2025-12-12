package data

import (
	"context"
	"fmt"

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
