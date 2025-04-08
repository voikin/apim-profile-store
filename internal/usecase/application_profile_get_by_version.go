package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (u *Usecase) GetApplicationProfileByVersion(ctx context.Context, applicationID uuid.UUID, version uint32) (*entity.ApplicationProfile, string, error) {
	applicationProfile, err := u.postgresRepo.GetApplicationProfileByVersion(ctx, applicationID, version)
	if err != nil {
		return nil, "", fmt.Errorf("u.postgresRepo.GetApplicationProfileByVersion: %w", err)
	}

	graph, err := u.neo4jRepo.GetAPIGraph(ctx, applicationProfile.GraphID)
	if err != nil {
		return nil, "", fmt.Errorf("u.neo4JRepo.GetAPIGraph: %w", err)
	}

	return applicationProfile, graph, nil
}
