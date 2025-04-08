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

func (r *Repo) CreateApplicationProfile(ctx context.Context, profile *entity.ApplicationProfile) (uuid.UUID, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	id := uuid.New()

	sql, args, err := sq.Insert(ApplicationProfileTableName).
		Columns(
			ApplicationProfileColumnID,
			ApplicationProfileColumnApplicationID,
			ApplicationProfileColumnVersion,
			ApplicationProfileColumnGraphID,
		).
		Values(
			id,
			profile.ApplicationID,
			profile.Version,
			profile.GraphID,
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
