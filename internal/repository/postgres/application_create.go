package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (r *Repo) CreateApplication(ctx context.Context, app *entity.Application) (uuid.UUID, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	id := uuid.New()

	sql, args, err := sq.Insert(ApplicationTableName).
		Columns(
			ApplicationColumnID,
			ApplicationColumnName,
		).
		Values(
			id.String(),
			app.Name,
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("sq.ToSql: %w", err)
	}

	_, err = conn.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
			return uuid.Nil, entity.ErrAlreadyExists
		}

		return uuid.Nil, fmt.Errorf("conn.Exec: %w", err)
	}

	return id, nil
}
