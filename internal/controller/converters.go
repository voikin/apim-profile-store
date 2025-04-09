package controller

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
	shared "github.com/voikin/apim-proto/gen/go/shared/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func applicationToAPI(application *entity.Application) *shared.Application {
	return &shared.Application{
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

func applicationProfileToAPI(applicationProfile *entity.ApplicationProfile, _ string) *shared.ApplicationProfile {
	apiApplicationProfile := &shared.ApplicationProfile{
		Id:            applicationProfile.ID.String(),
		ApplicationId: applicationProfile.ApplicationID.String(),
		Version:       applicationProfile.Version,
		CreatedAt:     timestamppb.New(applicationProfile.CreatedAt),
	}

	return apiApplicationProfile
}

func applicationProfileFromAPI(
	applicationID string,
	apiGraph *shared.APIGraph,
) (*entity.ApplicationProfile, string, error) {
	applicationUUID, err := uuid.Parse(applicationID)
	if err != nil {
		return nil, "", fmt.Errorf("uuid.Parse: application_id: %w", err)
	}

	raw, err := json.Marshal(apiGraph)
	if err != nil {
		return nil, "", fmt.Errorf("json.Marshal: %w", err)
	}

	return &entity.ApplicationProfile{
		ApplicationID: applicationUUID,
	}, string(raw), nil
}
