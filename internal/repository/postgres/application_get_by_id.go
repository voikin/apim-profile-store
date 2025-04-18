package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (r *Repo) GetApplication(ctx context.Context, id uuid.UUID) (*entity.Application, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	sql, args, err := sq.Select(ApplicationAllColumns...).
		From(ApplicationTableName).
		Where(sq.Eq{ApplicationColumnID: id.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("sq.ToSql: %w", err)
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("conn.Query: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, entity.ErrNotFound
	}

	app := &entity.Application{}
	err = rows.Scan(
		&app.ID,
		&app.Name,
		&app.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("rows.Scan: %w", err)
	}

	return app, nil
}
