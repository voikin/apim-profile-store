package usecase

import (
	"context"
	"fmt"

	"github.com/voikin/apim-profile-store/internal/entity"
)

func (u *Usecase) ListApplications(ctx context.Context) ([]*entity.Application, error) {
	applications, err := u.postgresRepo.ListApplications(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgresRepo.ListApplications: %w", err)
	}

	return applications, nil
}
