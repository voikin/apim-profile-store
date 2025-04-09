package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

type (
	TrManager interface {
		Do(ctx context.Context, fn func(ctx context.Context) error) error
	}

	PostgresRepo interface {
		CreateApplication(ctx context.Context, app *entity.Application) (uuid.UUID, error)
		GetApplication(ctx context.Context, id uuid.UUID) (*entity.Application, error)
		UpdateApplication(ctx context.Context, app *entity.Application, id uuid.UUID) error
		DeleteApplication(ctx context.Context, id uuid.UUID) error
		ListApplications(ctx context.Context) ([]*entity.Application, error)

		CreateApplicationProfile(ctx context.Context, profile *entity.ApplicationProfile) (uuid.UUID, error)
		DeleteApplicationProfile(ctx context.Context, id uuid.UUID) error
		GetApplicationProfileByID(ctx context.Context, id uuid.UUID) (*entity.ApplicationProfile, error)
		GetApplicationProfileByVersion(
			ctx context.Context,
			applicationID uuid.UUID,
			version uint32,
		) (*entity.ApplicationProfile, error)
		GetLatestApplicationProfile(
			ctx context.Context,
			applicationID uuid.UUID,
		) (*entity.ApplicationProfile, error)
		ListLatestApplicationProfiles(ctx context.Context) ([]*entity.ApplicationProfile, error)
		ListApplicationProfiles(ctx context.Context, applicationID uuid.UUID) ([]*entity.ApplicationProfile, error)

		GetLatestVersionForUpdate(ctx context.Context, applicationID uuid.UUID) (uint32, error)
		UpdateLatestProfileVersion(ctx context.Context, applicationID uuid.UUID, version uint32) error
	}

	Neo4jRepo interface {
		CreateAPIGraph(ctx context.Context, data string) (uuid.UUID, error)
		DeleteAPIGraph(ctx context.Context, id uuid.UUID) error
		GetAPIGraph(ctx context.Context, id uuid.UUID) (string, error)
	}
)

type Usecase struct {
	postgresRepo PostgresRepo
	neo4jRepo    Neo4jRepo

	trManager TrManager
}

func New(postgresRepo PostgresRepo, neo4jRepo Neo4jRepo, trManager TrManager) *Usecase {
	return &Usecase{
		postgresRepo: postgresRepo,
		neo4jRepo:    neo4jRepo,
		trManager:    trManager,
	}
}
