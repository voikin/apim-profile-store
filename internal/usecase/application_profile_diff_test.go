package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/voikin/apim-profile-store/internal/entity"
	mock_usecase "github.com/voikin/apim-profile-store/internal/usecase/mocks"
)

func TestUsecase_DiffApplicationProfiles(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		setupMocks   func(mc *minimock.Controller) (*mock_usecase.PostgresRepoMock, *mock_usecase.Neo4jRepoMock)
		expectedErr  string
		expectedAdd  []*entity.Operation
		expectedRem  []*entity.Operation
	}

	appID := uuid.New()
	oldProfileID := uuid.New()
	newProfileID := uuid.New()
	graphIDOld := uuid.New()
	graphIDNew := uuid.New()

	op1 := &entity.Operation{ID: uuid.NewString(), Method: "GET", PathSegmentID: uuid.NewString()}
	op2 := &entity.Operation{ID: uuid.NewString(), Method: "POST", PathSegmentID: uuid.NewString()}
	op3 := &entity.Operation{ID: uuid.NewString(), Method: "DELETE", PathSegmentID: uuid.NewString()}

	tests := []testCase{
		{
			name: "happy path with diff",
			setupMocks: func(mc *minimock.Controller) (*mock_usecase.PostgresRepoMock, *mock_usecase.Neo4jRepoMock) {
				postgres := mock_usecase.NewPostgresRepoMock(mc)
				neo4j := mock_usecase.NewNeo4jRepoMock(mc)

				postgres.GetApplicationMock.Expect(context.Background(), appID).
					Return(&entity.Application{ID: appID}, nil)

				postgres.GetApplicationProfileByIDMock.
					Expect(context.Background(), oldProfileID).Return(&entity.ApplicationProfile{
					ID:            oldProfileID,
					ApplicationID: appID,
					GraphID:       graphIDOld,
				}, nil)

				postgres.GetApplicationProfileByIDMock.
					Expect(context.Background(), newProfileID).Return(&entity.ApplicationProfile{
					ID:            newProfileID,
					ApplicationID: appID,
					GraphID:       graphIDNew,
				}, nil)

				neo4j.GetAPIGraphMock.Expect(context.Background(), graphIDOld).
					Return(&entity.APIGraph{Operations: []*entity.Operation{op1, op2}}, nil)

				neo4j.GetAPIGraphMock.Expect(context.Background(), graphIDNew).
					Return(&entity.APIGraph{Operations: []*entity.Operation{op2, op3}}, nil)

				return postgres, neo4j
			},
			expectedAdd: []*entity.Operation{op3},
			expectedRem: []*entity.Operation{op1},
		},
		{
			name: "error fetching application",
			setupMocks: func(mc *minimock.Controller) (*mock_usecase.PostgresRepoMock, *mock_usecase.Neo4jRepoMock) {
				postgres := mock_usecase.NewPostgresRepoMock(mc)
				neo4j := mock_usecase.NewNeo4jRepoMock(mc)

				postgres.GetApplicationMock.Expect(context.Background(), appID).
					Return(nil, errors.New("db error"))

				return postgres, neo4j
			},
			expectedErr: "u.postgresRepo.GetApplication: db error",
		},
		{
			name: "error fetching old profile",
			setupMocks: func(mc *minimock.Controller) (*mock_usecase.PostgresRepoMock, *mock_usecase.Neo4jRepoMock) {
				postgres := mock_usecase.NewPostgresRepoMock(mc)
				neo4j := mock_usecase.NewNeo4jRepoMock(mc)

				postgres.GetApplicationMock.Expect(context.Background(), appID).
					Return(&entity.Application{ID: appID}, nil)

				postgres.GetApplicationProfileByIDMock.
					Expect(context.Background(), oldProfileID).Return(nil, errors.New("not found"))

				return postgres, neo4j
			},
			expectedErr: "u.postgresRepo.GetApplicationProfileByID: old: not found",
		},
		{
			name: "old profile not owned by application",
			setupMocks: func(mc *minimock.Controller) (*mock_usecase.PostgresRepoMock, *mock_usecase.Neo4jRepoMock) {
				postgres := mock_usecase.NewPostgresRepoMock(mc)
				neo4j := mock_usecase.NewNeo4jRepoMock(mc)

				postgres.GetApplicationMock.Expect(context.Background(), appID).
					Return(&entity.Application{ID: appID}, nil)

				postgres.GetApplicationProfileByIDMock.
					Expect(context.Background(), oldProfileID).Return(&entity.ApplicationProfile{
					ID:            oldProfileID,
					ApplicationID: uuid.New(), // другое приложение
					GraphID:       graphIDOld,
				}, nil)

				return postgres, neo4j
			},
			expectedErr: "old profile is related to other application",
		},
		{
			name: "error fetching new profile",
			setupMocks: func(mc *minimock.Controller) (*mock_usecase.PostgresRepoMock, *mock_usecase.Neo4jRepoMock) {
				postgres := mock_usecase.NewPostgresRepoMock(mc)
				neo4j := mock_usecase.NewNeo4jRepoMock(mc)

				postgres.GetApplicationMock.Expect(context.Background(), appID).
					Return(&entity.Application{ID: appID}, nil)

				postgres.GetApplicationProfileByIDMock.
					Expect(context.Background(), oldProfileID).Return(&entity.ApplicationProfile{
					ID:            oldProfileID,
					ApplicationID: appID,
					GraphID:       graphIDOld,
				}, nil)

				postgres.GetApplicationProfileByIDMock.
					Expect(context.Background(), newProfileID).Return(nil, errors.New("not found"))

				return postgres, neo4j
			},
			expectedErr: "u.postgresRepo.GetApplicationProfileByID: new: not found",
		},
		{
			name: "new profile not owned by application",
			setupMocks: func(mc *minimock.Controller) (*mock_usecase.PostgresRepoMock, *mock_usecase.Neo4jRepoMock) {
				postgres := mock_usecase.NewPostgresRepoMock(mc)
				neo4j := mock_usecase.NewNeo4jRepoMock(mc)

				postgres.GetApplicationMock.Expect(context.Background(), appID).
					Return(&entity.Application{ID: appID}, nil)

				postgres.GetApplicationProfileByIDMock.
					Expect(context.Background(), oldProfileID).Return(&entity.ApplicationProfile{
					ID:            oldProfileID,
					ApplicationID: appID,
					GraphID:       graphIDOld,
				}, nil)

				postgres.GetApplicationProfileByIDMock.
					Expect(context.Background(), newProfileID).Return(&entity.ApplicationProfile{
					ID:            newProfileID,
					ApplicationID: uuid.New(), // другое приложение
					GraphID:       graphIDNew,
				}, nil)

				return postgres, neo4j
			},
			expectedErr: "new profile is related to other application",
		},
		{
			name: "error getting old graph",
			setupMocks: func(mc *minimock.Controller) (*mock_usecase.PostgresRepoMock, *mock_usecase.Neo4jRepoMock) {
				postgres := mock_usecase.NewPostgresRepoMock(mc)
				neo4j := mock_usecase.NewNeo4jRepoMock(mc)

				postgres.GetApplicationMock.Expect(context.Background(), appID).
					Return(&entity.Application{ID: appID}, nil)

				postgres.GetApplicationProfileByIDMock.
					Expect(context.Background(), oldProfileID).Return(&entity.ApplicationProfile{
					ID:            oldProfileID,
					ApplicationID: appID,
					GraphID:       graphIDOld,
				}, nil)

				postgres.GetApplicationProfileByIDMock.
					Expect(context.Background(), newProfileID).Return(&entity.ApplicationProfile{
					ID:            newProfileID,
					ApplicationID: appID,
					GraphID:       graphIDNew,
				}, nil)

				neo4j.GetAPIGraphMock.Expect(context.Background(), graphIDOld).
					Return(nil, errors.New("neo4j error"))

				return postgres, neo4j
			},
			expectedErr: "u.neo4jRepo.GetAPIGraph: old: neo4j error",
		},
		{
			name: "error getting new graph",
			setupMocks: func(mc *minimock.Controller) (*mock_usecase.PostgresRepoMock, *mock_usecase.Neo4jRepoMock) {
				postgres := mock_usecase.NewPostgresRepoMock(mc)
				neo4j := mock_usecase.NewNeo4jRepoMock(mc)

				postgres.GetApplicationMock.Expect(context.Background(), appID).
					Return(&entity.Application{ID: appID}, nil)

				postgres.GetApplicationProfileByIDMock.
					Expect(context.Background(), oldProfileID).Return(&entity.ApplicationProfile{
					ID:            oldProfileID,
					ApplicationID: appID,
					GraphID:       graphIDOld,
				}, nil)

				postgres.GetApplicationProfileByIDMock.
					Expect(context.Background(), newProfileID).Return(&entity.ApplicationProfile{
					ID:            newProfileID,
					ApplicationID: appID,
					GraphID:       graphIDNew,
				}, nil)

				neo4j.GetAPIGraphMock.Expect(context.Background(), graphIDOld).
					Return(&entity.APIGraph{Operations: []*entity.Operation{}}, nil)

				neo4j.GetAPIGraphMock.Expect(context.Background(), graphIDNew).
					Return(nil, errors.New("neo4j failed"))

				return postgres, neo4j
			},
			expectedErr: "u.neo4jRepo.GetAPIGraph: new: neo4j failed",
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			pg, neo := tc.setupMocks(mc)
			u := New(pg, neo, nil)

			added, removed, err := u.DiffApplicationProfiles(context.Background(), appID, oldProfileID, newProfileID)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
				require.ElementsMatch(t, tc.expectedAdd, added)
				require.ElementsMatch(t, tc.expectedRem, removed)
			}
		})
	}
}
