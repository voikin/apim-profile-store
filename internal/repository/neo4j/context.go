package neo4j

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ctxKey string

const txKey ctxKey = "neo4jTx"

func withTx(ctx context.Context, tx neo4j.ManagedTransaction) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func TxFromContext(ctx context.Context) (neo4j.ManagedTransaction, bool) {
	tx, ok := ctx.Value(txKey).(neo4j.ManagedTransaction)
	return tx, ok
}

func (r *Repo) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.trManager.Do(ctx, func(ctx context.Context) error {
		session := r.driver.NewSession(ctx, neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
		})
		defer session.Close(ctx)

		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			ctxWithTx := withTx(ctx, tx)
			return nil, fn(ctxWithTx)
		})
		if err != nil {
			return fmt.Errorf("neo4j.ExecuteWrite: %w", err)
		}

		return nil
	})
}
