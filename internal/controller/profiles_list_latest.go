package controller

import (
	"context"
	"fmt"

	profilestorepb "github.com/voikin/apim-proto/gen/go/apim_profile_store/v1"
	shared "github.com/voikin/apim-proto/gen/go/shared/v1"
)

func (c *Controller) ListLatestProfiles(
	ctx context.Context,
	_ *profilestorepb.ListLatestProfilesRequest,
) (*profilestorepb.ListLatestProfilesResponse, error) {
	profiles, err := c.usecase.ListLatestApplicationProfiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("c.usecase.ListLatestApplicationProfiles: %w", err)
	}

	apiProfiles := make([]*shared.ApplicationProfile, len(profiles))
	for i, profile := range profiles {
		apiProfiles[i] = applicationProfileToAPI(profile, nil)
	}

	return &profilestorepb.ListLatestProfilesResponse{
		Profiles: apiProfiles,
	}, nil
}
