package controller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	profilestorepb "github.com/voikin/apim-profile-store/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) DeleteProfile(
	ctx context.Context,
	req *profilestorepb.DeleteProfileRequest,
) (*profilestorepb.DeleteProfileResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := c.usecase.DeleteApplicationProfile(ctx, id); err != nil {
		return nil, fmt.Errorf("c.usecase.DeleteApplicationProfile: %w", err)
	}

	return &profilestorepb.DeleteProfileResponse{}, nil
}
