// Copyright (c) Outernet Council and Contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package graph

import (
	set "github.com/deckarep/golang-set/v2"

	npb "outernetcouncil.org/nmts/proto"
)

type stack []*Node

func (s *stack) Len() int     { return len(*s) }
func (s *stack) Push(n *Node) { *s = append(*s, n) }
func (s *stack) Pop() *Node {
	n := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return n
}

type DepthFirst struct {
	// Visit is called once for each Node that is visited.
	Visit func(*Node)

	// Traverse is called to determine whether the given Edge should be traversed from the given Node.
	Traverse func(*Node, *Edge) bool

	visited set.Set[*Node]
}

func (df *DepthFirst) reset() {
	df.visited = set.NewSet[*Node]()
}

// Walk walks the entire graph g starting from the from Node, respecting the Traverse function if
// provided. until is called after each visit, and if it returns true, the walk is terminated
// early.
func (df *DepthFirst) Walk(g *Graph, from *Node, until func(*Node) bool) {
	df.reset()

	s := &stack{}
	s.Push(from)
	for s.Len() != 0 {
		visiting := s.Pop()
		if df.visited.Contains(visiting) {
			continue
		}
		df.visited.Add(visiting)
		if df.Visit != nil {
			df.Visit(visiting)
		}
		if until != nil && until(visiting) {
			return
		}
		toVisits := g.Neighbors(visiting.ID())
		for _, toVisit := range toVisits {
			edges := g.Edges(visiting.ID(), toVisit.ID())
			for _, edge := range edges {
				if df.Traverse != nil && df.Traverse(visiting, edge) {
					s.Push(toVisit)
					break
				}
			}
		}
	}
}

type TraverseOpt func(*Node, *Edge) bool

// Traverse is a helper for creating Traverse functions that can be provided to a traversal object
// above.
func Traverse(opts ...TraverseOpt) func(*Node, *Edge) bool {
	return func(from *Node, edge *Edge) bool {
		for _, opt := range opts {
			if opt(from, edge) {
				return true
			}
		}
		return false
	}
}

// Edges will include traversal of all edges that have the given a and z node kinds and edge
// relationship kind, regardless of whether the direction of traversal matches the direction of the
// relationship.
func Edges(aKind string, relationshipKind npb.RK, zKind string) TraverseOpt {
	return func(from *Node, edge *Edge) bool {
		return edge.A.Kind() == aKind && edge.Kind() == relationshipKind && edge.Z.Kind() == zKind
	}
}

// EdgesFrom will include traversal of all edges that have the given a and z node kinds and edge
// relationship kind, but only if the edge is being traversed from a node with the given fromKind.
func EdgesFrom(fromKind, aKind string, relationshipKind npb.RK, zKind string) TraverseOpt {
	return func(from *Node, edge *Edge) bool {
		return from.Kind() == fromKind && edge.A.Kind() == aKind && edge.Kind() == relationshipKind && edge.Z.Kind() == zKind
	}
}
