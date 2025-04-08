package controller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	profilestorepb "github.com/voikin/apim-profile-store/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) ListProfilesByApplication(
	ctx context.Context,
	req *profilestorepb.ListProfilesByApplicationRequest,
) (*profilestorepb.ListProfilesByApplicationResponse, error) {
	applicationID, err := uuid.Parse(req.GetApplicationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	profiles, err := c.usecase.ListApplicationProfiles(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("c.usecase.ListApplicationProfiles: %w", err)
	}

	apiProfiles := make([]*profilestorepb.ApplicationProfile, len(profiles))
	for i, profile := range profiles {
		apiProfiles[i] = applicationProfileToAPI(profile, "")
	}

	return &profilestorepb.ListProfilesByApplicationResponse{
		Profiles: apiProfiles,
	}, nil
}
