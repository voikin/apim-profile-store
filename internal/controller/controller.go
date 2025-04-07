package controller

import (
	"context"

	"github.com/voikin/apim-profile-store/internal/entity"
	v1 "github.com/voikin/apim-profile-store/pkg/api/v1"
)

type Usecase interface {
	ListApplications(ctx context.Context) ([]*entity.Application, error)
}

type Controller struct {
	usecase Usecase

	v1.UnimplementedProfileStoreServiceServer
}

func New(usecase Usecase) *Controller {
	return &Controller{
		usecase: usecase,
	}
}
