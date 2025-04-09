package neo4j

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r *Repo) GetAPIGraph(ctx context.Context, id uuid.UUID) (string, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	data, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx,
			"MATCH (g:Graph {id: $id}) RETURN g.data",
			map[string]any{
				"id": id.String(),
			})
		if err != nil {
			return nil, fmt.Errorf("tx.Run: %w", err)
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, fmt.Errorf("result.Single: %w", err)
		}

		return record.Values[0], nil
	})
	if err != nil {
		return "", fmt.Errorf("session.ExecuteRead: %w", err)
	}

	res, _ := data.(string)

	return res, nil
}
