package controller

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
	sharedpb "github.com/voikin/apim-proto/gen/go/shared/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func applicationToAPI(application *entity.Application) *sharedpb.Application {
	return &sharedpb.Application{
		Id:        application.ID.String(),
		Name:      application.Name,
		CreatedAt: timestamppb.New(application.CreatedAt),
	}
}

func applicationFromAPI(name string) *entity.Application {
	return &entity.Application{
		Name: name,
	}
}

func applicationProfileToAPI(applicationProfile *entity.ApplicationProfile, apiGraph *entity.APIGraph) *sharedpb.ApplicationProfile {
	apiApplicationProfile := &sharedpb.ApplicationProfile{
		Id:            applicationProfile.ID.String(),
		ApplicationId: applicationProfile.ApplicationID.String(),
		Version:       applicationProfile.Version,
		CreatedAt:     timestamppb.New(applicationProfile.CreatedAt),
	}

	if apiGraph != nil {
		apiApplicationProfile.ApiGraph = ToProtoGraph(apiGraph)
	}

	return apiApplicationProfile
}

func applicationProfileFromAPI(
	applicationID string,
	apiGraph *sharedpb.APIGraph,
) (*entity.ApplicationProfile, *entity.APIGraph, error) {
	applicationUUID, err := uuid.Parse(applicationID)
	if err != nil {
		return nil, nil, fmt.Errorf("uuid.Parse: application_id: %w", err)
	}

	return &entity.ApplicationProfile{
		ApplicationID: applicationUUID,
	}, ToEntityGraph(apiGraph), nil
}

func ToEntityGraph(pb *sharedpb.APIGraph) *entity.APIGraph {
	if pb == nil {
		return nil
	}

	segments := make([]entity.PathSegment, len(pb.Segments))
	for i, s := range pb.Segments {
		segments[i] = toEntitySegment(s)
	}

	edges := make([]entity.Edge, len(pb.Edges))
	for i, e := range pb.Edges {
		edges[i] = entity.Edge{
			From: e.From,
			To:   e.To,
		}
	}

	operations := make([]entity.Operation, len(pb.Operations))
	for i, op := range pb.Operations {
		operations[i] = entity.Operation{
			ID:              op.Id,
			Method:          op.Method,
			PathSegmentID:   op.PathSegmentId,
			QueryParameters: toEntityParameters(op.QueryParameters),
			StatusCodes:     op.StatusCodes,
		}
	}

	return &entity.APIGraph{
		Segments:   segments,
		Edges:      edges,
		Operations: operations,
	}
}

func toEntitySegment(s *sharedpb.PathSegment) entity.PathSegment {
	switch seg := s.Segment.(type) {
	case *sharedpb.PathSegment_Static:
		return entity.PathSegment{
			Static: &entity.StaticSegment{
				ID:   seg.Static.Id,
				Name: seg.Static.Name,
			},
		}
	case *sharedpb.PathSegment_Param:
		return entity.PathSegment{
			Param: &entity.Parameter{
				ID:      seg.Param.Id,
				Name:    seg.Param.Name,
				Type:    toEntityParamType(seg.Param.Type),
				Example: seg.Param.Example,
			},
		}
	default:
		return entity.PathSegment{}
	}
}

func toEntityParameters(params []*sharedpb.Parameter) []entity.Parameter {
	result := make([]entity.Parameter, len(params))
	for i, p := range params {
		result[i] = entity.Parameter{
			ID:      p.Id,
			Name:    p.Name,
			Type:    toEntityParamType(p.Type),
			Example: p.Example,
		}
	}
	return result
}

func toEntityParamType(pt sharedpb.ParameterType) entity.ParameterType {
	switch pt {
	case sharedpb.ParameterType_PARAMETER_TYPE_INTEGER:
		return entity.ParameterTypeInteger
	case sharedpb.ParameterType_PARAMETER_TYPE_UUID:
		return entity.ParameterTypeUUID
	default:
		return entity.ParameterTypeUnspecified
	}
}

func ToProtoGraph(e *entity.APIGraph) *sharedpb.APIGraph {
	if e == nil {
		return nil
	}

	segments := make([]*sharedpb.PathSegment, len(e.Segments))
	for i, s := range e.Segments {
		segments[i] = toProtoSegment(s)
	}

	edges := make([]*sharedpb.Edge, len(e.Edges))
	for i, e := range e.Edges {
		edges[i] = &sharedpb.Edge{
			From: e.From,
			To:   e.To,
		}
	}

	operations := make([]*sharedpb.Operation, len(e.Operations))
	for i, op := range e.Operations {
		operations[i] = &sharedpb.Operation{
			Id:              op.ID,
			Method:          op.Method,
			PathSegmentId:   op.PathSegmentID,
			QueryParameters: toProtoParameters(op.QueryParameters),
			StatusCodes:     op.StatusCodes,
		}
	}

	return &sharedpb.APIGraph{
		Segments:   segments,
		Edges:      edges,
		Operations: operations,
	}
}

func toProtoSegment(s entity.PathSegment) *sharedpb.PathSegment {
	if s.Static != nil {
		return &sharedpb.PathSegment{
			Segment: &sharedpb.PathSegment_Static{
				Static: &sharedpb.StaticSegment{
					Id:   s.Static.ID,
					Name: s.Static.Name,
				},
			},
		}
	}

	if s.Param != nil {
		return &sharedpb.PathSegment{
			Segment: &sharedpb.PathSegment_Param{
				Param: &sharedpb.Parameter{
					Id:      s.Param.ID,
					Name:    s.Param.Name,
					Type:    toProtoParamType(s.Param.Type),
					Example: s.Param.Example,
				},
			},
		}
	}

	return &sharedpb.PathSegment{}
}

func toProtoParameters(params []entity.Parameter) []*sharedpb.Parameter {
	result := make([]*sharedpb.Parameter, len(params))
	for i, p := range params {
		result[i] = &sharedpb.Parameter{
			Id:      p.ID,
			Name:    p.Name,
			Type:    toProtoParamType(p.Type),
			Example: p.Example,
		}
	}
	return result
}

func toProtoParamType(pt entity.ParameterType) sharedpb.ParameterType {
	switch pt {
	case entity.ParameterTypeInteger:
		return sharedpb.ParameterType_PARAMETER_TYPE_INTEGER
	case entity.ParameterTypeUUID:
		return sharedpb.ParameterType_PARAMETER_TYPE_UUID
	default:
		return sharedpb.ParameterType_PARAMETER_TYPE_UNSPECIFIED
	}
}
