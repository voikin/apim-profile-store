package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
	profilestorepb "github.com/voikin/apim-profile-store/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) DeleteApplication(
	ctx context.Context,
	req *profilestorepb.DeleteApplicationRequest,
) (*profilestorepb.DeleteApplicationResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	err = c.usecase.DeleteApplication(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "application not found")
		}

		return nil, fmt.Errorf("c.usecase.DeleteApplication: %w", err)
	}

	return &profilestorepb.DeleteApplicationResponse{}, nil
}
