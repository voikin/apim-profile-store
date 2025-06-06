package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
	profilestorepb "github.com/voikin/apim-proto/gen/go/apim_profile_store/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) GetApplication(
	ctx context.Context,
	req *profilestorepb.GetApplicationRequest,
) (*profilestorepb.GetApplicationResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	application, err := c.usecase.GetApplication(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "application not found")
		}

		return nil, fmt.Errorf("c.usecase.GetApplication: %w", err)
	}

	return &profilestorepb.GetApplicationResponse{
		Application: applicationToAPI(application),
	}, nil
}
