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

package validation_test

import (
	"testing"

	"google.golang.org/protobuf/encoding/prototext"
	er "outernetcouncil.org/nmts/lib/entity_relationship"
	"outernetcouncil.org/nmts/lib/validation"
	npb "outernetcouncil.org/nmts/proto"
)

type testCase struct {
	entityA string
	rk      npb.RK
	entityZ string
}

var relationshipTestCases = []testCase{
	{
		entityA: `id: "antennaA" ek_antenna{}`,
		rk:      npb.RK_RK_SIGNAL_TRANSITS,
		entityZ: `id: "receiverZ" ek_receiver{}`,
	},

	{
		entityA: `id: "interfaceA" ek_interface{}`,
		rk:      npb.RK_RK_ORIGINATES,
		entityZ: `id: "logical_packet_linkZ" ek_logical_packet_link{}`,
	},
	{
		entityA: `id: "interfaceA" ek_interface{}`,
		rk:      npb.RK_RK_TERMINATES,
		entityZ: `id: "logical_packet_linkZ" ek_logical_packet_link{}`,
	},
	{
		entityA: `id: "interfaceA" ek_interface{}`,
		rk:      npb.RK_RK_TRAVERSES,
		entityZ: `id: "portZ" ek_port{}`,
	},
	{
		entityA: `id: "interfaceA" ek_interface{}`,
		rk:      npb.RK_RK_TRAVERSES,
		entityZ: `id: "interfaceZ" ek_interface{}`,
	},

	{
		entityA: `id: "logical_packet_linkA" ek_logical_packet_link{}`,
		rk:      npb.RK_RK_TRAVERSES,
		entityZ: `id: "physical_medium_linkZ" ek_physical_medium_link{}`,
	},
	{
		entityA: `id: "logical_packet_linkA" ek_logical_packet_link{}`,
		rk:      npb.RK_RK_TRAVERSES,
		entityZ: `id: "logical_packet_linkZ" ek_logical_packet_link{}`,
	},

	{
		entityA: `id: "modulatorA" ek_modulator{}`,
		rk:      npb.RK_RK_SIGNAL_TRANSITS,
		entityZ: `id: "signal_processing_chainZ" ek_signal_processing_chain{}`,
	},

	{
		entityA: `id: "network_nodeA" ek_network_node{}`,
		rk:      npb.RK_RK_CONTAINS,
		entityZ: `id: "interfaceZ" ek_interface{}`,
	},
	{
		entityA: `id: "network_nodeA" ek_network_node{}`,
		rk:      npb.RK_RK_CONTAINS,
		entityZ: `id: "route_fnZ" ek_route_fn{}`,
	},
	{
		entityA: `id: "network_nodeA" ek_network_node{}`,
		rk:      npb.RK_RK_CONTAINS,
		entityZ: `id: "switch_fnZ" ek_switch_fn{}`,
	},
	{
		entityA: `id: "network_nodeA" ek_network_node{}`,
		rk:      npb.RK_RK_CONTAINS,
		entityZ: `id: "sdn_agentZ" ek_sdn_agent{}`,
	},

	{
		entityA: `id: "platformA" ek_platform{}`,
		rk:      npb.RK_RK_CONTAINS,
		entityZ: `id: "antennaZ" ek_antenna{}`,
	},
	{
		entityA: `id: "platformA" ek_platform{}`,
		rk:      npb.RK_RK_CONTAINS,
		entityZ: `id: "demodulatorZ" ek_demodulator{}`,
	},
	{
		entityA: `id: "platformA" ek_platform{}`,
		rk:      npb.RK_RK_CONTAINS,
		entityZ: `id: "modulatorZ" ek_modulator{}`,
	},
	{
		entityA: `id: "platformA" ek_platform{}`,
		rk:      npb.RK_RK_CONTAINS,
		entityZ: `id: "platformZ" ek_platform{}`,
	},
	{
		entityA: `id: "platformA" ek_platform{}`,
		rk:      npb.RK_RK_CONTAINS,
		entityZ: `id: "portZ" ek_port{}`,
	},
	{
		entityA: `id: "platformA" ek_platform{}`,
		rk:      npb.RK_RK_CONTAINS,
		entityZ: `id: "network_nodeZ" ek_network_node{}`,
	},
	{
		entityA: `id: "platformA" ek_platform{}`,
		rk:      npb.RK_RK_CONTAINS,
		entityZ: `id: "receiverZ" ek_receiver{}`,
	},
	{
		entityA: `id: "platformA" ek_platform{}`,
		rk:      npb.RK_RK_CONTAINS,
		entityZ: `id: "signal_processing_chainZ" ek_signal_processing_chain{}`,
	},
	{
		entityA: `id: "platformA" ek_platform{}`,
		rk:      npb.RK_RK_CONTAINS,
		entityZ: `id: "transmitterZ" ek_transmitter{}`,
	},

	{
		entityA: `id: "portA" ek_port{}`,
		rk:      npb.RK_RK_ORIGINATES,
		entityZ: `id: "modulatorZ" ek_modulator{}`,
	},
	{
		entityA: `id: "portA" ek_port{}`,
		rk:      npb.RK_RK_ORIGINATES,
		entityZ: `id: "physical_medium_linkZ" ek_physical_medium_link{}`,
	},
	{
		entityA: `id: "portA" ek_port{}`,
		rk:      npb.RK_RK_TERMINATES,
		entityZ: `id: "demodulatorZ" ek_demodulator{}`,
	},
	{
		entityA: `id: "portA" ek_port{}`,
		rk:      npb.RK_RK_TERMINATES,
		entityZ: `id: "physical_medium_linkZ" ek_physical_medium_link{}`,
	},

	{
		entityA: `id: "receiverA" ek_receiver{}`,
		rk:      npb.RK_RK_TERMINATES,
		entityZ: `id: "signal_processing_chainZ" ek_signal_processing_chain{}`,
	},

	{
		entityA: `id: "sdn_agentA" ek_sdn_agent{}`,
		rk:      npb.RK_RK_CONTROLS,
		entityZ: `id: "antennaZ" ek_antenna{}`,
	},
	{
		entityA: `id: "sdn_agentA" ek_sdn_agent{}`,
		rk:      npb.RK_RK_CONTROLS,
		entityZ: `id: "demodulatorZ" ek_demodulator{}`,
	},
	{
		entityA: `id: "sdn_agentA" ek_sdn_agent{}`,
		rk:      npb.RK_RK_CONTROLS,
		entityZ: `id: "modulatorZ" ek_modulator{}`,
	},
	{
		entityA: `id: "sdn_agentA" ek_sdn_agent{}`,
		rk:      npb.RK_RK_CONTROLS,
		entityZ: `id: "network_nodeZ" ek_network_node{}`,
	},

	{
		entityA: `id: "signal_processing_chainA" ek_signal_processing_chain{}`,
		rk:      npb.RK_RK_SIGNAL_TRANSITS,
		entityZ: `id: "demodulatorZ" ek_demodulator{}`,
	},
	{
		entityA: `id: "signal_processing_chainA" ek_signal_processing_chain{}`,
		rk:      npb.RK_RK_SIGNAL_TRANSITS,
		entityZ: `id: "signal_processing_chainZ" ek_signal_processing_chain{}`,
	},
	{
		entityA: `id: "signal_processing_chainA" ek_signal_processing_chain{}`,
		rk:      npb.RK_RK_SIGNAL_TRANSITS,
		entityZ: `id: "transmitterZ" ek_transmitter{}`,
	},

	{
		entityA: `id: "transmitterA" ek_transmitter{}`,
		rk:      npb.RK_RK_SIGNAL_TRANSITS,
		entityZ: `id: "antennaZ" ek_antenna{}`,
	},
}

func TestEntityKindStringExamples(t *testing.T) {
	for _, tc := range relationshipTestCases {
		collection := er.NewCollection()
		validator := validation.DefaultValidator{}

		entityA := new(npb.Entity)
		if err := prototext.Unmarshal([]byte(tc.entityA), entityA); err != nil {
			t.Fatalf("failed to parse %q: %q", tc.entityA, err)
		}
		if err := collection.InsertEntity(entityA); err != nil {
			t.Fatalf("Failed to add entity %q to colleciton: %q", tc.entityA, err)
		}
		if err := validator.ValidateEntity(collection, entityA); err != nil {
			t.Fatalf("Entity validation error for %q: %q", tc.entityA, err)
		}

		entityZ := new(npb.Entity)
		if err := prototext.Unmarshal([]byte(tc.entityZ), entityZ); err != nil {
			t.Fatalf("failed to parse %q: %q", tc.entityZ, err)
		}
		if err := collection.InsertEntity(entityZ); err != nil {
			t.Fatalf("Failed to add entity %q to colleciton: %q", tc.entityZ, err)
		}
		if err := validator.ValidateEntity(collection, entityZ); err != nil {
			t.Fatalf("Entity validation error for %q: %q", tc.entityZ, err)
		}

		relationship := er.Relationship{
			A:    entityA.Id,
			Kind: tc.rk,
			Z:    entityZ.Id,
		}
		validator.ValidateRelationship(collection, relationship)
	}
}
