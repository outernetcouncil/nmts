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
	gcmp "github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	er "outernetcouncil.org/nmts/v1/lib/entityrelationship"
	npb "outernetcouncil.org/nmts/v1/proto"
	logicalpb "outernetcouncil.org/nmts/v1/proto/ek/logical"
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

func mustUpsertEntities(t *testing.T, g *Graph, entities []string) {
	for _, e := range entities {
		mustUpsertEntity(t, g, mustUnmarshalEntity(t, e))
	}
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

func mustRemoveEntities(t *testing.T, g *Graph, ids []string) {
	for _, id := range ids {
		if err := g.RemoveEntity(id); err != nil {
			t.Fatalf("unable to remove entity with ID: %s: %v", id, err)
		}
	}
}

func mustRemoveRelationships(t *testing.T, g *Graph, relationships []string) {
	for _, relationshipTxtPb := range relationships {
		relationship := mustUnmarshalRelationship(t, relationshipTxtPb)
		if err := g.RemoveRelationship(relationship); err != nil {
			t.Fatalf("unable to remove relationship %v: %v", relationship, err)
		}
	}
}

func idsSet(nodes []*Node) set.Set[string] {
	return set.NewSet(lo.Map(nodes, func(n *Node, _ int) string {
		return n.GetID()
	})...)
}

type nodeTestCase struct {
	desc       string
	node       *Node
	wantEntity *npb.Entity
	wantID     string
	wantKind   string
}

func (tc *nodeTestCase) run(t *testing.T) {
	gotEntity := tc.node.GetEntity()
	if diff := gcmp.Diff(tc.wantEntity, gotEntity, protocmp.Transform()); diff != "" {
		t.Errorf("unexpected entity (-want +got): %s", diff)
	}
	if gotID := tc.node.GetID(); tc.wantID != gotID {
		t.Errorf("unexpected ID; want: %v; got: %v", tc.wantID, gotID)
	}
	if gotKind := tc.node.GetKind(); tc.wantKind != gotKind {
		t.Errorf("unexpected kind; want: %v; got: %v", tc.wantKind, gotKind)
	}
}

var nodeTestCases = []nodeTestCase{
	{
		desc: "non-nil node",
		node: &Node{
			entity: &npb.Entity{
				Id: "node_id",
				Kind: &npb.Entity_EkNetworkNode{
					EkNetworkNode: &logicalpb.NetworkNode{},
				},
			},
		},
		wantID:   "node_id",
		wantKind: "EK_NETWORK_NODE",
		wantEntity: &npb.Entity{
			Id: "node_id",
			Kind: &npb.Entity_EkNetworkNode{
				EkNetworkNode: &logicalpb.NetworkNode{},
			},
		},
	},
	{
		desc:       "nil node",
		node:       nil,
		wantID:     "",
		wantKind:   "",
		wantEntity: nil,
	},
}

func TestNode(t *testing.T) {
	for _, tc := range nodeTestCases {
		t.Run(tc.desc, tc.run)
	}
}

type edgeTestCase struct {
	desc             string
	edge             *Edge
	wantRelationship *npb.Relationship
	wantKind         npb.RK
	wantA            string
	wantZ            string
}

func (tc *edgeTestCase) run(t *testing.T) {
	gotRelationship := tc.edge.GetRelationship()
	if diff := gcmp.Diff(tc.wantRelationship, gotRelationship, protocmp.Transform()); diff != "" {
		t.Errorf("unexpected relationship (-want +got): %s", diff)
	}
	if gotKind := tc.edge.GetKind(); tc.wantKind != gotKind {
		t.Errorf("unexpected kind; want: %v; got: %v", tc.wantKind, gotKind)
	}
	if gotA := tc.edge.GetA(); tc.wantA != gotA {
		t.Errorf("unexpected A; want: %v; got: %v", tc.wantA, gotA)
	}
	if gotZ := tc.edge.GetZ(); tc.wantZ != gotZ {
		t.Errorf("unexpected Z; want: %v; got: %v", tc.wantZ, gotZ)
	}
}

var edgeTestCases = []edgeTestCase{
	{
		desc: "non-nil edge",
		edge: &Edge{
			relationship: &npb.Relationship{
				Kind: npb.RK_RK_AGGREGATES,
				A:    "node_a",
				Z:    "node_z",
			},
		},
		wantRelationship: &npb.Relationship{
			Kind: npb.RK_RK_AGGREGATES,
			A:    "node_a",
			Z:    "node_z",
		},
		wantKind: npb.RK_RK_AGGREGATES,
		wantA:    "node_a",
		wantZ:    "node_z",
	},
	{
		desc:             "nil edge",
		edge:             nil,
		wantRelationship: nil,
		wantKind:         npb.RK_RK_UNSPECIFIED,
		wantA:            "",
		wantZ:            "",
	},
}

func TestEdge(t *testing.T) {
	for _, tc := range edgeTestCases {
		t.Run(tc.desc, tc.run)
	}
}

type nodesOfKindTestCase struct {
	desc           string
	entities       []string
	entityRemovals []string
	kindIDs        map[string]set.Set[string]
}

func (tc *nodesOfKindTestCase) Run(t *testing.T) {
	g := New()

	mustUpsertEntities(t, g, tc.entities)
	mustRemoveEntities(t, g, tc.entityRemovals)

	for kind, wantIDs := range tc.kindIDs {
		gotIDs := idsSet(g.NodesOfKind(kind))
		if !wantIDs.Equal(gotIDs) {
			t.Errorf("unexpected %s IDs; want: %v; got: %v", kind, wantIDs, gotIDs)
		}
	}
}

var nodesOfKindTestCases = []nodesOfKindTestCase{
	{
		desc: "returns nodes of given kind",
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
	{
		desc: "respects entity removals",
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
		entityRemovals: []string{"node_b", "transmitter_c", "demodulator_a", "demodulator_b"},
		kindIDs: map[string]set.Set[string]{
			"EK_NETWORK_NODE": set.NewSet("node_a", "node_c"),
			"EK_TRANSMITTER":  set.NewSet("transmitter_a", "transmitter_b"),
			"EK_DEMODULATOR":  set.NewSet[string](),
		},
	},
}

func TestNodesOfKind(t *testing.T) {
	for _, tc := range nodesOfKindTestCases {
		t.Run(tc.desc, tc.Run)
	}
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
			entity: entity,
			kind:   er.EntityKindStringFromProto(entity),
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

type removeEntityTestCase struct {
	desc          string
	entities      []string
	relationships []string
	remove        string
	wantErr       bool
}

func (tc *removeEntityTestCase) Run(t *testing.T) {
	g := New()

	mustUpsertEntities(t, g, tc.entities)
	mustAddRelationships(t, g, tc.relationships)

	err := g.RemoveEntity(tc.remove)
	if gotErr := err != nil; tc.wantErr != gotErr {
		t.Errorf("RemoveEntity; wanted error: %t; got error: %v", tc.wantErr, err)
	}
}

var removeEntityTestCases = []removeEntityTestCase{
	{
		desc: "removing existing entity succeeds",
		entities: []string{
			`id: "node" ek_network_node{}`,
		},
		remove: "node",
	},
	{
		desc:    "removing non-existent entity fails",
		remove:  "doesnt_exist",
		wantErr: true,
	},
	{
		desc: "removing entity that has edges fails",
		entities: []string{
			`id: "node" ek_network_node{}`,
			`id: "interface" ek_interface{}`,
		},
		relationships: []string{
			`a: "node" kind: RK_CONTAINS z: "interface"`,
		},
		remove:  "a",
		wantErr: true,
	},
}

func TestRemoveEntity(t *testing.T) {
	for _, tc := range removeEntityTestCases {
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

	mustUpsertEntities(t, g, tc.entities)

	var err error = nil
	for _, r := range tc.relationships {
		relationship := mustUnmarshalRelationship(t, r)

		var gotEdge *Edge
		if gotEdge, err = g.AddRelationship(relationship); err != nil {
			break
		}

		wantEdge := &Edge{
			relationship: relationship,
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
		desc: "adding relationship when a doesn't exist succeeds",
		entities: []string{
			`id: "interface" ek_interface{}`,
		},
		relationships: []string{
			`a: "node" kind: RK_CONTAINS z: "interface"`,
		},
	},
	{
		desc: "adding relationship when z doesn't exist succeeds",
		entities: []string{
			`id: "node" ek_network_node{}`,
		},
		relationships: []string{
			`a: "node" kind: RK_CONTAINS z: "interface"`,
		},
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

type removeRelationshipTestCase struct {
	desc          string
	entities      []string
	relationships []string
	remove        string
	wantErr       bool
}

func (tc *removeRelationshipTestCase) Run(t *testing.T) {
	g := New()

	mustUpsertEntities(t, g, tc.entities)
	mustAddRelationships(t, g, tc.relationships)

	toRemove := mustUnmarshalRelationship(t, tc.remove)
	err := g.RemoveRelationship(toRemove)
	if gotErr := err != nil; tc.wantErr != gotErr {
		t.Errorf("RemoveRelationship; wanted error: %t; got error: %v", tc.wantErr, err)
	}
}

var removeRelationshipTestCases = []removeRelationshipTestCase{
	{
		desc: "removing existing relationship succeeds",
		entities: []string{
			`id: "node" ek_network_node{}`,
			`id: "platform" ek_platform{}`,
		},
		relationships: []string{
			`a: "node" kind: RK_CONTAINS z: "platform"`,
		},
		remove: `a: "node" kind: RK_CONTAINS z: "platform"`,
	},
	{
		desc: "removing self-loop relationship succeeds",
		entities: []string{
			`id: "node" ek_network_node{}`,
		},
		relationships: []string{
			`a: "node" kind: RK_CONTAINS z: "node"`,
		},
		remove: `a: "node" kind: RK_CONTAINS z: "node"`,
	},
	{
		desc: "removing non-existent relationship fails",
		entities: []string{
			`id: "node" ek_network_node{}`,
			`id: "platform" ek_platform{}`,
		},
		remove:  `a: "node" kind: RK_CONTAINS z: "platform"`,
		wantErr: true,
	},
	{
		desc: "removing relationship where a doesn't correspond to any node fails",
		entities: []string{
			`id: "platform" ek_platform{}`,
		},
		remove:  `a: "node" kind: RK_CONTAINS z: "platform"`,
		wantErr: true,
	},
	{
		desc: "removing relationship where z doesn't correspond to any node fails",
		entities: []string{
			`id: "node" ek_network_node{}`,
		},
		remove:  `a: "node" kind: RK_CONTAINS z: "platform"`,
		wantErr: true,
	},
}

func TestRemoveRelationship(t *testing.T) {
	for _, tc := range removeRelationshipTestCases {
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
	relationshipRemovals []string
	neighbors            map[string]set.Set[string]
}

func (tc *neighborsTestCase) Run(t *testing.T) {
	g := New()

	mustUpsertEntities(t, g, tc.entities)
	mustAddRelationships(t, g, tc.relationships)
	mustRemoveRelationships(t, g, tc.relationshipRemovals)

	for nodeID, wantNeighbors := range tc.neighbors {
		gotNeighbors := set.NewSet(g.Neighbors(nodeID)...)
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
	{
		desc:          "respects relationship removals",
		graphEntities: testGraph,
		relationshipRemovals: []string{
			`a: "interface" kind: RK_TRAVERSES z: "port"`,
			`a: "port" kind: RK_TERMINATES z: "demodulator"`,
		},
		neighbors: map[string]set.Set[string]{
			"agent":       set.NewSet("node"),
			"node":        set.NewSet("agent", "interface"),
			"interface":   set.NewSet("node"),
			"port":        set.NewSet("modulator"),
			"modulator":   set.NewSet("port"),
			"demodulator": set.NewSet[string](),
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
	return e.a == graphEdge.GetA() && e.z == graphEdge.GetZ() &&
		e.kind == graphEdge.relationship.GetKind()
}

type edgesTestCase struct {
	desc string
	graphEntities
	relationshipRemovals []string
	edges                map[idPair][]*edge
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
	{
		desc:          "respects relationship removals",
		graphEntities: testGraph,
		relationshipRemovals: []string{
			`a: "interface" kind: RK_TRAVERSES z: "port"`,
			`a: "port" kind: RK_TERMINATES z: "demodulator"`,
		},
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
			{x: "interface", y: "port"}: {},
			{x: "port", y: "modulator"}: {
				{
					a:    "port",
					kind: npb.RK_RK_ORIGINATES,
					z:    "modulator",
				},
			},
			{x: "port", y: "demodulator"}: {},
		},
	},
}

func (tc *edgesTestCase) Run(t *testing.T) {
	g := New()

	mustUpsertEntities(t, g, tc.entities)
	mustAddRelationships(t, g, tc.relationships)
	mustRemoveRelationships(t, g, tc.relationshipRemovals)

	for ids, wantEdges := range tc.edges {
		validateEdges := func(first, second string) {
			gotEdges := g.Edges(first, second)
			slices.SortFunc(gotEdges, func(a, b *Edge) int {
				return cmp.Or(
					cmp.Compare(a.GetA(), b.GetA()),
					cmp.Compare(a.GetZ(), b.GetZ()),
					cmp.Compare(a.relationship.GetKind(), b.relationship.GetKind()),
				)
			})

			if len(wantEdges) != len(gotEdges) {
				t.Errorf("unexpected number of edges; want: %d; got: %d", len(wantEdges), len(gotEdges))
			}
			for i := 0; i < min(len(wantEdges), len(gotEdges)); i++ {
				wantEdge := wantEdges[i]
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

func TestTryAddRelationship(t *testing.T) {
	g := New()
	mustUpsertEntities(t, g, []string{
		`id: "node" ek_network_node{}`,
		`id: "interface" ek_interface{}`,
	})

	r := mustUnmarshalRelationship(t, `a: "node" kind: RK_CONTAINS z: "interface"`)
	edge, added := g.TryAddRelationship(r)
	if !added {
		t.Fatalf("TryAddRelationship(%v) returned added=false for a new relationship", r)
	}
	if edge.GetRelationship() != r {
		t.Errorf("TryAddRelationship returned edge with unexpected relationship; want: %v; got: %v", r, edge.GetRelationship())
	}

	duplicate := mustUnmarshalRelationship(t, `a: "node" kind: RK_CONTAINS z: "interface"`)
	edge, added = g.TryAddRelationship(duplicate)
	if added {
		t.Errorf("TryAddRelationship(%v) returned added=true for a duplicate relationship", duplicate)
	}
	if edge != nil {
		t.Errorf("TryAddRelationship returned non-nil edge for a duplicate relationship: %v", edge)
	}
	if gotEdges := g.Edges("node", "interface"); len(gotEdges) != 1 {
		t.Errorf("unexpected number of edges after duplicate TryAddRelationship; want: 1; got: %d", len(gotEdges))
	}
}

func TestAllNeighbors(t *testing.T) {
	g := New()
	mustUpsertEntities(t, g, testGraph.entities)
	mustAddRelationships(t, g, testGraph.relationships)

	wantNeighbors := map[string]set.Set[string]{
		"port":         set.NewSet("interface", "modulator", "demodulator"),
		"orphan":       set.NewSet[string](),
		"doesnt_exist": set.NewSet[string](),
	}
	for nodeID, want := range wantNeighbors {
		got := set.NewSet[string]()
		for neighbor, edges := range g.AllNeighbors(nodeID) {
			got.Add(neighbor)
			if len(edges) == 0 {
				t.Errorf("AllNeighbors(%s) yielded neighbor %s with no edges", nodeID, neighbor)
			}
			wantEdges := g.Edges(nodeID, neighbor)
			if !slices.Equal(edges, wantEdges) {
				t.Errorf("AllNeighbors(%s) yielded unexpected edges for %s; want: %v; got: %v", nodeID, neighbor, wantEdges, edges)
			}
		}
		if !want.Equal(got) {
			t.Errorf("AllNeighbors(%s) yielded unexpected neighbors; want: %v; got: %v", nodeID, want, got)
		}
	}

	yields := 0
	for range g.AllNeighbors("port") {
		yields++
		break
	}
	if yields != 1 {
		t.Errorf("AllNeighbors with early break yielded %d times; want: 1", yields)
	}
}

func TestAllNodesOfKind(t *testing.T) {
	g := New()
	mustUpsertEntities(t, g, testGraph.entities)

	got := set.NewSet[string]()
	for node := range g.AllNodesOfKind("EK_NETWORK_NODE") {
		got.Add(node.GetID())
	}
	want := set.NewSet("node", "orphan")
	if !want.Equal(got) {
		t.Errorf("AllNodesOfKind(EK_NETWORK_NODE) yielded unexpected IDs; want: %v; got: %v", want, got)
	}

	for node := range g.AllNodesOfKind("EK_NO_SUCH_KIND") {
		t.Errorf("AllNodesOfKind(EK_NO_SUCH_KIND) unexpectedly yielded %v", node)
	}

	yields := 0
	for range g.AllNodesOfKind("EK_NETWORK_NODE") {
		yields++
		break
	}
	if yields != 1 {
		t.Errorf("AllNodesOfKind with early break yielded %d times; want: 1", yields)
	}
}

func edgeKey(e *Edge) string {
	return e.GetA() + "|" + e.GetKind().String() + "|" + e.GetZ()
}

func TestAllEdges(t *testing.T) {
	g := New()
	mustUpsertEntities(t, g, testGraph.entities)
	mustAddRelationships(t, g, testGraph.relationships)

	want := set.NewSet[string]()
	for _, r := range testGraph.relationships {
		relationship := mustUnmarshalRelationship(t, r)
		want.Add(edgeKey(&Edge{relationship: relationship}))
	}

	got := set.NewSet[string]()
	yields := 0
	for e := range g.AllEdges() {
		got.Add(edgeKey(e))
		yields++
	}
	if !want.Equal(got) {
		t.Errorf("AllEdges yielded unexpected edges; want: %v; got: %v", want, got)
	}
	if yields != want.Cardinality() {
		t.Errorf("AllEdges yielded %d times for %d distinct edges", yields, want.Cardinality())
	}
}

func TestSelfLoopStoredOnce(t *testing.T) {
	g := New()
	mustUpsertEntities(t, g, []string{`id: "node" ek_network_node{}`})
	mustAddRelationships(t, g, []string{`a: "node" kind: RK_CONTAINS z: "node"`})

	if gotEdges := g.Edges("node", "node"); len(gotEdges) != 1 {
		t.Errorf("Edges(node, node) returned %d edges for a single self-loop; want: 1", len(gotEdges))
	}
	if _, err := g.AddRelationship(mustUnmarshalRelationship(t, `a: "node" kind: RK_CONTAINS z: "node"`)); err == nil {
		t.Errorf("AddRelationship of a duplicate self-loop unexpectedly succeeded")
	}
	if neighbors := g.Neighbors("node"); !slices.Equal(neighbors, []string{"node"}) {
		t.Errorf("Neighbors(node) returned %v for a self-loop; want: [node]", neighbors)
	}
}

func TestAllEdgesYieldsSelfLoopOnce(t *testing.T) {
	g := New()
	mustUpsertEntities(t, g, []string{`id: "node" ek_network_node{}`})
	mustAddRelationships(t, g, []string{`a: "node" kind: RK_CONTAINS z: "node"`})

	yields := 0
	for range g.AllEdges() {
		yields++
	}
	if yields != 1 {
		t.Errorf("AllEdges yielded a self-loop edge %d times; want: 1", yields)
	}
}

func benchmarkDuplicateGraph(b *testing.B) (*Graph, *npb.Relationship) {
	b.Helper()
	g := New()
	entity := &npb.Entity{}
	if err := prototext.Unmarshal([]byte(`id: "node" ek_network_node{}`), entity); err != nil {
		b.Fatalf("unable to unmarshal entity: %v", err)
	}
	if _, err := g.UpsertEntity(entity); err != nil {
		b.Fatalf("unable to upsert entity: %v", err)
	}
	r := &npb.Relationship{}
	if err := prototext.Unmarshal([]byte(`a: "node" kind: RK_CONTAINS z: "interface"`), r); err != nil {
		b.Fatalf("unable to unmarshal relationship: %v", err)
	}
	if _, err := g.AddRelationship(r); err != nil {
		b.Fatalf("unable to add relationship: %v", err)
	}
	return g, r
}

func BenchmarkAddRelationshipDuplicate(b *testing.B) {
	g, r := benchmarkDuplicateGraph(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AddRelationship(r)
	}
}

func BenchmarkTryAddRelationshipDuplicate(b *testing.B) {
	g, r := benchmarkDuplicateGraph(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.TryAddRelationship(r)
	}
}

func TestAllEdgesEarlyBreak(t *testing.T) {
	g := New()
	mustUpsertEntities(t, g, testGraph.entities)
	mustAddRelationships(t, g, testGraph.relationships)

	yields := 0
	for range g.AllEdges() {
		yields++
		break
	}
	if yields != 1 {
		t.Errorf("AllEdges with early break yielded %d times; want: 1", yields)
	}
}
