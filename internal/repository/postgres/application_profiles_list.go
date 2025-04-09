package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (r *Repo) ListApplicationProfiles(
	ctx context.Context,
	applicationID uuid.UUID,
) ([]*entity.ApplicationProfile, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	sql, args, err := sq.Select(ApplicationProfileAllColumns...).
		From(ApplicationProfileTableName).
		Where(sq.Eq{ApplicationProfileColumnApplicationID: applicationID}).
		OrderBy("version DESC").
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

	profiles := make([]*entity.ApplicationProfile, 0)
	for rows.Next() {
		profile := &entity.ApplicationProfile{}
		err = rows.Scan(
			&profile.ID,
			&profile.ApplicationID,
			&profile.Version,
			&profile.GraphID,
			&profile.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}
