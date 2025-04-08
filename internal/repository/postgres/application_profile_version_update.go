package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *Repo) UpdateLatestProfileVersion(ctx context.Context, applicationID uuid.UUID, version uint32) error {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	sql, args, err := sq.Update(ApplicationProfileVersionsTableName).
		Set(ApplicationProfileVersionsColumnLastVersion, version).
		Where(sq.Eq{ApplicationProfileVersionsColumnApplicationID: applicationID}).
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
		sql, args, err = sq.Insert(ApplicationProfileVersionsTableName).
			Columns(
				ApplicationProfileVersionsColumnApplicationID,
				ApplicationProfileVersionsColumnLastVersion,
			).
			Values(applicationID, version).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("sq.ToSql: %w", err)
		}

		_, err = conn.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("conn.Exec: %w", err)
		}
	}

	return nil
}
