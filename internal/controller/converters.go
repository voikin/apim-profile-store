package controller

import (
	"fmt"

	"github.com/google/uuid"
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

func applicationFromAPI(name string) *entity.Application {
	return &entity.Application{
		Name: name,
	}
}

func applicationProfileToAPI(applicationProfile *entity.ApplicationProfile, graph string) *profilestorepb.ApplicationProfile {
	apiApplicationProfile := &profilestorepb.ApplicationProfile{
		Id:            applicationProfile.ID.String(),
		ApplicationId: applicationProfile.ApplicationID.String(),
		Version:       applicationProfile.Version,
		CreatedAt:     timestamppb.New(applicationProfile.CreatedAt),
	}

	if graph != "" {
		apiApplicationProfile.Graph = &profilestorepb.ProfileGraph{
			Data: graph,
		}
	}

	return apiApplicationProfile
}

func applicationProfileFromAPI(applicationID string, graph *profilestorepb.ProfileGraph) (*entity.ApplicationProfile, string, error) {
	applicationUUID, err := uuid.Parse(applicationID)
	if err != nil {
		return nil, "", fmt.Errorf("uuid.Parse: application_id: %w", err)
	}

	return &entity.ApplicationProfile{
		ApplicationID: applicationUUID,
	}, graph.GetData(), nil
}
