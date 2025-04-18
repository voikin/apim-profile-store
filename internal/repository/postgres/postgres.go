package postgres

import (
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	TrManager interface {
		Do(ctx context.Context, fn func(ctx context.Context) error) error
	}
)

type Repo struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter

	trManager TrManager
}

func New(db *pgxpool.Pool, trManager TrManager, c *trmpgx.CtxGetter) *Repo {
	repo := &Repo{
		db:        db,
		getter:    c,
		trManager: trManager,
	}
	return repo
}
