package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23514" && pgErr.ConstraintName == "data_size_limit" {
				return queries.Datum{}, repository.ErrDataTooBig
			}
		}
		return queries.Datum{}, fmt.Errorf("%s: %w", op, err)
	}

	return model, err
}

func (r *Repository) GetByUserIDAndID(ctx context.Context, userID, ID int) (queries.Datum, error) {
	const op = "data.Repository.GetByUserIDAndID"

	data, err := r.querier.FindDataByUserIDAndID(ctx, r.db, queries.FindDataByUserIDAndIDParams{
		UserID: int64(userID),
		ID:     int64(ID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return queries.Datum{}, repository.ErrNoModel
		}
		return queries.Datum{}, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}

func (r *Repository) DeleteByUserIDAndID(ctx context.Context, userID, ID int) error {
	const op = "data.Repository.DeleteByUserIDAndID"

	rowsAffected, err := r.querier.DeleteDataByUserIDAndID(ctx, r.db, queries.DeleteDataByUserIDAndIDParams{
		UserID: int64(userID),
		ID:     int64(ID),
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, repository.ErrNoModel)
	}

	return nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID int) ([]queries.Datum, error) {
	const op = "data.Repository.ListByUserID"

	datas, err := r.querier.FindDataByUserID(ctx, r.db, int64(userID))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return datas, nil
}
