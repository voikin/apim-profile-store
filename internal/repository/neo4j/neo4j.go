package neo4j

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type TrManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Repo struct {
	driver    neo4j.DriverWithContext
	trManager TrManager
}

func New(driver neo4j.DriverWithContext, trManager TrManager) *Repo {
	return &Repo{
		driver:    driver,
		trManager: trManager,
	}
}
