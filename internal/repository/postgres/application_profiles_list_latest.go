package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (r *Repo) ListLatestApplicationProfiles(ctx context.Context) ([]*entity.ApplicationProfile, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	alias := "ap"

	subquery := sq.Select("MAX(" + ApplicationProfileColumnVersion + ")").
		From(ApplicationProfileTableName).
		Where(alias + "." + ApplicationProfileColumnApplicationID + " = " + ApplicationProfileTableName + "." + ApplicationProfileColumnApplicationID)

	query := sq.Select(ApplicationProfileAllColumns...).
		From(ApplicationProfileTableName + " " + alias).
		Where(sq.Expr(alias+"."+ApplicationProfileColumnVersion+" = (?)", subquery)).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
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
