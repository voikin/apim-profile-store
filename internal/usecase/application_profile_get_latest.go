package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (u *Usecase) GetLatestApplicationProfile(
	ctx context.Context,
	applicationID uuid.UUID,
) (*entity.ApplicationProfile, *entity.APIGraph, error) {
	applicationProfile, err := u.postgresRepo.GetLatestApplicationProfile(ctx, applicationID)
	if err != nil {
		return nil, nil, fmt.Errorf("u.postgresRepo.GetLatestApplicationProfile: %w", err)
	}

	graph, err := u.neo4jRepo.GetAPIGraph(ctx, applicationProfile.GraphID)
	if err != nil {
		return nil, nil, fmt.Errorf("u.neo4JRepo.GetAPIGraph: %w", err)
	}

	return applicationProfile, graph, nil
}
