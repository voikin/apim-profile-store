package controller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	profilestorepb "github.com/voikin/apim-proto/gen/go/apim_profile_store/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) GetLatestProfile(
	ctx context.Context,
	req *profilestorepb.GetLatestProfileRequest,
) (*profilestorepb.GetLatestProfileResponse, error) {
	applicationID, err := uuid.Parse(req.GetApplicationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	profile, graph, err := c.usecase.GetLatestApplicationProfile(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("c.usecase.GetLatestApplicationProfile: %w", err)
	}

	return &profilestorepb.GetLatestProfileResponse{
		Profile: applicationProfileToAPI(profile, graph),
	}, nil
}
