package controller

import (
	"context"

	profilestorepb "github.com/voikin/apim-profile-store/pkg/api/v1"
)

func (c *Controller) ListApplications(
	ctx context.Context,
	_ *profilestorepb.ListApplicationsRequest,
) (*profilestorepb.ListApplicationsResponse, error) {
	applications, err := c.usecase.ListApplications(ctx)
	if err != nil {
		return nil, err
	}

	apiApplications := make([]*profilestorepb.Application, len(applications))
	for i, application := range applications {
		apiApplications[i] = applicationToAPI(application)
	}

	return &profilestorepb.ListApplicationsResponse{
		Applications: apiApplications,
	}, nil
}
