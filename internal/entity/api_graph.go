package entity

type ParameterType int32

const (
	ParameterTypeUnspecified ParameterType = 0
	ParameterTypeInteger     ParameterType = 1
	ParameterTypeUUID        ParameterType = 2
)

type PathSegment struct {
	Static *StaticSegment
	Param  *Parameter
}

type StaticSegment struct {
	ID   string
	Name string
}

type Parameter struct {
	ID      string
	Name    string
	Type    ParameterType
	Example string
}

type Edge struct {
	From string
	To   string
}

type Operation struct {
	ID              string
	Method          string
	PathSegmentID   string
	QueryParameters []*Parameter
	StatusCodes     []int32
}

type Transition struct {
	From string
	To   string
}

type APIGraph struct {
	Segments    []*PathSegment
	Edges       []*Edge
	Operations  []*Operation
	Transitions []*Transition
}
