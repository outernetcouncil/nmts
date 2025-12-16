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

package utilities_test

import (
	"slices"
	"strings"
	"testing"

	set "github.com/deckarep/golang-set/v2"
	gcmp "github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/samber/lo"
	"google.golang.org/protobuf/testing/protocmp"
	"outernetcouncil.org/nmts/v1/lib/entityrelationship"
	"outernetcouncil.org/nmts/v1/lib/graph"
	graphutil "outernetcouncil.org/nmts/v1/lib/utilities"
	testutil "outernetcouncil.org/nmts/v1/lib/utilities/testing"
	nmtspb "outernetcouncil.org/nmts/v1/proto"
	logicalpb "outernetcouncil.org/nmts/v1/proto/ek/logical"
	ietfpb "outernetcouncil.org/nmts/v1/proto/types/ietf"
)

const (
	interfaceID      = "uuid(gs1/network_node/interfaces/antenna0)"
	bogusInterfaceID = "bogusInterface"
)

func getBasicWorkingNMTSFragment() *nmtspb.Fragment {
	const txtpb = `
##
# Ground station "gs1"
##
entity {
	id: "uuid(gs1/platform)"
	ek_platform {
		name: "gs1"
		category_tag: "Gateway"
		motion {
			entry {
				# Onizuka AFB, a.k.a. the "Blue Cube"
				geodetic_msl {
					latitude_deg: 37.4058958222506
					longitude_deg: -122.02679892112421
					height_msl_m: 30
				}
			}
		}
	}
}
entity {
	id: "uuid(gs1/platform/ports/unused_port0)"
	ek_port {
		name: "unused_port0"
	}
}
entity {
	id: "uuid(gs1/platform/ports/antenna0)"
	ek_port {
		name: "antenna0"
	}
}
relationship {
	a: "uuid(gs1/platform)"
	kind: RK_CONTAINS
	z: "uuid(gs1/platform/ports/antenna0)"
}
entity {
	id: "uuid(gs1/platform/modulators/0)"
	ek_modulator {}
}
relationship {
	a: "uuid(gs1/platform)"
	kind: RK_CONTAINS
	z: "uuid(gs1/platform/modulators/0)"
}
relationship {
	a: "uuid(gs1/platform/ports/antenna0)"
	kind: RK_ORIGINATES
	z: "uuid(gs1/platform/modulators/0)"
}
entity {
	id: "uuid(gs1/platform/sigproc/tx0)"
	ek_signal_processing_chain {}
}
relationship {
	a: "uuid(gs1/platform)"
	kind: RK_CONTAINS
	z: "uuid(gs1/platform/sigproc/tx0)"
}
relationship {
	a: "uuid(gs1/platform/modulators/0)"
	kind: RK_SIGNAL_TRANSITS
	z: "uuid(gs1/platform/sigproc/tx0)"
}
entity {
	id: "uuid(gs1/platform/transmitters/0)"
	ek_transmitter {}
}
relationship {
	a: "uuid(gs1/platform)"
	kind: RK_CONTAINS
	z: "uuid(gs1/platform/transmitters/0)"
}
relationship {
	a: "uuid(gs1/platform/sigproc/tx0)"
	kind: RK_SIGNAL_TRANSITS
	z: "uuid(gs1/platform/transmitters/0)"
}
entity {
	id: "uuid(gs1/platform/antennas/0)"
	ek_antenna {}
}
relationship {
	a: "uuid(gs1/platform)"
	kind: RK_CONTAINS
	z: "uuid(gs1/platform/antennas/0)"
}
relationship {
	a: "uuid(gs1/platform/transmitters/0)"
	kind: RK_SIGNAL_TRANSITS
	z: "uuid(gs1/platform/antennas/0)"
}
entity {
	id: "uuid(gs1/platform/receivers/0)"
	ek_receiver {}
}
relationship {
	a: "uuid(gs1/platform)"
	kind: RK_CONTAINS
	z: "uuid(gs1/platform/receivers/0)"
}
relationship {
	a: "uuid(gs1/platform/antennas/0)"
	kind: RK_SIGNAL_TRANSITS
	z: "uuid(gs1/platform/receivers/0)"
}
entity {
	id: "uuid(gs1/platform/sigproc/rx0)"
	ek_signal_processing_chain {}
}
relationship {
	a: "uuid(gs1/platform)"
	kind: RK_CONTAINS
	z: "uuid(gs1/platform/sigproc/rx0)"
}
relationship {
	a: "uuid(gs1/platform/receivers/0)"
	kind: RK_SIGNAL_TRANSITS
	z: "uuid(gs1/platform/sigproc/rx0)"
}
entity {
	id: "uuid(gs1/platform/demodulators/0)"
	ek_demodulator {}
}
relationship {
	a: "uuid(gs1/platform)"
	kind: RK_CONTAINS
	z: "uuid(gs1/platform/demodulators/0)"
}
relationship {
	a: "uuid(gs1/platform/sigproc/rx0)"
	kind: RK_SIGNAL_TRANSITS
	z: "uuid(gs1/platform/demodulators/0)"
}
relationship {
	a: "uuid(gs1/platform/ports/antenna0)"
	kind: RK_TERMINATES
	z: "uuid(gs1/platform/demodulators/0)"
}
entity {
	id: "uuid(gs1/network_node)"
	ek_network_node {
		name: "gs1"
		category_tag: "Gateway"
	}
}
relationship {
	a: "uuid(gs1/platform)"
	kind: RK_CONTAINS
	z: "uuid(gs1/network_node)"
}
entity {
	id: "uuid(gs1/network_node/interfaces/antenna0)"
	ek_interface {
		name: "antenna0"
		gse {}
	}
}
relationship {
	a: "uuid(gs1/network_node)"
	kind: RK_CONTAINS
	z: "uuid(gs1/network_node/interfaces/antenna0)"
}
relationship {
	a: "uuid(gs1/network_node/interfaces/antenna0)"
	kind: RK_TRAVERSES
	z: "uuid(gs1/platform/ports/antenna0)"
}
relationship {
	a: "uuid(gs1/network_node/interfaces/antenna0)"
	kind: RK_TRAVERSES
	z: "uuid(gs1/platform/ports/unused_port0)"
}
entity {
	id: "uuid(gs1/network_node/interfaces/antenna0.mpls)"
	ek_interface {
		name: "antenna0"
		admin_status: IF_ADMIN_STATUS_UP
		mpls {
		}
	}
}
relationship {
	a: "uuid(gs1/network_node)"
	kind: RK_CONTAINS
	z: "uuid(gs1/network_node/interfaces/antenna0.mpls)"
}
relationship {
	a: "uuid(gs1/network_node/interfaces/antenna0.mpls)"
	kind: RK_TRAVERSES
	z: "uuid(gs1/network_node/interfaces/antenna0)"
}
entity {
	id: "uuid(gs1/network_node/logical_packet_links/antenna0.mpls/0)"
	ek_logical_packet_link {
		sr {
			enabled: true
			adjacency_sid: { mpls: 3001 }
		}
		max_data_rate_bps: 10737418240
	}
}
relationship {
	a: "uuid(gs1/network_node/interfaces/antenna0.mpls)"
	kind: RK_ORIGINATES
	z: "uuid(gs1/network_node/logical_packet_links/antenna0.mpls/0)"
}
entity {
	id: "uuid(gs1/network_node/interfaces/eth0)"
	ek_interface {
		name: "eth0"
		eth {
			mac_addr {
				str: "00:01:02:03:04:05"
			}
		}
	}
}
relationship {
	a: "uuid(gs1/network_node)"
	kind: RK_CONTAINS
	z: "uuid(gs1/network_node/interfaces/eth0)"
}
entity {
	id: "uuid(gs1/network_node/interfaces/eth0.mpls)"
	ek_interface {
		name: "eth0"
		admin_status: IF_ADMIN_STATUS_UP
		mpls {
		}
	}
}
relationship {
	a: "uuid(gs1/network_node)"
	kind: RK_CONTAINS
	z: "uuid(gs1/network_node/interfaces/eth0.mpls)"
}
relationship {
	a: "uuid(gs1/network_node/interfaces/eth0.mpls)"
	kind: RK_TRAVERSES
	z: "uuid(gs1/network_node/interfaces/eth0)"
}
entity {
	id: "uuid(gs1/network_node/logical_packet_links/eth0/0)"
	ek_logical_packet_link {
		max_data_rate_bps: 10737418240
	}
}
relationship {
	a: "uuid(gs1/network_node/interfaces/eth0)"
	kind: RK_ORIGINATES
	z: "uuid(gs1/network_node/logical_packet_links/eth0/0)"
}
entity {
	id: "uuid(gs1/network_node/logical_packet_links/eth0.mpls/0)"
	ek_logical_packet_link {
		sr {
			enabled: true
			adjacency_sid: { mpls: 1001 }
		}
		max_data_rate_bps: 10737418240
	}
}
relationship {
	a: "uuid(gs1/network_node/interfaces/eth0.mpls)"
	kind: RK_ORIGINATES
	z: "uuid(gs1/network_node/logical_packet_links/eth0.mpls/0)"
}
relationship {
	a: "uuid(gs1/network_node/logical_packet_links/eth0.mpls/0)"
	kind: RK_TRAVERSES
	z: "uuid(gs1/network_node/logical_packet_links/eth0/0)"
}
entity {
	id: "uuid(gs1/network_node/ipvrf0)"
	ek_route_fn {
		router_id: {
			dotted_quad: {
				str: "1.1.1.1"
			}
		}
		sr {
			enabled: true
			node_sid: { mpls: 1111 }
		}
	}
}
relationship {
	kind: RK_CONTAINS
	a: "uuid(gs1/network_node)"
	z: "uuid(gs1/network_node/ipvrf0)"
}
relationship {
	kind: RK_CONTROLS
	a: "uuid(gs1/network_node/ipvrf0)"
	z: "uuid(gs1/network_node)"
}
entity {
	id: "uuid(gs1/network_node/bpagent)"
	ek_bp_agent_fn {
		max_capacity_bytes: 10737418240
	}
}
relationship {
	kind: RK_CONTAINS
	a: "uuid(gs1/network_node)"
	z: "uuid(gs1/network_node/bpagent)"
}
##
# PoP "pop1"
#
entity {
	id: "uuid(pop1/network_node)"
	ek_network_node {
		name: "pop1"
		category_tag: "PoP"
	}
}
entity {
	id: "uuid(pop1/network_node/interfaces/eth0)"
	ek_interface {
		name: "eth0"
		eth {
			mac_addr {
				str: "00:02:03:04:05:06"
			}
		}
	}
}
relationship {
	a: "uuid(pop1/network_node)"
	kind: RK_CONTAINS
	z: "uuid(pop1/network_node/interfaces/eth0)"
}
entity {
	id: "uuid(pop1/network_node/interfaces/eth0.mpls)"
	ek_interface {
		name: "eth0"
		admin_status: IF_ADMIN_STATUS_UP
		mpls {
		}
	}
}
relationship {
	a: "uuid(pop1/network_node)"
	kind: RK_CONTAINS
	z: "uuid(pop1/network_node/interfaces/eth0.mpls)"
}
relationship {
	a: "uuid(pop1/network_node/interfaces/eth0.mpls)"
	kind: RK_TRAVERSES
	z: "uuid(pop1/network_node/interfaces/eth0)"
}
entity {
	id: "uuid(pop1/network_node/logical_packet_links/eth0.mpls/0)"
	ek_logical_packet_link {
		sr {
			enabled: true
			adjacency_sid: { mpls: 5001 }
		}
		max_data_rate_bps: 10737418240
	}
}
relationship {
	a: "uuid(pop1/network_node/interfaces/eth0.mpls)"
	kind: RK_ORIGINATES
	z: "uuid(pop1/network_node/logical_packet_links/eth0.mpls/0)"
}
##
# Links between pop1 and gs1
##
relationship {
	a: "uuid(pop1/network_node/interfaces/eth0.mpls)"
	kind: RK_TERMINATES
	z: "uuid(gs1/network_node/logical_packet_links/eth0.mpls/0)"
}
relationship {
	a: "uuid(gs1/network_node/interfaces/eth0.mpls)"
	kind: RK_TERMINATES
	z: "uuid(pop1/network_node/logical_packet_links/eth0.mpls/0)"
}
##
# An SDN Agent responsible for most elements in the model.
##
entity {
	id: "uuid(sdn_agents/0)"
	ek_sdn_agent{}
}
relationship {
	a: "uuid(sdn_agents/0)"
	kind: RK_CONTROLS
	z: "uuid(gs1/network_node/bpagent)"
}
relationship {
	a: "uuid(sdn_agents/0)"
	kind: RK_CONTROLS
	z: "uuid(gs1/network_node/ipvrf0)"
}
relationship {
	a: "uuid(sdn_agents/0)"
	kind: RK_CONTROLS
	z: "uuid(pop1/network_node)"
}
relationship {
	a: "uuid(sdn_agents/0)"
	kind: RK_CONTROLS
	z: "uuid(gs1/platform/modulators/0)"
}
`

	return lo.Must(testutil.FragmentFrom(txtpb))
}

func getBasicBentPipeSatelliteNMTSFragment() []*nmtspb.Fragment {
	fragments := []*nmtspb.Fragment{}

	const sat1Txtpb = `
##
# Satellite "sat1"
##
entity { id: "uuid(sat1/platform)"              ek_platform {} }
entity { id: "uuid(sat1/platform/rx_antenna)"   ek_antenna {} }
entity { id: "uuid(sat1/platform/tx_antenna)"   ek_antenna {} }
relationship { a: "uuid(sat1/platform)" kind: RK_CONTAINS z: "uuid(sat1/platform/rx_antenna)" }
relationship { a: "uuid(sat1/platform)" kind: RK_CONTAINS z: "uuid(sat1/platform/tx_antenna)" }
`
	fragments = append(fragments, lo.Must(testutil.FragmentFrom(sat1Txtpb)))

	const transponderTemplateTxtpb = `
entity { id: "uuid(sat1/platform/receiver{{idx}})"    ek_receiver {} }
entity { id: "uuid(sat1/platform/rxrfchain{{idx}})"   ek_signal_processing_chain {} }
entity { id: "uuid(sat1/platform/freqmixer{{idx}})"   ek_signal_processing_chain {} }
entity { id: "uuid(sat1/platform/txrfchain{{idx}})"   ek_signal_processing_chain {} }
entity { id: "uuid(sat1/platform/transmitter{{idx}})" ek_transmitter {} }
relationship { a: "uuid(sat1/platform/rx_antenna)"         kind: RK_SIGNAL_TRANSITS z: "uuid(sat1/platform/receiver{{idx}})" }
relationship { a: "uuid(sat1/platform/receiver{{idx}})"    kind: RK_SIGNAL_TRANSITS z: "uuid(sat1/platform/rxrfchain{{idx}})" }
relationship { a: "uuid(sat1/platform/rxrfchain{{idx}})"   kind: RK_SIGNAL_TRANSITS z: "uuid(sat1/platform/freqmixer{{idx}})" }
relationship { a: "uuid(sat1/platform/freqmixer{{idx}})"   kind: RK_SIGNAL_TRANSITS z: "uuid(sat1/platform/txrfchain{{idx}})" }
relationship { a: "uuid(sat1/platform/txrfchain{{idx}})"   kind: RK_SIGNAL_TRANSITS z: "uuid(sat1/platform/transmitter{{idx}})" }
relationship { a: "uuid(sat1/platform/transmitter{{idx}})" kind: RK_SIGNAL_TRANSITS z: "uuid(sat1/platform/tx_antenna)" }
`
	for _, idx := range []string{"0", "1", "2"} {
		fragments = append(fragments, lo.Must(testutil.FragmentFrom(strings.ReplaceAll(transponderTemplateTxtpb, "{{idx}}", idx))))
	}

	return fragments
}

func getBentPipeTransponderNMTSFragments() []*nmtspb.Fragment {
	fragments := []*nmtspb.Fragment{}

	const satelliteTxtpb = `
entity { id: "satellite/platform"     ek_platform{} }
entity { id: "satellite/network_node" ek_network_node{} }
entity { id: "satellite/rx_antenna"   ek_antenna{} }
entity { id: "satellite/receiver"     ek_receiver{} }
entity { id: "satellite/transmitter"  ek_transmitter{} }
entity { id: "satellite/tx_antenna"   ek_antenna{} }
entity { id: "satellite/sigproc"      ek_signal_processing_chain {} }
relationship{ a: "satellite/platform"    kind: RK_CONTAINS        z: "satellite/network_node" }
relationship{ a: "satellite/platform"    kind: RK_CONTAINS        z: "satellite/rx_antenna" }
relationship{ a: "satellite/platform"    kind: RK_CONTAINS        z: "satellite/receiver" }
relationship{ a: "satellite/platform"    kind: RK_CONTAINS        z: "satellite/transmitter" }
relationship{ a: "satellite/platform"    kind: RK_CONTAINS        z: "satellite/tx_antenna" }
relationship{ a: "satellite/rx_antenna"  kind: RK_SIGNAL_TRANSITS z: "satellite/receiver" }
relationship{ a: "satellite/receiver"    kind: RK_SIGNAL_TRANSITS z: "satellite/sigproc" }
relationship{ a: "satellite/sigproc"     kind: RK_SIGNAL_TRANSITS z: "satellite/transmitter" }
relationship{ a: "satellite/transmitter" kind: RK_SIGNAL_TRANSITS z: "satellite/tx_antenna" }
`
	fragments = append(fragments, lo.Must(testutil.FragmentFrom(satelliteTxtpb)))

	const terminalTxtpb = `
entity { id: "{{ID}}/platform"     ek_platform{} }
entity { id: "{{ID}}/network_node" ek_network_node{} }
entity { id: "{{ID}}/port"         ek_port{} }
entity { id: "{{ID}}/interface"    ek_interface{} }
entity { id: "{{ID}}/modulator"    ek_modulator{} }
entity { id: "{{ID}}/transmitter"  ek_transmitter{} }
entity { id: "{{ID}}/antenna"      ek_antenna{} }
entity { id: "{{ID}}/tx_sigproc"   ek_signal_processing_chain{} }
entity { id: "{{ID}}/demodulator"  ek_demodulator{} }
entity { id: "{{ID}}/receiver"     ek_receiver{} }
entity { id: "{{ID}}/rx_sigproc"   ek_signal_processing_chain{} }
relationship{ a: "{{ID}}/platform"     kind: RK_CONTAINS        z: "{{ID}}/network_node" }
relationship{ a: "{{ID}}/platform"     kind: RK_CONTAINS        z: "{{ID}}/port" }
relationship{ a: "{{ID}}/network_node" kind: RK_CONTAINS        z: "{{ID}}/interface" }
relationship{ a: "{{ID}}/interface"    kind: RK_TRAVERSES       z: "{{ID}}/port" }
relationship{ a: "{{ID}}/platform"     kind: RK_CONTAINS        z: "{{ID}}/modulator" }
relationship{ a: "{{ID}}/platform"     kind: RK_CONTAINS        z: "{{ID}}/transmitter" }
relationship{ a: "{{ID}}/platform"     kind: RK_CONTAINS        z: "{{ID}}/antenna" }
relationship{ a: "{{ID}}/platform"     kind: RK_CONTAINS        z: "{{ID}}/receiver" }
relationship{ a: "{{ID}}/platform"     kind: RK_CONTAINS        z: "{{ID}}/demodulator" }
relationship{ a: "{{ID}}/port"         kind: RK_ORIGINATES      z: "{{ID}}/modulator" }
relationship{ a: "{{ID}}/modulator"    kind: RK_SIGNAL_TRANSITS z: "{{ID}}/tx_sigproc" }
relationship{ a: "{{ID}}/tx_sigproc"   kind: RK_SIGNAL_TRANSITS z: "{{ID}}/transmitter" }
relationship{ a: "{{ID}}/transmitter"  kind: RK_SIGNAL_TRANSITS z: "{{ID}}/antenna" }
relationship{ a: "{{ID}}/antenna"      kind: RK_SIGNAL_TRANSITS z: "{{ID}}/receiver" }
relationship{ a: "{{ID}}/receiver"     kind: RK_SIGNAL_TRANSITS z: "{{ID}}/rx_sigproc" }
relationship{ a: "{{ID}}/rx_sigproc"   kind: RK_SIGNAL_TRANSITS z: "{{ID}}/demodulator" }
relationship{ a: "{{ID}}/port"         kind: RK_TERMINATES      z: "{{ID}}/demodulator" }
`
	for _, terminalID := range []string{"tx_terminal", "rx_terminal"} {
		fragments = append(fragments, lo.Must(testutil.FragmentFrom(strings.ReplaceAll(terminalTxtpb, "{{ID}}", terminalID))))
	}

	const physicalMediumUplinkTxtpb = `
entity { id: "physical-uplink"   ek_physical_medium_link{} }
relationship { a: "tx_terminal/antenna"  kind: RK_ORIGINATES z: "physical-uplink" }
relationship { a: "satellite/rx_antenna" kind: RK_TERMINATES z: "physical-uplink" }
`
	fragments = append(fragments, lo.Must(testutil.FragmentFrom(physicalMediumUplinkTxtpb)))

	const physicalMediumDownlinkTxtpb = `
entity { id: "physical-downlink" ek_physical_medium_link{} }
relationship { a: "satellite/tx_antenna" kind: RK_ORIGINATES z: "physical-downlink" }
relationship { a: "rx_terminal/antenna"  kind: RK_TERMINATES z: "physical-downlink" }
`
	fragments = append(fragments, lo.Must(testutil.FragmentFrom(physicalMediumDownlinkTxtpb)))

	const logicalPacketUplinkTxtpb = `
entity { id: "logical-uplink"   ek_logical_packet_link{} }
relationship { a: "tx_terminal/interface" kind: RK_ORIGINATES z: "logical-uplink" }
relationship { a: "rx_terminal/interface" kind: RK_TERMINATES z: "logical-uplink" }
relationship { a: "logical-uplink"        kind: RK_TRAVERSES  z: "physical-uplink" }
relationship { a: "logical-uplink"        kind: RK_TRAVERSES  z: "physical-downlink" }
`
	fragments = append(fragments, lo.Must(testutil.FragmentFrom(logicalPacketUplinkTxtpb)))

	return fragments
}

func getBasicWorkingGraph() *graph.Graph {
	return lo.Must(testutil.GraphFromFragments(getBasicWorkingNMTSFragment()))
}

func getBasicBentPipeSatellite() *graph.Graph {
	return lo.Must(testutil.GraphFromFragments(getBasicBentPipeSatelliteNMTSFragment()...))
}

// Testing sets of pointers to graph Edges for equality is a bit of
// a pain. Convert to sets of strings for comparison instead.
func toStringSet(relationships []*nmtspb.Relationship) set.Set[string] {
	strs := set.NewSet[string]()
	for _, relationship := range relationships {
		rel := entityrelationship.RelationshipFromProto(relationship)
		strs.Add(rel.String())
	}
	return strs
}

func TestEdgesIncidentToAntennaPort(t *testing.T) {
	g := getBasicWorkingGraph()
	want := []*nmtspb.Relationship{
		{
			A:    "uuid(gs1/platform)",
			Kind: nmtspb.RK_RK_CONTAINS,
			Z:    "uuid(gs1/platform/ports/antenna0)",
		},
		{
			A:    "uuid(gs1/platform/ports/antenna0)",
			Kind: nmtspb.RK_RK_ORIGINATES,
			Z:    "uuid(gs1/platform/modulators/0)",
		},
		{
			A:    "uuid(gs1/platform/ports/antenna0)",
			Kind: nmtspb.RK_RK_TERMINATES,
			Z:    "uuid(gs1/platform/demodulators/0)",
		},
		{
			A:    "uuid(gs1/network_node/interfaces/antenna0)",
			Kind: nmtspb.RK_RK_TRAVERSES,
			Z:    "uuid(gs1/platform/ports/antenna0)",
		},
	}
	got := graphutil.EdgesIncidentTo(g, "uuid(gs1/platform/ports/antenna0)")
	gotRelationships := lo.Map(got.ToSlice(), func(edge *graph.Edge, _ int) *nmtspb.Relationship {
		return edge.GetRelationship()
	})
	diff := toStringSet(want).Difference(toStringSet(gotRelationships))
	if !diff.IsEmpty() {
		t.Errorf("want: %q, got: %q, diff: %q", want, got, diff)
	}
}

func TestOutEdgesFromAntennaPort(t *testing.T) {
	g := getBasicWorkingGraph()
	want := []*nmtspb.Relationship{
		{
			A:    "uuid(gs1/platform/ports/antenna0)",
			Kind: nmtspb.RK_RK_ORIGINATES,
			Z:    "uuid(gs1/platform/modulators/0)",
		},
		{
			A:    "uuid(gs1/platform/ports/antenna0)",
			Kind: nmtspb.RK_RK_TERMINATES,
			Z:    "uuid(gs1/platform/demodulators/0)",
		},
	}
	got := graphutil.OutEdgesFrom(g, "uuid(gs1/platform/ports/antenna0)")
	gotRelationships := lo.Map(got.ToSlice(), func(edge *graph.Edge, _ int) *nmtspb.Relationship {
		return edge.GetRelationship()
	})
	diff := toStringSet(want).Difference(toStringSet(gotRelationships))
	if !diff.IsEmpty() {
		t.Errorf("want: %q, got: %q, diff: %q", want, got, diff)
	}
}

func TestInEdgeToAntennaPort(t *testing.T) {
	g := getBasicWorkingGraph()
	want := []*nmtspb.Relationship{
		{
			A:    "uuid(gs1/platform)",
			Kind: nmtspb.RK_RK_CONTAINS,
			Z:    "uuid(gs1/platform/ports/antenna0)",
		},
		{
			A:    "uuid(gs1/network_node/interfaces/antenna0)",
			Kind: nmtspb.RK_RK_TRAVERSES,
			Z:    "uuid(gs1/platform/ports/antenna0)",
		},
	}
	got := graphutil.InEdgesTo(g, "uuid(gs1/platform/ports/antenna0)")
	gotRelationships := lo.Map(got.ToSlice(), func(edge *graph.Edge, _ int) *nmtspb.Relationship {
		return edge.GetRelationship()
	})
	diff := toStringSet(want).Difference(toStringSet(gotRelationships))
	if !diff.IsEmpty() {
		t.Errorf("want: %q, got: %q, diff: %q", want, got, diff)
	}
}

func TestEdgeSetFilterForSomeAntennaPortRelationships(t *testing.T) {
	g := getBasicWorkingGraph()
	want := []*nmtspb.Relationship{
		{
			A:    "uuid(gs1/platform/ports/antenna0)",
			Kind: nmtspb.RK_RK_ORIGINATES,
			Z:    "uuid(gs1/platform/modulators/0)",
		},
		{
			A:    "uuid(gs1/network_node/interfaces/antenna0)",
			Kind: nmtspb.RK_RK_TRAVERSES,
			Z:    "uuid(gs1/platform/ports/antenna0)",
		},
	}
	got := graphutil.FilterFor(graphutil.EdgesIncidentTo(g, "uuid(gs1/platform/ports/antenna0)"), nmtspb.RK_RK_ORIGINATES, nmtspb.RK_RK_TRAVERSES)
	gotRelationships := lo.Map(got.ToSlice(), func(edge *graph.Edge, _ int) *nmtspb.Relationship {
		return edge.GetRelationship()
	})
	diff := toStringSet(want).Difference(toStringSet(gotRelationships))
	if !diff.IsEmpty() {
		t.Errorf("want: %q, got: %q, diff: %q", want, got, diff)
	}
}

func TestFindEncompassingPlatformForPlatform(t *testing.T) {
	g := getBasicWorkingGraph()
	const platformID = "uuid(gs1/platform)"
	if p := graphutil.FindEncompassingPlatform(g, platformID); p != platformID {
		t.Errorf("want: %q, got: %q", platformID, p)
	}
}

func TestFindEncompassingPlatformForAntenna(t *testing.T) {
	g := getBasicWorkingGraph()
	const platformID = "uuid(gs1/platform)"
	if p := graphutil.FindEncompassingPlatform(g, "uuid(gs1/platform/antennas/0)"); p != platformID {
		t.Errorf("want: %q, got: %q", platformID, p)
	}

	// Remove the direct relationship between the platform and
	// the antenna port. A path should still be found via other
	// connected elements.
	g.RemoveRelationship(&nmtspb.Relationship{
		A:    platformID,
		Kind: nmtspb.RK_RK_CONTAINS,
		Z:    "uuid(gs1/platform/ports/antenna0)",
	})
	if p := graphutil.FindEncompassingPlatform(g, "uuid(gs1/platform/antennas/0)"); p != platformID {
		t.Errorf("want: %q, got: %q", platformID, p)
	}
}

func FindEncompassingNetworkNodeForNetworkNode(t *testing.T) {
	g := getBasicWorkingGraph()
	const networkNodeID = "uuid(gs1/network_node)"
	if n := graphutil.FindEncompassingPlatform(g, networkNodeID); n != networkNodeID {
		t.Errorf("want: %q, got: %q", networkNodeID, n)
	}
}

func FindEncompassingNetworkNodeForAntenna(t *testing.T) {
	g := getBasicWorkingGraph()
	const networkNodeID = "uuid(gs1/network_node)"
	if n := graphutil.FindEncompassingPlatform(g, "uuid(gs1/platform/antennas/0)"); n != networkNodeID {
		t.Errorf("want: %q, got: %q", networkNodeID, n)
	}
}

func FindEncompassingNetworkNodeForDemodulator(t *testing.T) {
	g := getBasicWorkingGraph()
	const networkNodeID = "uuid(gs1/network_node)"
	if n := graphutil.FindEncompassingPlatform(g, "uuid(gs1/platform/demodulators/0)"); n != networkNodeID {
		t.Errorf("want: %q, got: %q", networkNodeID, n)
	}
}

func FindEncompassingNetworkNodeForRouteFn(t *testing.T) {
	g := getBasicWorkingGraph()
	const networkNodeID = "uuid(gs1/network_node)"
	if n := graphutil.FindEncompassingPlatform(g, "uuid(gs1/network_node/ipvrf0)"); n != networkNodeID {
		t.Errorf("want: %q, got: %q", networkNodeID, n)
	}
}

func TestFindAssociatedPortForAntenna(t *testing.T) {
	g := getBasicWorkingGraph()
	const portID = "uuid(gs1/platform/ports/antenna0)"
	if p := graphutil.FindAssociatedPort(g, "uuid(gs1/platform/antennas/0)"); p != portID {
		t.Errorf("want: %q, got: %q", portID, p)
	}
}

func TestFindBentPipeReceiverForTransmitter(t *testing.T) {
	type expectation struct {
		start    string
		receiver string
	}
	expectations := []expectation{
		{"uuid(sat1/platform/transmitter0)", "uuid(sat1/platform/receiver0)"},
		{"uuid(sat1/platform/transmitter1)", "uuid(sat1/platform/receiver1)"},
		{"uuid(sat1/platform/transmitter2)", "uuid(sat1/platform/receiver2)"},
	}

	g := getBasicBentPipeSatellite()

	for _, expect := range expectations {
		if p := graphutil.FindBentPipeReceiverFromTransmitter(g, expect.start); p != expect.receiver {
			t.Errorf("start: %q, want: %q, got: %q", expect.start, expect.receiver, p)
		}
	}
}

func TestFindBentPipeReceiverForTransmitterReturnsEmptyString(t *testing.T) {
	type expectation struct {
		start    string
		receiver string
	}
	expectations := []expectation{
		{"uuid(sat1/platform/transmitter0)", ""},
		{"uuid(sat1/platform/transmitter1)", "uuid(sat1/platform/receiver1)"},
		{"uuid(sat1/platform/transmitter2)", "uuid(sat1/platform/receiver2)"},
	}

	g := getBasicBentPipeSatellite()
	g.RemoveEntity("uuid(sat1/platform/freqmixer0)")
	g.RemoveRelationship(&nmtspb.Relationship{
		A:    "uuid(sat1/platform/freqmixer0)",
		Kind: nmtspb.RK_RK_SIGNAL_TRANSITS,
		Z:    "uuid(sat1/platform/txrfchain0)",
	})
	g.RemoveRelationship(&nmtspb.Relationship{
		A:    "uuid(sat1/platform/rxrfchain0)",
		Kind: nmtspb.RK_RK_SIGNAL_TRANSITS,
		Z:    "uuid(sat1/platform/freqmixer0)",
	})

	for _, expect := range expectations {
		if p := graphutil.FindBentPipeReceiverFromTransmitter(g, expect.start); p != expect.receiver {
			t.Errorf("start: %q, want: %q, got: %q", expect.start, expect.receiver, p)
		}
	}
}

func TestFindEncompassingEntitiesForBentPipeTransponderTopology(t *testing.T) {
	type EncompassingExpectation struct {
		start         string
		ekPlatform    string
		ekNetworkNode string
	}

	expectations := []EncompassingExpectation{
		{start: "satellite/receiver", ekPlatform: "satellite/platform", ekNetworkNode: "satellite/network_node"},
		{start: "satellite/transmitter", ekPlatform: "satellite/platform", ekNetworkNode: "satellite/network_node"},
		{start: "tx_terminal/antenna", ekPlatform: "tx_terminal/platform", ekNetworkNode: "tx_terminal/network_node"},
		{start: "rx_terminal/antenna", ekPlatform: "rx_terminal/platform", ekNetworkNode: "rx_terminal/network_node"},
	}

	for _, expect := range expectations {
		g := lo.Must(testutil.GraphFromFragments(getBentPipeTransponderNMTSFragments()...))
		if n := graphutil.FindEncompassingPlatform(g, expect.start); n != expect.ekPlatform {
			t.Errorf("from: %q want: %q, got: %q", expect.start, expect.ekPlatform, n)
		}
		// removing the expected EK_PLATFORM should yield "".
		if err := g.RemoveEntity(expect.ekPlatform); err != nil {
			t.Errorf("failed to remove %q from graph", expect.ekPlatform)
		}
		if n := graphutil.FindEncompassingPlatform(g, expect.start); n != "" {
			t.Errorf("from: %q want: '', got: %q", expect.start, n)
		}

		g = lo.Must(testutil.GraphFromFragments(getBentPipeTransponderNMTSFragments()...))
		if n := graphutil.FindEncompassingNetworkNode(g, expect.start); n != expect.ekNetworkNode {
			t.Errorf("from: %q want: %q, got: %q", expect.start, expect.ekNetworkNode, n)
		}
		// removing the expected EK_NETWORK_NODE should yield "".
		if err := g.RemoveEntity(expect.ekNetworkNode); err != nil {
			t.Errorf("failed to remove %q from graph", expect.ekNetworkNode)
		}
		if n := graphutil.FindEncompassingNetworkNode(g, expect.start); n != "" {
			t.Errorf("from: %q want: '', got: %q", expect.start, n)
		}
	}
}

func TestFindRootInterfaceBeneathBaseInterface(t *testing.T) {
	g := getBasicWorkingGraph()
	const baseInterfaceID = "uuid(gs1/network_node/interfaces/antenna0)"
	if i := graphutil.FindRootInterfaceBeneath(g, baseInterfaceID); i != baseInterfaceID {
		t.Errorf("want: %q, got: %q", baseInterfaceID, i)
	}
}

func TestFindRootInterfaceBeneathUpperInterface(t *testing.T) {
	g := getBasicWorkingGraph()
	const baseInterfaceID = "uuid(gs1/network_node/interfaces/antenna0)"
	if i := graphutil.FindRootInterfaceBeneath(g, "uuid(gs1/network_node/interfaces/antenna0.mpls)"); i != baseInterfaceID {
		t.Errorf("want: %q, got: %q", baseInterfaceID, i)
	}
}

func TestFindRootLogicalPacketLinkBeneathBaseLogicalPacketLink(t *testing.T) {
	g := getBasicWorkingGraph()
	const baseLPLID = "uuid(gs1/network_node/logical_packet_links/eth0/0)"
	if lpl := graphutil.FindRootLogicalPacketLinkBeneath(g, baseLPLID); lpl != baseLPLID {
		t.Errorf("want: %q, got: %q", baseLPLID, lpl)
	}
}

func TestFindRootLogicalPacketLinkBeneathUpperLogicalPacketLink(t *testing.T) {
	g := getBasicWorkingGraph()
	const baseLPLID = "uuid(gs1/network_node/logical_packet_links/eth0/0)"
	if lpl := graphutil.FindRootLogicalPacketLinkBeneath(g, "uuid(gs1/network_node/logical_packet_links/eth0.mpls/0)"); lpl != baseLPLID {
		t.Errorf("want: %q, got: %q", baseLPLID, lpl)
	}
}

func TestGetAgentIDControllingThisEntity(t *testing.T) {
	g := getBasicWorkingGraph()
	const entityID = "uuid(gs1/platform/modulators/0)"
	const expectedAgentID = "uuid(sdn_agents/0)"
	if agentID, ok := graphutil.GetAgentIDControllingThisEntity(g, entityID); expectedAgentID != agentID || !ok {
		t.Errorf("want: %q, got: %q", expectedAgentID, agentID)
	}
}

func TestGetAgentIDControllingThisEntityReturnsError(t *testing.T) {
	g := getBasicWorkingGraph()
	const entityID = "bogus"
	if _, ok := graphutil.GetAgentIDControllingThisEntity(g, entityID); ok {
		t.Errorf("want ok to be false for bogus entityID")
	}
}

func TestGetAgentIDsFromNetworkInterfaceIDViaModems(t *testing.T) {
	g := getBasicWorkingGraph()
	expectedAgentIDs := []string{"uuid(sdn_agents/0)"}
	if agentIDs := graphutil.GetAgentIDsFromNetworkInterfaceIDViaModems(g, interfaceID); !slices.Equal(expectedAgentIDs, agentIDs) {
		t.Errorf("want: %q, got: %q", expectedAgentIDs, agentIDs)
	}
}

func TestGetAgentIDsFromNetworkInterfaceIDViaModemsReturnsEmptyListIfFails(t *testing.T) {
	g := getBasicWorkingGraph()
	expectedAgentIDs := []string{}
	if agentIDs := graphutil.GetAgentIDsFromNetworkInterfaceIDViaModems(g, bogusInterfaceID); !slices.Equal(expectedAgentIDs, agentIDs) {
		t.Errorf("want: %q, got: %q", expectedAgentIDs, agentIDs)
	}
}

func TestGetPortIDsFromInterfaceID(t *testing.T) {
	g := getBasicWorkingGraph()
	expectedPortIDs := []string{"uuid(gs1/platform/ports/antenna0)", "uuid(gs1/platform/ports/unused_port0)"}
	if portIDs := graphutil.GetPortIDsFromInterfaceID(g, interfaceID); !gcmp.Equal(expectedPortIDs, portIDs, cmpopts.SortSlices(func(a, b string) bool { return a < b })) {
		t.Errorf("want: %q, got: %q", expectedPortIDs, portIDs)
	}
}

func TestGetModulatorIDsFromInterfaceIDViaPort(t *testing.T) {
	g := getBasicWorkingGraph()
	expectedModemIDs := []string{"uuid(gs1/platform/modulators/0)"}
	if modemIDs := graphutil.GetModulatorIDsFromInterfaceIDViaPort(g, interfaceID); !slices.Equal(expectedModemIDs, modemIDs) {
		t.Errorf("want: %q, got: %q", expectedModemIDs, modemIDs)
	}
}

func TestGetModulatorIDsFromInterfaceIDViaPortReturnsEmptyListIfFails(t *testing.T) {
	g := getBasicWorkingGraph()
	expectedModemIDs := []string{}
	if modemIDs := graphutil.GetModulatorIDsFromInterfaceIDViaPort(g, bogusInterfaceID); !slices.Equal(expectedModemIDs, modemIDs) {
		t.Errorf("want: %q, got: %q", expectedModemIDs, modemIDs)
	}
}

func TestGetDemodulatorIDsFromInterfaceIDViaPort(t *testing.T) {
	g := getBasicWorkingGraph()
	expectedModemIDs := []string{"uuid(gs1/platform/demodulators/0)"}
	if modemIDs := graphutil.GetDemodulatorIDsFromInterfaceIDViaPort(g, interfaceID); !slices.Equal(expectedModemIDs, modemIDs) {
		t.Errorf("want: %q, got: %q", expectedModemIDs, modemIDs)
	}
}

func TestGetDemodulatorIDsFromInterfaceIDViaPortReturnsEmptyListIfFails(t *testing.T) {
	g := getBasicWorkingGraph()
	expectedModemIDs := []string{}
	if modemIDs := graphutil.GetDemodulatorIDsFromInterfaceIDViaPort(g, bogusInterfaceID); !slices.Equal(expectedModemIDs, modemIDs) {
		t.Errorf("want: %q, got: %q", expectedModemIDs, modemIDs)
	}
}

// EK_NETWORK_NODE (<--RK_CONTROLS | RK_CONTAINS-->) EK_ROUTE_FN <-- RK_CONTROLS-- EK_SDN_AGENT
// The routeFn that this uses is "uuid(gs1/network_node/ipvrf0)" and that is set up to satisfy both of the conditions above
func TestGetAgentIDFromNetworkNodeIDViaRouteFn(t *testing.T) {
	g := getBasicWorkingGraph()
	const networkInterfaceID = "uuid(gs1/network_node)"
	const expectedAgentID = "uuid(sdn_agents/0)"
	if agentID, ok := graphutil.GetAgentIDFromNetworkNodeIDViaRouteFn(g, networkInterfaceID); expectedAgentID != agentID || !ok {
		t.Errorf("want: %q, got: %q", expectedAgentID, agentID)
	}
}

func TestGetAgentIDFromInterfaceIDViaRouteFn(t *testing.T) {
	g := getBasicWorkingGraph()
	const networkInterfaceID = "uuid(gs1/network_node/interfaces/antenna0)"
	const expectedAgentID = "uuid(sdn_agents/0)"
	if agentID, ok := graphutil.GetAgentIDFromInterfaceIDViaRouteFn(g, networkInterfaceID); expectedAgentID != agentID || !ok {
		t.Errorf("want: %q, got: %q", expectedAgentID, agentID)
	}
}

func TestGetLogicalPacketLinksOriginatingFromInterface(t *testing.T) {
	g := getBasicWorkingGraph()
	const networkInterfaceID = "uuid(gs1/network_node/interfaces/antenna0.mpls)"
	expectedLogicalPacketLinks := []*nmtspb.Entity{
		{
			Id: "uuid(gs1/network_node/logical_packet_links/antenna0.mpls/0)",
			Kind: &nmtspb.Entity_EkLogicalPacketLink{
				EkLogicalPacketLink: &logicalpb.LogicalPacketLink{
					Sr: &logicalpb.LogicalPacketLink_SegmentRouting{
						Enabled: true,
						AdjacencySid: &ietfpb.SegmentId{
							Mpls: 3001,
						},
					},
					MaxDataRateBps: 10737418240,
				},
			},
		},
	}
	logicalPacketLinks := graphutil.GetLogicalPacketLinksOriginatingFromInterface(g, networkInterfaceID)
	if diff := gcmp.Diff(expectedLogicalPacketLinks, logicalPacketLinks, protocmp.Transform()); diff != "" {
		t.Errorf("mismatch (want -> got):\n%s", diff)
	}
}

func TestGetRouteFnsFromInterface(t *testing.T) {
	g := getBasicWorkingGraph()
	const networkInterfaceID = "uuid(gs1/network_node/interfaces/eth0)"
	expectedRouteFns := []*nmtspb.Entity{
		{
			Id: "uuid(gs1/network_node/ipvrf0)",
			Kind: &nmtspb.Entity_EkRouteFn{
				EkRouteFn: &logicalpb.RouteFn{
					Sr: &logicalpb.RouteFn_SegmentRouting{
						Enabled: true,
						NodeSid: &ietfpb.SegmentId{
							Mpls: 1111,
						},
					},
					RouterId: &ietfpb.RouterId{
						Type: &ietfpb.RouterId_DottedQuad{
							DottedQuad: &ietfpb.DottedQuad{
								Str: "1.1.1.1",
							},
						},
					},
				},
			},
		},
	}
	routeFns := graphutil.GetRouteFnsFromInterface(g, networkInterfaceID)
	if diff := gcmp.Diff(expectedRouteFns, routeFns, protocmp.Transform()); diff != "" {
		t.Errorf("mismatch (want -> got):\n%s", diff)
	}
}

func TestGetTransitivelyAffectedIDsForPort(t *testing.T) {
	const portInterfaceTxtpb = `
entity { id: "port"   ek_port{} }
entity { id: "interface"	ek_interface{} }
relationship { a: "interface" kind: RK_TRAVERSES z: "port" }
`
	fragment := lo.Must(testutil.FragmentFrom(portInterfaceTxtpb))
	g := lo.Must(testutil.GraphFromFragments(fragment))
	want := []string{"interface"}
	got := graphutil.ComputeTransitivelyAffectedIDsForFault(g, []string{"port"})
	if diff := gcmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected transitively affected IDs (-want +got): %s", diff)
	}
}
