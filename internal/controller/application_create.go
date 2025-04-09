package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/voikin/apim-profile-store/internal/entity"
	profilestorepb "github.com/voikin/apim-proto/gen/go/apim_profile_store/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) CreateApplication(
	ctx context.Context,
	req *profilestorepb.CreateApplicationRequest,
) (*profilestorepb.CreateApplicationResponse, error) {
	application := applicationFromAPI(req.GetName())

	createdApp, err := c.usecase.CreateApplication(ctx, application)
	if err != nil {
		if errors.Is(err, entity.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "application already exists")
		}

		return nil, fmt.Errorf(" c.usecase.CreateApplication: %w", err)
	}

	return &profilestorepb.CreateApplicationResponse{
		Application: applicationToAPI(createdApp),
	}, nil
}
