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
	"iter"

	"github.com/samber/lo"

	er "outernetcouncil.org/nmts/v1/lib/entityrelationship"
	npb "outernetcouncil.org/nmts/v1/proto"
)

type Node struct {
	entity *npb.Entity
	// kind caches the EK string for the entity; computing it on demand costs a
	// reflection call per lookup, which dominates large-graph traversals.
	// Empty means either "not populated via UpsertEntity" or "entity has no
	// kind set" — GetKind falls back to computing for both, which is cheap for
	// the latter (a nil-kind entity short-circuits before any reflection).
	kind string
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

// GetKind returns the EK string for the underlying NMTS entity, as defined by
// entityrelationship.EntityKindStringFromProto. For nodes built via UpsertEntity the value is
// cached at upsert time; mutating the entity's kind oneof in place afterwards is not supported.
func (n *Node) GetKind() string {
	if n == nil {
		return ""
	}
	if n.kind != "" {
		return n.kind
	}
	// Compute without caching: a lazy write here would race with concurrent readers.
	return er.EntityKindStringFromProto(n.entity)
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
	newEK := er.EntityKindStringFromProto(entity)
	node := g.nodes[entity.GetId()]
	if node != nil {
		if node.GetKind() != newEK {
			return nil, fmt.Errorf("node for ID %s already existed and had a different EK; old EK: %s, new EK: %s", node.GetID(), node.GetKind(), newEK)
		}
	} else {
		node = &Node{}
	}
	node.entity = entity
	node.kind = newEK
	g.nodes[node.GetID()] = node

	nodesOfKind := g.nodesByKind[newEK]
	if nodesOfKind == nil {
		nodesOfKind = map[string]*Node{}
	}
	nodesOfKind[node.GetID()] = node
	g.nodesByKind[newEK] = nodesOfKind

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
	edge, added := g.TryAddRelationship(relationship)
	if !added {
		return nil, fmt.Errorf(
			"already contains an edge of kind %s between %s and %s",
			relationship.GetKind(),
			relationship.GetA(),
			relationship.GetZ(),
		)
	}
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

// TryAddRelationship adds the given relationship to the graph. It returns the new edge and true if
// the relationship was added, or nil and false if the graph already contained an edge representing
// the same relationship (as determined by Edge.Same). Unlike AddRelationship, the duplicate path
// performs no allocations, making it suitable for bulk loads where duplicates are expected.
func (g *Graph) TryAddRelationship(relationship *npb.Relationship) (*Edge, bool) {
	// The probe stays a stack value so the duplicate path allocates nothing.
	probe := Edge{relationship: relationship}
	for _, existing := range g.Edges(probe.GetA(), probe.GetZ()) {
		if probe.Same(existing) {
			return nil, false
		}
	}
	edge := &Edge{
		relationship: relationship,
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
	if edge.GetA() != edge.GetZ() {
		// Self-loops live in a single (A, A) slot; mirroring would store them twice.
		addMappingToEdge(edge.GetZ(), edge.GetA())
	}

	return edge, true
}

// AllNeighbors returns an iterator over the IDs of all nodes adjacent to the node with the given
// ID, paired with the edges connecting them, regardless of the direction of the relationships.
// The graph must not be modified during iteration.
//
// NOTE: Like Neighbors, this is based entirely on the relationships that have been loaded into the
// graph; the yielded IDs may not correspond to actual graph nodes.
func (g *Graph) AllNeighbors(id string) iter.Seq2[string, []*Edge] {
	return func(yield func(string, []*Edge) bool) {
		for neighbor, edges := range g.edges[id] {
			if !yield(neighbor, edges) {
				return
			}
		}
	}
}

// AllNodesOfKind returns an iterator over all nodes of the given kind. The graph must not be
// modified during iteration.
func (g *Graph) AllNodesOfKind(ek string) iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for _, node := range g.nodesByKind[ek] {
			if !yield(node) {
				return
			}
		}
	}
}

// AllEdges returns an iterator over all edges in the graph, yielding each edge exactly once. The
// graph must not be modified during iteration.
func (g *Graph) AllEdges() iter.Seq[*Edge] {
	return func(yield func(*Edge) bool) {
		for x, edgesByNeighbor := range g.edges {
			for y, edges := range edgesByNeighbor {
				if x > y {
					// This pair's edges are also stored under (y, x); yield them there.
					continue
				}
				for _, edge := range edges {
					if !yield(edge) {
						return
					}
				}
			}
		}
	}
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
