package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
	profilestorepb "github.com/voikin/apim-proto/gen/go/apim_profile_store/v1"
)

type Usecase interface {
	CreateApplication(ctx context.Context, app *entity.Application) (*entity.Application, error)
	GetApplication(ctx context.Context, id uuid.UUID) (*entity.Application, error)
	UpdateApplication(ctx context.Context, app *entity.Application, id uuid.UUID) error
	DeleteApplication(ctx context.Context, id uuid.UUID) error
	ListApplications(ctx context.Context) ([]*entity.Application, error)

	CreateApplicationProfile(
		ctx context.Context,
		profile *entity.ApplicationProfile,
		apiGraph *entity.APIGraph,
	) (*entity.ApplicationProfile, *entity.APIGraph, error)
	DeleteApplicationProfile(ctx context.Context, id uuid.UUID) error
	GetApplicationProfileByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.ApplicationProfile, *entity.APIGraph, error)
	GetApplicationProfileByVersion(
		ctx context.Context,
		applicationID uuid.UUID,
		version uint32,
	) (*entity.ApplicationProfile, *entity.APIGraph, error)
	GetLatestApplicationProfile(
		ctx context.Context,
		applicationID uuid.UUID,
	) (*entity.ApplicationProfile, *entity.APIGraph, error)
	ListLatestApplicationProfiles(ctx context.Context) ([]*entity.ApplicationProfile, error)
	ListApplicationProfiles(
		ctx context.Context,
		applicationID uuid.UUID,
	) ([]*entity.ApplicationProfile, error)
}

type Controller struct {
	usecase Usecase

	profilestorepb.UnimplementedProfileStoreServiceServer
}

func New(usecase Usecase) *Controller {
	return &Controller{
		usecase: usecase,
	}
}
