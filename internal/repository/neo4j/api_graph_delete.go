package neo4j

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repo) DeleteAPIGraph(ctx context.Context, id uuid.UUID) error {
	tx, ok := TxFromContext(ctx)
	if !ok {
		return errors.New("no neo4j transaction in context")
	}

	result, err := tx.Run(ctx,
		"MATCH (g:Graph {id: $id}) DELETE g",
		map[string]any{
			"id": id.String(),
		})
	if err != nil {
		return fmt.Errorf("tx.Run: %w", err)
	}

	if _, err := result.Consume(ctx); err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	return nil
}
