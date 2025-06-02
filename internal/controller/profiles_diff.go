package controller

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
	profilestorepb "github.com/voikin/apim-proto/gen/go/apim_profile_store/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) DiffProfiles(
	ctx context.Context,
	req *profilestorepb.DiffProfilesRequest,
) (*profilestorepb.DiffProfilesResponse, error) {
	appID, err := uuid.Parse(req.ApplicationId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	oldID, err := uuid.Parse(req.OldProfileId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	newID, err := uuid.Parse(req.NewProfileId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	added, removed, err := c.usecase.DiffApplicationProfiles(ctx, appID, oldID, newID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
	}

	return &profilestorepb.DiffProfilesResponse{
		Added:   toProtoOperations(added),
		Removed: toProtoOperations(removed),
	}, nil
}
