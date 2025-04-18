package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (u *Usecase) DeleteApplicationProfile(ctx context.Context, id uuid.UUID) error {
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		applicationProfile, err := u.postgresRepo.GetApplicationProfileByID(ctx, id)
		if err != nil {
			return fmt.Errorf("u.postgresRepo.GetApplicationProfileByID: %w", err)
		}

		err = u.neo4jRepo.DeleteAPIGraph(ctx, applicationProfile.GraphID)
		if err != nil {
			return fmt.Errorf("u.neo4jRepo.DeleteAPIGraph: %w", err)
		}

		err = u.postgresRepo.DeleteApplicationProfile(ctx, id)
		if err != nil {
			return fmt.Errorf("postgresRepo.DeleteApplicationProfile: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("trManager.Do: %w", err)
	}

	return nil
}
