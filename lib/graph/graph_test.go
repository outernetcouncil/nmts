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
	"cmp"
	"slices"
	"testing"

	set "github.com/deckarep/golang-set/v2"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	npb "outernetcouncil.org/nmts/proto"
)

func mustUnmarshal(t *testing.T, txtPb string, m proto.Message) {
	if err := prototext.Unmarshal([]byte(txtPb), m); err != nil {
		t.Fatalf("unable to unmarshal prototext: %s", txtPb)
	}
}

func mustUnmarshalEntity(t *testing.T, txtPb string) *npb.Entity {
	e := &npb.Entity{}
	mustUnmarshal(t, txtPb, e)
	return e
}

func mustUnmarshalRelationship(t *testing.T, txtPb string) *npb.Relationship {
	r := &npb.Relationship{}
	mustUnmarshal(t, txtPb, r)
	return r
}

func mustUpsertEntity(t *testing.T, g *Graph, e *npb.Entity) *Node {
	node, err := g.UpsertEntity(e)
	if err != nil {
		t.Fatalf("unable to upsert entity %v: %v", e, err)
	}
	return node
}

func mustAddEntities(t *testing.T, g *Graph, entities []string) map[string]*Node {
	nodes := map[string]*Node{}
	for _, e := range entities {
		entity := mustUnmarshalEntity(t, e)
		nodes[entity.GetId()] = mustUpsertEntity(t, g, entity)
	}
	return nodes
}

func mustAddRelationship(t *testing.T, g *Graph, r *npb.Relationship) {
	if _, err := g.AddRelationship(r); err != nil {
		t.Fatalf("unable to upsert relationship %v: %v", r, err)
	}
}

func mustAddRelationships(t *testing.T, g *Graph, relationships []string) {
	for _, r := range relationships {
		mustAddRelationship(t, g, mustUnmarshalRelationship(t, r))
	}
}

func idsSet(nodes []*Node) set.Set[string] {
	return set.NewSet(lo.Map(nodes, func(n *Node, _ int) string {
		return n.ID()
	})...)
}

type upsertEntityTestCase struct {
	desc     string
	entities []string
	wantErr  bool
}

func (tc *upsertEntityTestCase) Run(t *testing.T) {
	g := New()

	var err error = nil
	for _, e := range tc.entities {
		entity := mustUnmarshalEntity(t, e)

		var gotNode *Node
		if gotNode, err = g.UpsertEntity(entity); err != nil {
			break
		}

		wantNode := &Node{
			Entity: entity,
		}
		if *wantNode != *gotNode {
			t.Errorf("unexpected node; want: %v; got: %v", wantNode, gotNode)
		}
	}

	if (err != nil) != tc.wantErr {
		t.Errorf("UpsertEntity; wanted error: %t; got error: %v", tc.wantErr, err)
	}
}

var upsertEntityTestCases = []upsertEntityTestCase{
	{
		desc: "upserting new entities succeeds",
		entities: []string{
			`id: "a" ek_network_node{ name: "nodeA" }`,
			`id: "b" ek_network_node{ name: "nodeB" }`,
			`id: "c" ek_network_node{ name: "nodeC" }`,
		},
	},
	{
		desc: "upserting existing entities updates the node's entity",
		entities: []string{
			`id: "a" ek_network_node{ name: "nodeA1" }`,
			`id: "a" ek_network_node{ name: "nodeA2" }`,
		},
	},
	{
		desc: "upserting existing entities with new EK fails",
		entities: []string{
			`id: "a" ek_network_node{ name: "nodeA1" }`,
			`id: "a" ek_platform{}`,
		},
		wantErr: true,
	},
}

func TestUpsertEntity(t *testing.T) {
	for _, tc := range upsertEntityTestCases {
		t.Run(tc.desc, tc.Run)
	}
}

type nodesOfKindTestCase struct {
	desc     string
	entities []string
	kindIDs  map[string]set.Set[string]
}

func (tc *nodesOfKindTestCase) Run(t *testing.T) {
	g := New()

	mustAddEntities(t, g, tc.entities)

	for kind, wantIDs := range tc.kindIDs {
		gotIDs := idsSet(g.NodesOfKind(kind))
		if !wantIDs.Equal(gotIDs) {
			t.Errorf("unexpected %s IDs; want: %v; got: %v", kind, wantIDs, gotIDs)
		}
	}
}

var nodesOfKindTestCases = []nodesOfKindTestCase{
	{
		desc: "happy path",
		entities: []string{
			`id: "node_a" ek_network_node{}`,
			`id: "node_b" ek_network_node{}`,
			`id: "node_c" ek_network_node{}`,
			`id: "transmitter_a" ek_transmitter{}`,
			`id: "transmitter_b" ek_transmitter{}`,
			`id: "transmitter_c" ek_transmitter{}`,
			`id: "demodulator_a" ek_demodulator{}`,
			`id: "demodulator_b" ek_demodulator{}`,
		},
		kindIDs: map[string]set.Set[string]{
			"EK_NETWORK_NODE": set.NewSet("node_a", "node_b", "node_c"),
			"EK_TRANSMITTER":  set.NewSet("transmitter_a", "transmitter_b", "transmitter_c"),
			"EK_DEMODULATOR":  set.NewSet("demodulator_a", "demodulator_b"),
		},
	},
}

func TestNodesOfKind(t *testing.T) {
	for _, tc := range nodesOfKindTestCases {
		t.Run(tc.desc, tc.Run)
	}
}

type addRelationshipTestCase struct {
	desc          string
	entities      []string
	relationships []string
	wantErr       bool
}

func (tc *addRelationshipTestCase) Run(t *testing.T) {
	g := New()

	nodes := mustAddEntities(t, g, tc.entities)

	var err error = nil
	for _, r := range tc.relationships {
		relationship := mustUnmarshalRelationship(t, r)

		var gotEdge *Edge
		if gotEdge, err = g.AddRelationship(relationship); err != nil {
			break
		}

		wantEdge := &Edge{
			Relationship: relationship,
			A:            nodes[relationship.GetA()],
			Z:            nodes[relationship.GetZ()],
		}
		if *wantEdge != *gotEdge {
			t.Errorf("unexpected edge; want: %v; got: %v", wantEdge, gotEdge)
		}
	}

	if (err != nil) != tc.wantErr {
		t.Errorf("AddRelationship; wanted error: %t; got error: %v", tc.wantErr, err)
	}
}

var addRelationshipTestCases = []addRelationshipTestCase{
	{
		desc: "adding relationship when both nodes exist succeeds",
		entities: []string{
			`id: "node" ek_network_node{}`,
			`id: "interface" ek_interface{}`,
		},
		relationships: []string{
			`a: "node" kind: RK_CONTAINS z: "interface"`,
		},
	},
	{
		desc: "adding relationship when a doesn't exist fails",
		entities: []string{
			`id: "interface" ek_interface{}`,
		},
		relationships: []string{
			`a: "node" kind: RK_CONTAINS z: "interface"`,
		},
		wantErr: true,
	},
	{
		desc: "adding relationship when z doesn't exist fails",
		entities: []string{
			`id: "node" ek_network_node{}`,
		},
		relationships: []string{
			`a: "node" kind: RK_CONTAINS z: "interface"`,
		},
		wantErr: true,
	},
	{
		desc: "same relationship kind in opposite directions are allowed",
		entities: []string{
			`id: "node" ek_network_node{}`,
			`id: "interface" ek_interface{}`,
		},
		relationships: []string{
			`a: "node" kind: RK_CONTAINS z: "interface"`,
			`a: "interface" kind: RK_CONTAINS z: "node"`,
		},
	},
	{
		desc: "same relationship kind in same direction are not allowed",
		entities: []string{
			`id: "node" ek_network_node{}`,
			`id: "interface" ek_interface{}`,
		},
		relationships: []string{
			`a: "node" kind: RK_CONTAINS z: "interface"`,
			`a: "node" kind: RK_CONTAINS z: "interface"`,
		},
		wantErr: true,
	},
}

func TestAddRelationship(t *testing.T) {
	for _, tc := range addRelationshipTestCases {
		t.Run(tc.desc, tc.Run)
	}
}

type graphEntities struct {
	entities      []string
	relationships []string
}

var testGraph = graphEntities{
	entities: []string{
		`id: "agent" ek_sdn_agent{}`,
		`id: "node" ek_network_node{}`,
		`id: "interface" ek_interface{}`,
		`id: "port" ek_port{}`,
		`id: "modulator" ek_modulator{}`,
		`id: "demodulator" ek_demodulator{}`,
		`id: "orphan" ek_network_node{}`,
	},
	relationships: []string{
		`a: "agent" kind: RK_CONTROLS z: "node"`,
		`a: "node" kind: RK_CONTAINS z: "agent"`,
		`a: "node" kind: RK_CONTAINS z: "interface"`,
		`a: "interface" kind: RK_TRAVERSES z: "port"`,
		`a: "port" kind: RK_ORIGINATES z: "modulator"`,
		`a: "port" kind: RK_TERMINATES z: "demodulator"`,
	},
}

type neighborsTestCase struct {
	desc string
	graphEntities
	neighbors map[string]set.Set[string]
}

func (tc *neighborsTestCase) Run(t *testing.T) {
	g := New()

	mustAddEntities(t, g, tc.entities)
	mustAddRelationships(t, g, tc.relationships)

	for nodeID, wantNeighbors := range tc.neighbors {
		gotNeighbors := idsSet(g.Neighbors(nodeID))
		if !wantNeighbors.Equal(gotNeighbors) {
			t.Errorf(
				"Neighbors(%s) returned unexpected neighbors; want: %v; got: %v",
				nodeID,
				wantNeighbors,
				gotNeighbors,
			)
		}
	}
}

var neighborsTestCases = []*neighborsTestCase{
	{
		desc:          "returns all neighbors regardless of relationship direction",
		graphEntities: testGraph,
		neighbors: map[string]set.Set[string]{
			"agent":        set.NewSet("node"),
			"node":         set.NewSet("agent", "interface"),
			"interface":    set.NewSet("node", "port"),
			"port":         set.NewSet("interface", "modulator", "demodulator"),
			"modulator":    set.NewSet("port"),
			"demodulator":  set.NewSet("port"),
			"orphan":       set.NewSet[string](),
			"doesnt_exist": set.NewSet[string](),
		},
	},
}

func TestNeighbors(t *testing.T) {
	for _, tc := range neighborsTestCases {
		t.Run(tc.desc, tc.Run)
	}
}

type idPair struct {
	x string
	y string
}

type edge struct {
	a    string
	kind npb.RK
	z    string
}

func (e *edge) equals(graphEdge *Edge) bool {
	return e.a == graphEdge.A.ID() && e.z == graphEdge.Z.ID() &&
		e.kind == graphEdge.Relationship.GetKind()
}

type edgesTestCase struct {
	desc string
	graphEntities
	edges map[idPair][]*edge
}

var edgesTestCases = []edgesTestCase{
	{
		desc:          "returns all edges between node pairs regardless of direction",
		graphEntities: testGraph,
		edges: map[idPair][]*edge{
			{x: "node", y: "agent"}: {
				{
					a:    "agent",
					kind: npb.RK_RK_CONTROLS,
					z:    "node",
				},
				{
					a:    "node",
					kind: npb.RK_RK_CONTAINS,
					z:    "agent",
				},
			},
			{x: "node", y: "interface"}: {
				{
					a:    "node",
					kind: npb.RK_RK_CONTAINS,
					z:    "interface",
				},
			},
			{x: "interface", y: "port"}: {
				{
					a:    "interface",
					kind: npb.RK_RK_TRAVERSES,
					z:    "port",
				},
			},
			{x: "port", y: "modulator"}: {
				{
					a:    "port",
					kind: npb.RK_RK_ORIGINATES,
					z:    "modulator",
				},
			},
			{x: "port", y: "demodulator"}: {
				{
					a:    "port",
					kind: npb.RK_RK_TERMINATES,
					z:    "demodulator",
				},
			},
			{x: "orphan", y: "node"}:       {},
			{x: "node", y: "doesnt_exist"}: {},
		},
	},
}

func (tc *edgesTestCase) Run(t *testing.T) {
	g := New()

	mustAddEntities(t, g, tc.entities)
	mustAddRelationships(t, g, tc.relationships)

	for ids, wantEdges := range tc.edges {
		validateEdges := func(first, second string) {
			gotEdges := g.Edges(first, second)
			slices.SortFunc(gotEdges, func(a, b *Edge) int {
				return cmp.Or(
					cmp.Compare(a.A.ID(), b.A.ID()),
					cmp.Compare(a.Z.ID(), b.Z.ID()),
					cmp.Compare(a.Relationship.GetKind(), b.Relationship.GetKind()),
				)
			})

			for i, wantEdge := range wantEdges {
				gotEdge := gotEdges[i]
				if !wantEdge.equals(gotEdge) {
					t.Errorf(
						"unexpected edge between %s and %s; want: %v; got: %v",
						first,
						second,
						wantEdge,
						gotEdge,
					)
				}
			}
		}

		validateEdges(ids.x, ids.y)
		validateEdges(ids.y, ids.x)
	}
}

func TestEdges(t *testing.T) {
	for _, tc := range edgesTestCases {
		t.Run(tc.desc, tc.Run)
	}
}
