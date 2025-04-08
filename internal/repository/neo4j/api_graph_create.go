package neo4j

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repo) CreateAPIGraph(ctx context.Context, data string) (uuid.UUID, error) {
	id := uuid.New()

	tx, ok := TxFromContext(ctx)
	if !ok {
		return uuid.Nil, errors.New("no neo4j transaction in context")
	}

	result, err := tx.Run(ctx,
		"CREATE (g:Graph {id: $id, data: $data})",
		map[string]any{
			"id":   id.String(),
			"data": data,
		})
	if err != nil {
		return uuid.Nil, fmt.Errorf("tx.Run: %w", err)
	}

	if _, err := result.Consume(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("result.Consume: %w", err)
	}

	return id, nil
}
