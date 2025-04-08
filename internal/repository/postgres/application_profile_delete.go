package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (r *Repo) DeleteApplicationProfile(ctx context.Context, id uuid.UUID) error {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	sql, args, err := sq.Delete(ApplicationProfileTableName).
		Where(sq.Eq{ApplicationProfileColumnID: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("sq.ToSql: %w", err)
	}

	result, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return entity.ErrNotFound
	}

	return nil
}
