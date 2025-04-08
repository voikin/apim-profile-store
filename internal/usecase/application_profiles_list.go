package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (u *Usecase) ListApplicationProfiles(ctx context.Context, applicationID uuid.UUID) ([]*entity.ApplicationProfile, error) {
	profiles, err := u.postgresRepo.ListApplicationProfiles(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("postgresRepo.ListApplicationProfiles: %w", err)
	}

	return profiles, nil
}
