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

	er "outernetcouncil.org/nmts/v1/lib/entityrelationship"
	npb "outernetcouncil.org/nmts/v1/proto"
)

type Node struct {
	entity *npb.Entity
}

func (n *Node) GetEntity() *npb.Entity {
	if n == nil {
		return nil
	}
	return n.entity
}

func (n *Node) GetID() string {
	return n.GetEntity().GetId()
}

// GetKind returns the EK string returned by calling entityrelationship.EntityKindStringFromProto on
// the underlying NMTS entity.
func (n *Node) GetKind() string {
	return er.EntityKindStringFromProto(n.GetEntity())
}

type Edge struct {
	relationship *npb.Relationship
}

func (e *Edge) GetRelationship() *npb.Relationship {
	if e == nil {
		return nil
	}
	return e.relationship
}

func (e *Edge) GetKind() npb.RK {
	return e.GetRelationship().GetKind()
}

func (e *Edge) GetA() string {
	return e.GetRelationship().GetA()
}

func (e *Edge) GetZ() string {
	return e.GetRelationship().GetZ()
}

// Same returns whether this edge represents the same relationship as the other edge.
// NOTE: It does not do a full equality check, it simply returns whether the endpoints and kind are
// the same.
func (e *Edge) Same(other *Edge) bool {
	return e.GetA() == other.GetA() && e.GetKind() == other.GetKind() && e.GetZ() == other.GetZ()
}

// Graph represents a graph of NMTS entites and relationships.
// NOTE: Graph is not thread-safe.
type Graph struct {
	// Maps node ID to the node
	nodes map[string]*Node

	// Maps entity kind string -> node ID for all nodes of that kind -> the Node
	nodesByKind map[string]map[string]*Node

	// Maps node ID -> adjacent node ID -> all edges connecting them, regardless of direction
	edges map[string]map[string][]*Edge
}

func New() *Graph {
	return &Graph{
		nodes:       map[string]*Node{},
		nodesByKind: map[string]map[string]*Node{},
		edges:       map[string]map[string][]*Edge{},
	}
}

func (g *Graph) Node(id string) *Node {
	return g.nodes[id]
}

func (g *Graph) NodesOfKind(ek string) []*Node {
	return lo.Values(g.nodesByKind[ek])
}

func (g *Graph) UpsertEntity(entity *npb.Entity) (*Node, error) {
	node := g.nodes[entity.GetId()]
	if node != nil {
		if newEK := er.EntityKindStringFromProto(entity); node.GetKind() != newEK {
			return nil, fmt.Errorf("node for ID %s already existed and had a different EK; old EK: %s, new EK: %s", node.GetID(), node.GetKind(), newEK)
		}
	} else {
		node = &Node{}
	}
	node.entity = entity
	g.nodes[node.GetID()] = node

	kind := node.GetKind()
	nodesOfKind := g.nodesByKind[kind]
	if nodesOfKind == nil {
		nodesOfKind = map[string]*Node{}
	}
	nodesOfKind[node.GetID()] = node
	g.nodesByKind[kind] = nodesOfKind

	return node, nil
}

func (g *Graph) RemoveEntity(id string) error {
	node := g.nodes[id]
	if node == nil {
		return fmt.Errorf("no corresponding node found")
	}

	delete(g.nodes, id)

	nodesOfKind := g.nodesByKind[node.GetKind()]
	delete(nodesOfKind, node.GetID())
	if len(nodesOfKind) == 0 {
		delete(g.nodesByKind, node.GetKind())
	}
	return nil
}

func (g *Graph) AddRelationship(relationship *npb.Relationship) (*Edge, error) {
	edge := &Edge{
		relationship: relationship,
	}

	edges := g.Edges(edge.GetA(), edge.GetZ())
	if slices.ContainsFunc(edges, edge.Same) {
		return nil, fmt.Errorf(
			"already contains an edge of kind %s between %s and %s",
			edge.GetKind(),
			edge.GetA(),
			edge.GetZ(),
		)
	}

	addMappingToEdge := func(x, y string) {
		edgesByNeighbor := g.edges[x]
		if edgesByNeighbor == nil {
			edgesByNeighbor = map[string][]*Edge{}
		}
		edgesToY := edgesByNeighbor[y]
		if edgesToY == nil {
			edgesToY = []*Edge{}
		}
		edgesToY = append(edgesToY, edge)
		edgesByNeighbor[y] = edgesToY
		g.edges[x] = edgesByNeighbor
	}
	addMappingToEdge(edge.GetA(), edge.GetZ())
	addMappingToEdge(edge.GetZ(), edge.GetA())

	return edge, nil
}

func (g *Graph) RemoveRelationship(relationship *npb.Relationship) error {
	edgeToRemove := &Edge{
		relationship: relationship,
	}

	removedAnEdge := false
	removeMappingToEdge := func(x, y string) {
		edgesByNeighbor := g.edges[x]
		if edgesByNeighbor == nil {
			return
		}
		edgesToY := edgesByNeighbor[y]
		if edgesToY == nil {
			return
		}
		filteredEdgesToY := lo.Filter(edgesToY, func(e *Edge, _ int) bool {
			return !edgeToRemove.Same(e)
		})
		removedAnEdge = len(filteredEdgesToY) < len(edgesToY)
		if len(filteredEdgesToY) == 0 {
			delete(edgesByNeighbor, y)
			if len(edgesByNeighbor) == 0 {
				delete(g.edges, x)
			}
		} else {
			edgesByNeighbor[y] = filteredEdgesToY
		}
	}
	removeMappingToEdge(edgeToRemove.GetA(), edgeToRemove.GetZ())
	removeMappingToEdge(edgeToRemove.GetZ(), edgeToRemove.GetA())

	if !removedAnEdge {
		return fmt.Errorf("no corresponding edge found")
	}
	return nil
}

// Neighbors returns the IDs of all nodes that are adjacent to the node with the given ID. It
// returns all adjacent nodes, regardless of the direction of the relationships connecting them.
//
// NOTE: This is based entirely on the relationships that have been loaded into the graph. There are
// no guarantees that the returned node IDs all correspond to actual graph nodes.
func (g *Graph) Neighbors(id string) []string {
	edgesByNeighbor := g.edges[id]
	return lo.Keys(edgesByNeighbor)
}

// Edges returns all edges connecting the nodes with the given IDs, regardless of the order of the
// IDs and the direction of the edges (ie. calling Edges("a", "z") and Edges("z", "a") will return
// the same edges, though not necessarily in the same order).
func (g *Graph) Edges(x, y string) []*Edge {
	xEdgesByNeighbor := g.edges[x]
	if xEdgesByNeighbor == nil {
		return nil
	}
	return xEdgesByNeighbor[y]
}
