package neo4j

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (r *Repo) GetAPIGraph(ctx context.Context, graphID uuid.UUID) (*entity.APIGraph, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	graph, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		prefix := fmt.Sprintf("g%s::", graphID.String())

		// Извлечение сегментов
		segmentsQuery := `
			MATCH (s:PathSegment)-[:BELONGS_TO]->(g:Graph {id: $graphID})
			RETURN s.id AS id, s.name AS name, s.is_param AS is_param, s.type AS type, s.example AS example
		`

		segmentResult, err := tx.Run(ctx, segmentsQuery, map[string]any{"graphID": graphID.String()})
		if err != nil {
			return nil, fmt.Errorf("get segments: %w", err)
		}

		var segments []*entity.PathSegment
		for segmentResult.Next(ctx) {
			record := segmentResult.Record()
			id, _ := record.Get("id")
			name, _ := record.Get("name")
			isParam, _ := record.Get("is_param")
			typ, _ := record.Get("type")
			example, _ := record.Get("example")

			originalID := strings.TrimPrefix(id.(string), prefix)

			if isParam.(bool) {
				segments = append(segments, &entity.PathSegment{
					Param: &entity.Parameter{
						ID:      originalID,
						Name:    name.(string),
						Type:    entity.ParameterType(int32(typ.(int64))),
						Example: stringOrEmpty(example),
					},
				})
			} else {
				segments = append(segments, &entity.PathSegment{
					Static: &entity.StaticSegment{
						ID:   originalID,
						Name: name.(string),
					},
				})
			}
		}
		if err := segmentResult.Err(); err != nil {
			return nil, fmt.Errorf("segment iteration: %w", err)
		}

		// Извлечение рёбер между сегментами
		edgesQuery := `
			MATCH (from:PathSegment)-[:FOLLOWS]->(to:PathSegment)
			WHERE from.id STARTS WITH $prefix AND to.id STARTS WITH $prefix
			RETURN from.id AS from, to.id AS to
		`
		edgeResult, err := tx.Run(ctx, edgesQuery, map[string]any{"prefix": prefix})
		if err != nil {
			return nil, fmt.Errorf("get edges: %w", err)
		}

		var edges []*entity.Edge
		for edgeResult.Next(ctx) {
			record := edgeResult.Record()
			from, _ := record.Get("from")
			to, _ := record.Get("to")

			edges = append(edges, &entity.Edge{
				From: strings.TrimPrefix(from.(string), prefix),
				To:   strings.TrimPrefix(to.(string), prefix),
			})
		}
		if err := edgeResult.Err(); err != nil {
			return nil, fmt.Errorf("edge iteration: %w", err)
		}

		// Извлечение операций
		opQuery := `
			MATCH (op:Operation)-[:ON]->(seg:PathSegment), (op)-[:BELONGS_TO]->(g:Graph {id: $graphID})
			RETURN op.id AS id, op.method AS method, seg.id AS segment_id, op.status_codes AS status_codes
		`
		opResult, err := tx.Run(ctx, opQuery, map[string]any{"graphID": graphID.String()})
		if err != nil {
			return nil, fmt.Errorf("get operations: %w", err)
		}

		var operations []*entity.Operation
		for opResult.Next(ctx) {
			record := opResult.Record()
			id, _ := record.Get("id")
			method, _ := record.Get("method")
			segmentID, _ := record.Get("segment_id")
			statusCodesRaw, _ := record.Get("status_codes")

			var statusCodes []int32
			if statusCodesRaw != nil {
				for _, val := range statusCodesRaw.([]any) {
					statusCodes = append(statusCodes, int32(val.(int64)))
				}
			}

			operations = append(operations, &entity.Operation{
				ID:            strings.TrimPrefix(id.(string), prefix),
				Method:        method.(string),
				PathSegmentID: strings.TrimPrefix(segmentID.(string), prefix),
				StatusCodes:   statusCodes,
				// Если будут добавлены дополнительные параметры запроса, их можно здесь восстановить
				QueryParameters: nil,
			})
		}
		if err := opResult.Err(); err != nil {
			return nil, fmt.Errorf("operation iteration: %w", err)
		}

		// Извлечение переходов между операциями
		transitionQuery := `
			MATCH (from:Operation)-[:NEXT]->(to:Operation)
			WHERE from.id STARTS WITH $prefix AND to.id STARTS WITH $prefix
			RETURN from.id AS from, to.id AS to
		`
		transitionResult, err := tx.Run(ctx, transitionQuery, map[string]any{"prefix": prefix})
		if err != nil {
			return nil, fmt.Errorf("get transitions: %w", err)
		}

		var transitions []*entity.Transition
		for transitionResult.Next(ctx) {
			record := transitionResult.Record()
			from, _ := record.Get("from")
			to, _ := record.Get("to")

			transitions = append(transitions, &entity.Transition{
				From: strings.TrimPrefix(from.(string), prefix),
				To:   strings.TrimPrefix(to.(string), prefix),
			})
		}
		if err := transitionResult.Err(); err != nil {
			return nil, fmt.Errorf("transition iteration: %w", err)
		}

		return &entity.APIGraph{
			Segments:    segments,
			Edges:       edges,
			Operations:  operations,
			Transitions: transitions,
		}, nil
	})

	if err != nil {
		return nil, fmt.Errorf("session.ExecuteRead: %w", err)
	}

	return graph.(*entity.APIGraph), nil
}

func stringOrEmpty(v any) string {
	if v == nil {
		return ""
	}
	return v.(string)
}
