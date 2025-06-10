package algorithms

import (
	"context"
	"github.com/AlexShmak/golang-graphs/internal/graph"
)

type BaselineSolver struct {
	graph *graph.Graph
	k     int
	found bool
	ctx   context.Context
}

func NewBaselineSolver(g *graph.Graph, k int, ctx context.Context) *BaselineSolver {
	return &BaselineSolver{graph: g, k: k, ctx: ctx}
}

func (s *BaselineSolver) Solve() bool {
	if s.graph.NumVertices > 1 && len(s.graph.Edges) < s.graph.NumVertices-1 {
		return false
	}
	d := graph.NewDSU(s.graph.NumVertices)
	s.backtrack(0, [][2]int{}, d)
	return s.found
}

func (s *BaselineSolver) backtrack(edgeIndex int, currentEdges [][2]int, ds *graph.DSU) {
	select {
	case <-s.ctx.Done():
		return
	default:
	}
	if s.found {
		return
	}
	if len(currentEdges) == s.graph.NumVertices-1 {
		if s.countLeaves(currentEdges) == s.k {
			s.found = true
		}
		return
	}
	if edgeIndex >= len(s.graph.Edges) {
		return
	}
	edge := s.graph.Edges[edgeIndex]
	u, v := edge[0], edge[1]
	if ds.Find(u) != ds.Find(v) {
		parentCopy := make([]int, len(ds.Parent()))
		copy(parentCopy, ds.Parent())
		ds.Union(u, v)
		s.backtrack(edgeIndex+1, append(currentEdges, edge), ds)
		if s.found {
			return
		}
		ds.SetParent(parentCopy)
	}
	s.backtrack(edgeIndex+1, currentEdges, ds)
}

func (s *BaselineSolver) countLeaves(edges [][2]int) int {
	n := s.graph.NumVertices
	if n <= 1 {
		return n
	}
	if len(edges) == 0 {
		return 0
	}
	degrees := make(map[int]int)
	all := make(map[int]bool)
	for _, e := range edges {
		degrees[e[0]]++
		degrees[e[1]]++
		all[e[0]] = true
		all[e[1]] = true
	}
	if len(all) != n {
		return 0
	}
	leafCount := 0
	for _, deg := range degrees {
		if deg == 1 {
			leafCount++
		}
	}
	return leafCount
}
