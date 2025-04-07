package controller

import (
	"github.com/voikin/apim-profile-store/internal/entity"
	profilestorepb "github.com/voikin/apim-profile-store/pkg/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func applicationToAPI(application *entity.Application) *profilestorepb.Application {
	return &profilestorepb.Application{
		Id:        application.ID.String(),
		Name:      application.Name,
		CreatedAt: timestamppb.New(application.CreatedAt),
	}
}
