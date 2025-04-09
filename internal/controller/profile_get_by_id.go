package controller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	profilestorepb "github.com/voikin/apim-proto/gen/go/apim_profile_store/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) GetProfileByID(
	ctx context.Context,
	req *profilestorepb.GetProfileByIDRequest,
) (*profilestorepb.GetProfileByIDResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	profile, graph, err := c.usecase.GetApplicationProfileByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("c.usecase.GetApplicationProfileByID: %w", err)
	}

	return &profilestorepb.GetProfileByIDResponse{
		Profile: applicationProfileToAPI(profile, graph),
	}, nil
}
