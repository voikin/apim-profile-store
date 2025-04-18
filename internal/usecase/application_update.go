package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (u *Usecase) UpdateApplication(ctx context.Context, app *entity.Application, id uuid.UUID) error {
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		if err := u.postgresRepo.UpdateApplication(ctx, app, id); err != nil {
			return fmt.Errorf("postgresRepo.UpdateApplication: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("trManager.Do: %w", err)
	}

	return nil
}
