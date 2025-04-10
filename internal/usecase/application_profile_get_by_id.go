package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (u *Usecase) GetApplicationProfileByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.ApplicationProfile, *entity.APIGraph, error) {
	applicationProfile, err := u.postgresRepo.GetApplicationProfileByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("u.postgresRepo.GetApplicationProfileByID: %w", err)
	}

	apiGraph, err := u.neo4jRepo.GetAPIGraph(ctx, applicationProfile.GraphID)
	if err != nil {
		return nil, nil, fmt.Errorf("u.neo4JRepo.GetAPIGraph: %w", err)
	}

	return applicationProfile, apiGraph, nil
}
