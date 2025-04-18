package neo4j

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (r *Repo) CreateAPIGraph(ctx context.Context, apiGraph *entity.APIGraph) (uuid.UUID, error) {
	id := uuid.New()
	prefix := fmt.Sprintf("g%s::", id.String())

	tx, ok := TxFromContext(ctx)
	if !ok {
		return uuid.Nil, errors.New("no neo4j transaction in context")
	}

	_, err := tx.Run(ctx,
		`CREATE (g:Graph {id: $graphID})`,
		map[string]any{"graphID": id.String()},
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create Graph: %w", err)
	}

	segmentIDMap := make(map[string]string)

	for _, segment := range apiGraph.Segments {
		var originalID, name string
		params := map[string]any{
			"graphID":  id.String(),
			"is_param": segment.Param != nil,
			"type":     nil,
			"example":  nil,
		}

		if segment.Static != nil {
			originalID = segment.Static.ID
			name = segment.Static.Name
			params["name"] = name
		} else if segment.Param != nil {
			originalID = segment.Param.ID
			name = segment.Param.Name
			params["name"] = name
			params["type"] = int(segment.Param.Type)
			params["example"] = segment.Param.Example
		}

		uniqueID := prefix + originalID
		params["id"] = uniqueID
		segmentIDMap[originalID] = uniqueID

		_, err := tx.Run(ctx,
			`MATCH (g:Graph {id: $graphID})
			 CREATE (s:PathSegment {
				id: $id,
				name: $name,
				is_param: $is_param,
				type: $type,
				example: $example
			 })-[:BELONGS_TO]->(g)`,
			params,
		)
		if err != nil {
			return uuid.Nil, fmt.Errorf("create PathSegment: %w", err)
		}
	}

	for _, edge := range apiGraph.Edges {
		fromID := prefix + edge.From
		toID := prefix + edge.To

		_, err := tx.Run(ctx,
			`MATCH (from:PathSegment {id: $from}), (to:PathSegment {id: $to})
			 CREATE (from)-[:FOLLOWS]->(to)`,
			map[string]any{
				"from": fromID,
				"to":   toID,
			},
		)
		if err != nil {
			return uuid.Nil, fmt.Errorf("create Edge: %w", err)
		}
	}

	for _, op := range apiGraph.Operations {
		originalSegID := op.PathSegmentID
		uniqueSegID := segmentIDMap[originalSegID]
		uniqueOpID := prefix + op.ID

		_, err := tx.Run(ctx,
			`MATCH (g:Graph {id: $graphID}), (seg:PathSegment {id: $segmentID})
			 CREATE (op:Operation {
				id: $id,
				method: $method,
				status_codes: $status_codes
			 })-[:ON]->(seg)
			 CREATE (op)-[:BELONGS_TO]->(g)`,
			map[string]any{
				"graphID":      id.String(),
				"id":           uniqueOpID,
				"method":       op.Method,
				"segmentID":    uniqueSegID,
				"status_codes": op.StatusCodes,
			},
		)
		if err != nil {
			return uuid.Nil, fmt.Errorf("create Operation: %w", err)
		}
	}

	return id, nil
}
