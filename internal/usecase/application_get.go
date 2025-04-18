package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (u *Usecase) GetApplication(ctx context.Context, id uuid.UUID) (*entity.Application, error) {
	app, err := u.postgresRepo.GetApplication(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("postgresRepo.GetApplication: %w", err)
	}

	return app, nil
}
