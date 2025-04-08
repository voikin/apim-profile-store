package usecase

import (
	"context"
	"fmt"

	"github.com/voikin/apim-profile-store/internal/entity"
)

func (u *Usecase) ListLatestApplicationProfiles(ctx context.Context) ([]*entity.ApplicationProfile, error) {
	profiles, err := u.postgresRepo.ListLatestApplicationProfiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgresRepo.ListLatestApplicationProfiles: %w", err)
	}

	return profiles, nil
}
