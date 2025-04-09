package controller

import (
	"context"

	profilestorepb "github.com/voikin/apim-proto/gen/go/apim_profile_store/v1"
	shared "github.com/voikin/apim-proto/gen/go/shared/v1"
)

func (c *Controller) ListApplications(
	ctx context.Context,
	_ *profilestorepb.ListApplicationsRequest,
) (*profilestorepb.ListApplicationsResponse, error) {
	applications, err := c.usecase.ListApplications(ctx)
	if err != nil {
		return nil, err
	}

	apiApplications := make([]*shared.Application, len(applications))
	for i, application := range applications {
		apiApplications[i] = applicationToAPI(application)
	}

	return &profilestorepb.ListApplicationsResponse{
		Applications: apiApplications,
	}, nil
}
