package controller

import (
	"context"
	"fmt"

	profilestorepb "github.com/voikin/apim-profile-store/pkg/api/v1"
)

func (c *Controller) ListLatestProfiles(
	ctx context.Context,
	req *profilestorepb.ListLatestProfilesRequest,
) (*profilestorepb.ListLatestProfilesResponse, error) {
	profiles, err := c.usecase.ListLatestApplicationProfiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("c.usecase.ListLatestApplicationProfiles: %w", err)
	}

	apiProfiles := make([]*profilestorepb.ApplicationProfile, len(profiles))
	for i, profile := range profiles {
		apiProfiles[i] = applicationProfileToAPI(profile, "")
	}

	return &profilestorepb.ListLatestProfilesResponse{
		Profiles: apiProfiles,
	}, nil
}
