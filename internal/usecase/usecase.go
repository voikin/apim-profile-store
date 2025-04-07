package usecase

import (
	"context"

	"github.com/voikin/apim-profile-store/internal/entity"
)

type (
	TrManager interface {
		Do(ctx context.Context, fn func(ctx context.Context) error) error
	}

	PostgresRepo interface {
		ListApplications(ctx context.Context) ([]*entity.Application, error)
	}
)

type Usecase struct {
	postgresRepo PostgresRepo

	trManager TrManager
}

func New(postgresRepo PostgresRepo, trManager TrManager) *Usecase {
	return &Usecase{
		postgresRepo: postgresRepo,
		trManager:    trManager,
	}
}
