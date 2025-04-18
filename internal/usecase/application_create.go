package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (u *Usecase) CreateApplication(ctx context.Context, app *entity.Application) (*entity.Application, error) {
	var appID uuid.UUID

	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		id, err := u.postgresRepo.CreateApplication(ctx, app)
		if err != nil {
			return fmt.Errorf("postgresRepo.CreateApplication: %w", err)
		}

		appID = id

		err = u.postgresRepo.UpdateLatestProfileVersion(ctx, appID, 0)
		if err != nil {
			return fmt.Errorf("u.postgresRepo.UpdateLatestProfileVersion: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("trManager.Do: %w", err)
	}

	application, err := u.postgresRepo.GetApplication(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("u.postgresRepo.GetApplication: %w", err)
	}

	return application, nil
}
