package usecase

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/voikin/apim-profile-store/internal/entity"
)

func (u *Usecase) DiffApplicationProfiles(ctx context.Context, applicationID, old, new uuid.UUID) (added []*entity.Operation, removed []*entity.Operation, err error) {
	_, err = u.postgresRepo.GetApplication(ctx, applicationID)
	if err != nil {
		return nil, nil, fmt.Errorf("u.postgresRepo.GetApplication: %w", err)
	}

	oldProfile, err := u.postgresRepo.GetApplicationProfileByID(ctx, old)
	if err != nil {
		return nil, nil, fmt.Errorf("u.postgresRepo.GetApplicationProfileByID: old: %w", err)
	}
	if oldProfile.ApplicationID != applicationID {
		return nil, nil, fmt.Errorf("old profile is related to other application")
	}

	oldAPIGraph, err := u.neo4jRepo.GetAPIGraph(ctx, oldProfile.GraphID)
	if err != nil {
		return nil, nil, fmt.Errorf("u.neo4jRepo.GetAPIGraph: old: %w", err)
	}

	newProfile, err := u.postgresRepo.GetApplicationProfileByID(ctx, new)
	if err != nil {
		return nil, nil, fmt.Errorf("u.postgresRepo.GetApplicationProfileByID: new: %w", err)
	}
	if newProfile.ApplicationID != applicationID {
		return nil, nil, fmt.Errorf("new profile is related to other application")
	}

	newAPIGraph, err := u.neo4jRepo.GetAPIGraph(ctx, newProfile.GraphID)
	if err != nil {
		return nil, nil, fmt.Errorf("u.neo4jRepo.GetAPIGraph: new: %w", err)
	}

	added, removed = MyersDiff(oldAPIGraph.Operations, newAPIGraph.Operations)

	return
}

func EqualOperation(a, b *entity.Operation) bool {
	if a.Method != b.Method || a.PathSegmentID != b.PathSegmentID {
		return false
	}
	if len(a.QueryParameters) != len(b.QueryParameters) || len(a.StatusCodes) != len(b.StatusCodes) {
		return false
	}

	compareParams := func(p1, p2 *entity.Parameter) int {
		if p1.Name < p2.Name {
			return -1
		}
		if p1.Name > p2.Name {
			return 1
		}
		if p1.Type < p2.Type {
			return -1
		}
		if p1.Type > p2.Type {
			return 1
		}
		return 0
	}

	slices.SortFunc(a.QueryParameters, compareParams)
	slices.SortFunc(b.QueryParameters, compareParams)

	for i := range a.QueryParameters {
		if a.QueryParameters[i].Name != b.QueryParameters[i].Name || a.QueryParameters[i].Type != b.QueryParameters[i].Type {
			return false
		}
	}

	slices.Sort(a.StatusCodes)
	slices.Sort(b.StatusCodes)

	for i := range a.StatusCodes {
		if a.StatusCodes[i] != b.StatusCodes[i] {
			return false
		}
	}
	return true
}

func MyersDiff(a, b []*entity.Operation) (added, removed []*entity.Operation) {
	n, m := len(a), len(b)
	max := n + m
	v := map[int]int{1: 0}
	trace := make([]map[int]int, 0, max+1)

	for d := 0; d <= max; d++ {
		vd := make(map[int]int)
		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || (k != d && safeGet(v, k-1) < safeGet(v, k+1)) {
				x = safeGet(v, k+1)
			} else {
				x = safeGet(v, k-1) + 1
			}
			y := x - k

			for x < n && y < m && EqualOperation(a[x], b[y]) {
				x++
				y++
			}
			vd[k] = x

			if x >= n && y >= m {
				trace = append(trace, copyMap(vd))
				goto backtrack
			}
		}
		trace = append(trace, copyMap(vd))
		v = vd
	}

backtrack:
	added = []*entity.Operation{}
	removed = []*entity.Operation{}
	x, y := n, m

	for d := len(trace) - 1; d >= 0; d-- {
		v := trace[d]
		k := x - y

		var prevK int
		if k == -d || (k != d && safeGet(v, k-1) < safeGet(v, k+1)) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		var prevX, prevY int
		if d > 0 {
			prevX = safeGet(trace[d-1], prevK)
			prevY = prevX - prevK
		}

		for x > prevX && y > prevY {
			x--
			y--
		}

		if x == prevX && y > prevY {
			for i := prevY; i < y; i++ {
				added = append([]*entity.Operation{b[i]}, added...)
			}
		} else if y == prevY && x > prevX {
			for i := prevX; i < x; i++ {
				removed = append([]*entity.Operation{a[i]}, removed...)
			}
		}

		x = prevX
		y = prevY
	}

	return
}

func copyMap(m map[int]int) map[int]int {
	cp := make(map[int]int, len(m))
	for k, v := range m {
		cp[k] = v
	}

	return cp
}

func safeGet(m map[int]int, k int) int {
	if val, ok := m[k]; ok {
		return val
	}

	return 0
}
