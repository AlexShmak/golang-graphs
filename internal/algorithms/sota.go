package algorithms

import (
	"context"
	"github.com/AlexShmak/golang-graphs/internal/graph"
)

const (
	UNDECIDED = 0
	LEAF      = 1
	INTERNAL  = 2
)

type SOTASolver struct {
	graph       *graph.Graph
	k           int
	numVertices int
	roles       []int
	parent      map[int]int
	found       bool
	ctx         context.Context
}

func NewSOTASolver(g *graph.Graph, k int, ctx context.Context) *SOTASolver {
	return &SOTASolver{
		graph:       g,
		k:           k,
		numVertices: g.NumVertices,
		ctx:         ctx,
	}
}

func (s *SOTASolver) Solve() bool {
	if s.k > s.numVertices || s.k < 0 {
		return false
	}
	if s.graph.NumVertices > 1 && (s.k < 2 || len(s.graph.Edges) < s.graph.NumVertices-1) {
		return false
	}
	if s.numVertices <= 1 {
		return s.k == s.numVertices
	}
	s.roles = make([]int, s.numVertices)
	s.parent = make(map[int]int)
	s.backtrack(0, 0)
	return s.found
}

func (s *SOTASolver) backtrack(leafCount, internalCount int) {
	select {
	case <-s.ctx.Done():
		return
	default:
	}
	if s.found {
		return
	}
	if leafCount+internalCount == s.numVertices {
		if leafCount == s.k && s.validateSolution() {
			s.found = true
		}
		return
	}
	if leafCount > s.k {
		return
	}
	remaining := s.numVertices - leafCount - internalCount
	if leafCount+remaining < s.k {
		return
	}
	bestV := -1
	maxDeg := -1
	if internalCount > 0 {
		for i := 0; i < s.numVertices; i++ {
			if s.roles[i] == INTERNAL {
				for _, nb := range s.graph.Adj[i] {
					if s.roles[nb] == UNDECIDED {
						if len(s.graph.Adj[nb]) > maxDeg {
							maxDeg = len(s.graph.Adj[nb])
							bestV = nb
						}
					}
				}
			}
		}
	} else {
		for i := 0; i < s.numVertices; i++ {
			if s.roles[i] == UNDECIDED {
				if len(s.graph.Adj[i]) > maxDeg {
					maxDeg = len(s.graph.Adj[i])
					bestV = i
				}
			}
		}
	}
	if bestV == -1 {
		return
	}
	v := bestV
	if internalCount+1 <= s.numVertices-s.k {
		var pCand int
		if internalCount > 0 {
			for _, nb := range s.graph.Adj[v] {
				if s.roles[nb] == INTERNAL {
					pCand = nb
					break
				}
			}
		} else {
			pCand = v
		}
		if pCand >= 0 {
			s.roles[v] = INTERNAL
			if pCand != v {
				s.parent[v] = pCand
			}
			s.backtrack(leafCount, internalCount+1)
			if s.found {
				return
			}
			delete(s.parent, v)
			s.roles[v] = UNDECIDED
		}
	}
	pCand2 := -1
	for _, nb := range s.graph.Adj[v] {
		if s.roles[nb] == INTERNAL {
			pCand2 = nb
			break
		}
	}
	if pCand2 >= 0 {
		s.roles[v] = LEAF
		s.parent[v] = pCand2
		s.backtrack(leafCount+1, internalCount)
		if s.found {
			return
		}
		delete(s.parent, v)
		s.roles[v] = UNDECIDED
	}
}

func (s *SOTASolver) validateSolution() bool {
	if len(s.parent) != s.numVertices-1 {
		return false
	}
	visited := make([]bool, s.numVertices)
	queue := []int{0}
	visited[0] = true
	count := 1
	head := 0
	for head < len(queue) {
		u := queue[head]
		head++
		if p, ok := s.parent[u]; ok && !visited[p] {
			visited[p] = true
			queue = append(queue, p)
			count++
		}
		for child, par := range s.parent {
			if par == u && !visited[child] {
				visited[child] = true
				queue = append(queue, child)
				count++
			}
		}
	}
	return count == s.numVertices
}
