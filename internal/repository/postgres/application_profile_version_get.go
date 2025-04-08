package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (r *Repo) GetLatestVersionForUpdate(ctx context.Context, applicationID uuid.UUID) (uint32, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	sql, args, err := sq.Select(ApplicationProfileVersionsColumnLastVersion).
		From(ApplicationProfileVersionsTableName).
		Where(sq.Eq{ApplicationProfileVersionsColumnApplicationID: applicationID}).
		Suffix("FOR UPDATE").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("sq.ToSql: %w", err)
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf("conn.Query: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, entity.ErrNotFound
	}

	var version uint32
	err = rows.Scan(
		&version,
	)
	if err != nil {
		return 0, fmt.Errorf("rows.Scan: %w", err)
	}

	return version, nil
}
