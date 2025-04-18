package controller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	profilestorepb "github.com/voikin/apim-proto/gen/go/apim_profile_store/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) GetProfileByVersion(
	ctx context.Context,
	req *profilestorepb.GetProfileByVersionRequest,
) (*profilestorepb.GetProfileByVersionResponse, error) {
	applicationID, err := uuid.Parse(req.GetApplicationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	profile, graph, err := c.usecase.GetApplicationProfileByVersion(ctx, applicationID, req.GetVersion())
	if err != nil {
		return nil, fmt.Errorf("c.usecase.GetApplicationProfileByVersion: %w", err)
	}

	return &profilestorepb.GetProfileByVersionResponse{
		Profile: applicationProfileToAPI(profile, graph),
	}, nil
}
