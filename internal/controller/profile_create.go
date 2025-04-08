package controller

import (
	"context"
	"fmt"

	profilestorepb "github.com/voikin/apim-profile-store/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) AddProfile(
	ctx context.Context,
	req *profilestorepb.AddProfileRequest,
) (*profilestorepb.AddProfileResponse, error) {
	profile, graphData, err := applicationProfileFromAPI(req.GetApplicationId(), req.GetGraph())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	createdProfile, graph, err := c.usecase.CreateApplicationProfile(ctx, profile, graphData)
	if err != nil {
		return nil, fmt.Errorf("c.usecase.CreateApplicationProfile: %w", err)
	}

	return &profilestorepb.AddProfileResponse{
		Profile: applicationProfileToAPI(createdProfile, graph),
	}, nil
}
