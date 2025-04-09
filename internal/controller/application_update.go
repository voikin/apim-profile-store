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

func (c *Controller) UpdateApplication(
	ctx context.Context,
	req *profilestorepb.UpdateApplicationRequest,
) (*profilestorepb.UpdateApplicationResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	application := applicationFromAPI(req.GetName())

	err = c.usecase.UpdateApplication(ctx, application, id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			return nil, status.Error(codes.NotFound, "application not found")
		case errors.Is(err, entity.ErrAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, "application already exists")
		default:
		}

		return nil, fmt.Errorf(" c.usecase.UpdateApplication: %w", err)
	}

	return &profilestorepb.UpdateApplicationResponse{}, nil
}
