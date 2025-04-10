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

	query := `
		MATCH (g:Graph {id: $id})
		OPTIONAL MATCH (g)<-[:BELONGS_TO]-(s:PathSegment)
		OPTIONAL MATCH (g)<-[:BELONGS_TO]-(o:Operation)
		DETACH DELETE g, s, o
	`

	result, err := tx.Run(ctx, query, map[string]any{
		"id": id.String(),
	})
	if err != nil {
		return fmt.Errorf("tx.Run: %w", err)
	}

	_, err = result.Consume(ctx)
	if err != nil {
		return fmt.Errorf("result.Consume: %w", err)
	}

	return nil
}
