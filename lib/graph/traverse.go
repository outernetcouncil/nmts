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

type stack[T comparable] []T

func (s *stack[T]) Len() int { return len(*s) }
func (s *stack[T]) Push(n T) { *s = append(*s, n) }
func (s *stack[T]) Pop() T {
	n := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return n
}

type DepthFirst struct {
	// Visit is called once for each Node that is visited.
	Visit func(*Graph, string)

	// Traverse is called to determine whether the given Edge should be traversed from the given Node.
	Traverse func(*Graph, string, *Edge) bool

	visited set.Set[string]
}

func (df *DepthFirst) reset() {
	df.visited = set.NewSet[string]()
}

// Walk walks the entire graph g starting from the from Node, respecting the Traverse function if
// provided. until is called after each visit, and if it returns true, the walk is terminated
// early.
func (df *DepthFirst) Walk(g *Graph, from string, until func(string) bool) {
	df.reset()

	s := &stack[string]{}
	s.Push(from)
	for s.Len() != 0 {
		visiting := s.Pop()
		if df.visited.Contains(visiting) {
			continue
		}
		df.visited.Add(visiting)
		if df.Visit != nil {
			df.Visit(g, visiting)
		}
		if until != nil && until(visiting) {
			return
		}
		toVisits := g.Neighbors(visiting)
		for _, toVisit := range toVisits {
			edges := g.Edges(visiting, toVisit)
			for _, edge := range edges {
				if df.Traverse != nil && df.Traverse(g, visiting, edge) {
					s.Push(toVisit)
					break
				}
			}
		}
	}
}

type TraverseOpt func(*Graph, string, *Edge) bool

// Traverse is a helper for creating Traverse functions that can be provided to a traversal object
// above.
func Traverse(opts ...TraverseOpt) func(*Graph, string, *Edge) bool {
	return func(g *Graph, from string, edge *Edge) bool {
		for _, opt := range opts {
			if opt(g, from, edge) {
				return true
			}
		}
		return false
	}
}

// Edges will include traversal of all edges that have the given a and z node kinds and edge
// relationship kind, regardless of whether the direction of traversal matches the direction of the
// relationship.
//
// NOTE: If an edge corresponds to a node that isn't loaded into the graph, it will not be
// traversed.
func Edges(aKind string, relationshipKind npb.RK, zKind string) TraverseOpt {
	return func(g *Graph, _ string, edge *Edge) bool {
		a, z := g.Node(edge.A()), g.Node(edge.Z())
		if a == nil || z == nil {
			return false
		}

		return a.Kind() == aKind && edge.Kind() == relationshipKind && z.Kind() == zKind
	}
}

// EdgesFrom will include traversal of all edges that have the given a and z node kinds and edge
// relationship kind, but only if the edge is being traversed from a node with the given fromKind.
//
// NOTE: If an edge corresponds to a node that isn't loaded into the graph, it will not be
// traversed.
func EdgesFrom(fromKind, aKind string, relationshipKind npb.RK, zKind string) TraverseOpt {
	return func(g *Graph, fromID string, edge *Edge) bool {
		from, a, z := g.Node(fromID), g.Node(edge.A()), g.Node(edge.Z())
		if from == nil || a == nil || z == nil {
			return false
		}

		return from.Kind() == fromKind && a.Kind() == aKind && edge.Kind() == relationshipKind && z.Kind() == zKind
	}
}
