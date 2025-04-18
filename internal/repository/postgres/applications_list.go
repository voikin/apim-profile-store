package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (r *Repo) ListApplications(ctx context.Context) ([]*entity.Application, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	sql, args, err := sq.Select(ApplicationAllColumns...).
		From(ApplicationTableName).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("sq.ToSql: %w", err)
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("conn.Query: %w", err)
	}

	applications := make([]*entity.Application, 0)
	for rows.Next() {
		application := &entity.Application{}
		err = rows.Scan(
			&application.ID,
			&application.Name,
			&application.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		applications = append(applications, application)
	}

	return applications, nil
}
