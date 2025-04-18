package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (r *Repo) GetApplicationProfileByID(ctx context.Context, id uuid.UUID) (*entity.ApplicationProfile, error) {
	conn := r.getter.DefaultTrOrDB(ctx, r.db)

	sql, args, err := sq.Select(ApplicationProfileAllColumns...).
		From(ApplicationProfileTableName).
		Where(sq.Eq{ApplicationProfileColumnID: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("sq.ToSql: %w", err)
	}

	profile := &entity.ApplicationProfile{}
	err = conn.QueryRow(ctx, sql, args...).Scan(
		&profile.ID,
		&profile.ApplicationID,
		&profile.Version,
		&profile.GraphID,
		&profile.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("conn.QueryRow: %w", err)
	}

	return profile, nil
}
