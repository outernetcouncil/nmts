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
	"fmt"
	"slices"

	"github.com/samber/lo"

	er "outernetcouncil.org/nmts/lib/entityrelationship"
	npb "outernetcouncil.org/nmts/proto"
)

type Node struct {
	Entity *npb.Entity
}

func (n *Node) ID() string {
	return n.Entity.GetId()
}

// Kind returns the EK string returned by calling entityrelationship.EntityKindStringFromProto on
// the underlying NMTS entity.
func (n *Node) Kind() string {
	return er.EntityKindStringFromProto(n.Entity)
}

type Edge struct {
	Relationship *npb.Relationship
	A            *Node
	Z            *Node
}

func (e *Edge) Kind() npb.RK {
	return e.Relationship.GetKind()
}

// Same returns whether this edge represents the same relationship as the other edge.
// NOTE: It does not do a full equality check, it simply returns whether the endpoints and kind are
// the same.
func (e *Edge) Same(other *Edge) bool {
	return e.A.ID() == other.A.ID() && e.Kind() == other.Kind() && e.Z.ID() == other.Z.ID()
}

// Graph represents a graph of NMTS entites and relationships.
// NOTE: Graph is not thread-safe.
type Graph struct {
	// Maps node ID to the node
	nodes map[string]*Node

	// Maps entity kind string to all nodes of that kind
	nodesByKind map[string][]*Node

	// Maps node -> adjacent node -> all edges connecting them, regardless of direction
	edges map[*Node]map[*Node][]*Edge
}

func New() *Graph {
	return &Graph{
		nodes:       map[string]*Node{},
		nodesByKind: map[string][]*Node{},
		edges:       map[*Node]map[*Node][]*Edge{},
	}
}

func (g *Graph) NodesOfKind(ek string) []*Node {
	return g.nodesByKind[ek]
}

func (g *Graph) UpsertEntity(entity *npb.Entity) (*Node, error) {
	node := g.nodes[entity.GetId()]
	if node != nil {
		if newEK := er.EntityKindStringFromProto(entity); node.Kind() != newEK {
			return nil, fmt.Errorf("node for ID %s already existed and had a different EK; old EK: %s, new EK: %s", node.ID(), node.Kind(), newEK)
		}
	} else {
		node = &Node{}
	}
	node.Entity = entity
	g.nodes[node.ID()] = node

	kind := node.Kind()
	nodesOfKind := g.nodesByKind[kind]
	if nodesOfKind == nil {
		nodesOfKind = []*Node{}
	}
	nodesOfKind = append(nodesOfKind, node)
	g.nodesByKind[kind] = nodesOfKind

	return node, nil
}

func (g *Graph) AddRelationship(relationship *npb.Relationship) (*Edge, error) {
	a := g.nodes[relationship.GetA()]
	if a == nil {
		return nil, fmt.Errorf("a ID %s doesn't correspond to any node", relationship.GetA())
	}

	z := g.nodes[relationship.GetZ()]
	if z == nil {
		return nil, fmt.Errorf("z ID %s doesn't correspond to any node", relationship.GetZ())
	}

	edge := &Edge{
		Relationship: relationship,
		A:            a,
		Z:            z,
	}

	edges := g.Edges(a.ID(), z.ID())
	if slices.ContainsFunc(edges, edge.Same) {
		return nil, fmt.Errorf(
			"already contains an edge of kind %s between %s and %s",
			edge.Kind(),
			edge.A.ID(),
			edge.Z.ID(),
		)
	}

	addMappingToEdge := func(x, y *Node) {
		edgesByNeighbor := g.edges[x]
		if edgesByNeighbor == nil {
			edgesByNeighbor = map[*Node][]*Edge{}
		}
		edgesToY := edgesByNeighbor[y]
		if edgesToY == nil {
			edgesToY = []*Edge{}
		}
		edgesToY = append(edgesToY, edge)
		edgesByNeighbor[y] = edgesToY
		g.edges[x] = edgesByNeighbor
	}
	addMappingToEdge(a, z)
	addMappingToEdge(z, a)

	return edge, nil
}

// Neighbors returns all nodes that are adjacent to the node with the given ID. It returns all adjacent
// nodes, regardless of the direction of the relationships connecting them.
func (g *Graph) Neighbors(id string) []*Node {
	node := g.nodes[id]
	edgesByNeighbor := g.edges[node]
	return lo.Keys(edgesByNeighbor)
}

// Edges returns all edges connecting the nodes with the given IDs, regardless of the order of the
// IDs and the direction of the edges (ie. calling Edges("a", "z") and Edges("z", "a") will return
// the same edges, though not necessarily in the same order).
func (g *Graph) Edges(x, y string) []*Edge {
	xNode, yNode := g.nodes[x], g.nodes[y]
	if xNode == nil || yNode == nil {
		return nil
	}
	xEdgesByNeighbor := g.edges[xNode]
	if xEdgesByNeighbor == nil {
		return nil
	}
	return xEdgesByNeighbor[yNode]
}
