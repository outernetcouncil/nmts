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
	"testing"

	set "github.com/deckarep/golang-set/v2"

	npb "outernetcouncil.org/nmts/v2alpha/proto"
)

type walkTestCase struct {
	desc string
	graphEntities
	traverseFuncs []TraverseFunc
	// Maps starting node ID to the expected visited nodes.
	wantVisits map[string]set.Set[string]
}

func (tc *walkTestCase) Run(t *testing.T) {
	g := New()
	mustUpsertEntities(t, g, tc.entities)
	mustAddRelationships(t, g, tc.relationships)

	for from, wantVisits := range tc.wantVisits {
		gotVisits := set.NewSet[string]()
		dfs := DepthFirst{
			Visit:    func(_ *Graph, n string) { gotVisits.Add(n) },
			Traverse: TraverseAll(tc.traverseFuncs...),
		}
		dfs.Walk(g, from, nil)

		if !wantVisits.Equal(gotVisits) {
			t.Errorf("unexpected visits from %s; want: %v; got: %v", from, wantVisits, gotVisits)
		}
	}
}

var walkTestCases = []walkTestCase{
	{
		desc:          "traverse single edge type in either direction",
		graphEntities: testGraph,
		traverseFuncs: []TraverseFunc{
			Edges("EK_NETWORK_NODE", npb.RK_RK_CONTAINS, "EK_SDN_AGENT"),
		},
		wantVisits: map[string]set.Set[string]{
			"agent": set.NewSet("agent", "node"),
			"node":  set.NewSet("node", "agent"),
		},
	},
	{
		desc:          "traverse single edge type in same direction as relationship",
		graphEntities: testGraph,
		traverseFuncs: []TraverseFunc{
			EdgesFrom("EK_NETWORK_NODE", "EK_NETWORK_NODE", npb.RK_RK_CONTAINS, "EK_SDN_AGENT"),
		},
		wantVisits: map[string]set.Set[string]{
			"agent": set.NewSet("agent"),
			"node":  set.NewSet("node", "agent"),
		},
	},
	{
		desc:          "traverse single edge type in opposite direction as relationship",
		graphEntities: testGraph,
		traverseFuncs: []TraverseFunc{
			EdgesFrom("EK_SDN_AGENT", "EK_NETWORK_NODE", npb.RK_RK_CONTAINS, "EK_SDN_AGENT"),
		},
		wantVisits: map[string]set.Set[string]{
			"agent": set.NewSet("agent", "node"),
			"node":  set.NewSet("node"),
		},
	},
	{
		desc:          "traverse many edge types",
		graphEntities: testGraph,
		traverseFuncs: []TraverseFunc{
			EdgesFrom("EK_SDN_AGENT", "EK_NETWORK_NODE", npb.RK_RK_CONTAINS, "EK_SDN_AGENT"),
			Edges("EK_NETWORK_NODE", npb.RK_RK_CONTAINS, "EK_INTERFACE"),
			Edges("EK_INTERFACE", npb.RK_RK_TRAVERSES, "EK_PORT"),
			Edges("EK_INTERFACE", npb.RK_RK_TRAVERSES, "EK_PORT"),
			EdgesFrom("EK_PORT", "EK_PORT", npb.RK_RK_ORIGINATES, "EK_MODULATOR"),
			EdgesFrom("EK_DEMODULATOR", "EK_PORT", npb.RK_RK_TERMINATES, "EK_DEMODULATOR"),
		},
		wantVisits: map[string]set.Set[string]{
			"agent":       set.NewSet("agent", "node", "interface", "port", "modulator"),
			"node":        set.NewSet("node", "interface", "port", "modulator"),
			"interface":   set.NewSet("node", "interface", "port", "modulator"),
			"port":        set.NewSet("node", "interface", "port", "modulator"),
			"modulator":   set.NewSet("modulator"),
			"demodulator": set.NewSet("node", "interface", "port", "modulator", "demodulator"),
		},
	},
}

func TestDepthFirstWalk(t *testing.T) {
	for _, tc := range walkTestCases {
		t.Run(tc.desc, tc.Run)
	}
}
