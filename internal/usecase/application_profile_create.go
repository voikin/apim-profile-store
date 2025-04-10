package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (u *Usecase) CreateApplicationProfile(
	ctx context.Context,
	profile *entity.ApplicationProfile,
	apiGraph *entity.APIGraph,
) (*entity.ApplicationProfile, *entity.APIGraph, error) {
	var id uuid.UUID
	err := u.trManager.Do(ctx, func(ctx context.Context) error {
		graphID, err := u.neo4jRepo.CreateAPIGraph(ctx, apiGraph)
		if err != nil {
			return fmt.Errorf("u.neo4JRepo.CreateAPIGraph: %w", err)
		}

		profile.GraphID = graphID

		latestVersion, err := u.postgresRepo.GetLatestVersionForUpdate(ctx, profile.ApplicationID)
		if err != nil {
			return fmt.Errorf("u.postgresRepo.GetLatestVersionForUpdate: %w", err)
		}

		latestVersion++
		profile.Version = latestVersion

		id, err = u.postgresRepo.CreateApplicationProfile(ctx, profile)
		if err != nil {
			return fmt.Errorf("postgresRepo.CreateApplicationProfile: %w", err)
		}

		err = u.postgresRepo.UpdateLatestProfileVersion(ctx, profile.ApplicationID, latestVersion)
		if err != nil {
			return fmt.Errorf("u.postgresRepo.UpdateLatestProfileVersion: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("trManager.Do: %w", err)
	}

	applicationProfile, apiGraph, err := u.GetApplicationProfileByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("u.GetApplicationProfileByID: %w", err)
	}

	return applicationProfile, apiGraph, nil
}
