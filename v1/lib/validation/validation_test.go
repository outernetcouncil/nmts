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
	er "outernetcouncil.org/nmts/v1/lib/entityrelationship"
	"outernetcouncil.org/nmts/v1/lib/graph"
	"outernetcouncil.org/nmts/v1/lib/validation"
	npb "outernetcouncil.org/nmts/v1/proto"
)

type testCase struct {
	entityA string
	rk      npb.RK
	entityZ string
}

var relationshipTestCases = []testCase{
	{
		entityA: `id: "antennaA" ek_antenna{}`,
		rk:      npb.RK_RK_ORIGINATES,
		entityZ: `id: "physical_medium_linkZ" ek_physical_medium_link{}`,
	},
	{
		entityA: `id: "antennaA" ek_antenna{}`,
		rk:      npb.RK_RK_SIGNAL_TRANSITS,
		entityZ: `id: "receiverZ" ek_receiver{}`,
	},
	{
		entityA: `id: "antennaA" ek_antenna{}`,
		rk:      npb.RK_RK_TERMINATES,
		entityZ: `id: "physical_medium_linkZ" ek_physical_medium_link{}`,
	},

	{
		entityA: `id: "interfaceA" ek_interface{}`,
		rk:      npb.RK_RK_DATA_TRANSITS,
		entityZ: `id: "internal_fabricZ" ek_internal_fabric{}`,
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
		entityA: `id: "internal_fabricA" ek_internal_fabric{}`,
		rk:      npb.RK_RK_DATA_TRANSITS,
		entityZ: `id: "interfaceZ" ek_interface{}`,
	},
	{
		entityA: `id: "internal_fabricA" ek_internal_fabric{}`,
		rk:      npb.RK_RK_DATA_TRANSITS,
		entityZ: `id: "internal_fabricZ" ek_internal_fabric{}`,
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
		entityZ: `id: "bp_agent_fnZ" ek_bp_agent_fn{}`,
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
	// Prohibited for now; may be revisited in the future.
	//
	// {
	// entityA: `id: "platformA" ek_platform{}`,
	// rk:      npb.RK_RK_CONTAINS,
	// entityZ: `id: "platformZ" ek_platform{}`,
	// },
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
		rk:      npb.RK_RK_SIGNAL_TRANSITS,
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
		entityZ: `id: "bp_agent_fnZ" ek_bp_agent_fn{}`,
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
		entityA: `id: "sdn_agentA" ek_sdn_agent{}`,
		rk:      npb.RK_RK_CONTROLS,
		entityZ: `id: "platformZ" ek_platform{}`,
	},
	{
		entityA: `id: "sdn_agentA" ek_sdn_agent{}`,
		rk:      npb.RK_RK_CONTROLS,
		entityZ: `id: "route_fnZ" ek_route_fn{}`,
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

	{
		entityA: `id: "transmitterA" ek_transmitter{}`,
		rk:      npb.RK_RK_SUPPORTS,
		entityZ: `id: "carrier_configurationZ" ek_carrier_configuration{}`,
	},
	{
		entityA: `id: "receiverA" ek_receiver{}`,
		rk:      npb.RK_RK_SUPPORTS,
		entityZ: `id: "carrier_configurationZ" ek_carrier_configuration{}`,
	},
}

func TestSimpleCollectionEntityRelationshipValidation(t *testing.T) {
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
		if err := validator.ValidateRelationship(collection, relationship); err != nil {
			t.Errorf("Relationship validation error for %q: %q", relationship, err)
		}
	}
}

func TestSimpleGraphEntityRelationshipValidation(t *testing.T) {
	for _, tc := range relationshipTestCases {
		g := graph.New()
		validator := validation.DefaultGraphValidator{}

		entityA := new(npb.Entity)
		if err := prototext.Unmarshal([]byte(tc.entityA), entityA); err != nil {
			t.Fatalf("failed to parse %q: %q", tc.entityA, err)
		}
		if _, err := g.UpsertEntity(entityA); err != nil {
			t.Fatalf("Failed to add entity %q to graph: %q", tc.entityA, err)
		}
		if err := validation.IsEntityMinimallyWellFormed(entityA); err != nil {
			t.Fatalf("Entity validation error for %q: %q", tc.entityA, err)
		}

		entityZ := new(npb.Entity)
		if err := prototext.Unmarshal([]byte(tc.entityZ), entityZ); err != nil {
			t.Fatalf("failed to parse %q: %q", tc.entityZ, err)
		}
		if _, err := g.UpsertEntity(entityZ); err != nil {
			t.Fatalf("Failed to add entity %q to graph: %q", tc.entityZ, err)
		}
		if err := validation.IsEntityMinimallyWellFormed(entityZ); err != nil {
			t.Fatalf("Entity validation error for %q: %q", tc.entityZ, err)
		}

		relationship := er.Relationship{
			A:    entityA.Id,
			Kind: tc.rk,
			Z:    entityZ.Id,
		}
		if err := validator.ValidateRelationship(g, relationship); err != nil {
			t.Errorf("Relationship validation error for %q: %q", relationship, err)
		}
	}
}

func TestValidateAntenna(t *testing.T) {
	tests := []struct {
		name    string
		entity  string
		wantErr bool
	}{
		{
			name: "antenna with only G/T is valid",
			entity: `id: "my-antenna"
				ek_antenna {
					g_over_t_db_per_k: 10.0
				}`,
			wantErr: false,
		},
		{
			name: "antenna with receive pattern and noise temperature is valid",
			entity: `id: "my-antenna"
				ek_antenna {
					antenna_pattern {
						receive_frequency_range_to_gain_patterns {
							min_frequency: 0
							max_frequency: 300000000000
						}
					}
					antenna_noise_temperature_k: 290.0
				}`,
			wantErr: false,
		},
		{
			name: "antenna with G/T and receive pattern is invalid",
			entity: `id: "my-antenna"
				ek_antenna {
					g_over_t_db_per_k: 10.0
					antenna_pattern {
						receive_frequency_range_to_gain_patterns {
							min_frequency: 0
							max_frequency: 300000000000
						}
					}
				}`,
			wantErr: true,
		},
		{
			name: "antenna with G/T and noise temperature is invalid",
			entity: `id: "my-antenna"
				ek_antenna {
					g_over_t_db_per_k: 10.0
					antenna_noise_temperature_k: 290.0
				}`,
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			entity := new(npb.Entity)
			if err := prototext.Unmarshal([]byte(tc.entity), entity); err != nil {
				t.Fatalf("failed to parse entity: %v", err)
			}
			err := validation.ValidateAntenna(entity)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateAntenna() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestIsEntityMinimallyWellFormedFailsNil(t *testing.T) {
	if err := validation.IsEntityMinimallyWellFormed(nil); err == nil {
		t.Fatalf("failed to minimally invalidate nil")
	}
}

func TestIsEntityMinimallyWellFormedChecksUnicodeNormalizationForm(t *testing.T) {
	// Code points example taken from:
	// https://www.tomdalling.com/blog/coding-tips/when-a-cafe-is-not-a-cafe-a-short-lesson-in-unicode-featuring-nsstring/
	workingAscii := `id: "cafe" ek_network_node{}`
	entity := new(npb.Entity)
	if err := prototext.Unmarshal([]byte(workingAscii), entity); err != nil {
		t.Fatalf("failed to parse %q: %q", workingAscii, err)
	}
	if err := validation.IsEntityMinimallyWellFormed(entity); err != nil {
		t.Fatalf("failed to minimally validate %q: %q", workingAscii, err)
	}

	workingUnicodeAcuteE := `id: "café" ek_network_node{}`
	if err := prototext.Unmarshal([]byte(workingUnicodeAcuteE), entity); err != nil {
		t.Fatalf("failed to parse %q: %q", workingUnicodeAcuteE, err)
	}
	if err := validation.IsEntityMinimallyWellFormed(entity); err != nil {
		t.Fatalf("failed to minimally validate %q: %q", workingUnicodeAcuteE, err)
	}

	brokenUnicodeCombiningAcute := `id: "cafe\u0301" ek_network_node{}`
	if err := prototext.Unmarshal([]byte(brokenUnicodeCombiningAcute), entity); err != nil {
		t.Fatalf("failed to parse %q: %q", brokenUnicodeCombiningAcute, err)
	}
	if err := validation.IsEntityMinimallyWellFormed(entity); err == nil {
		t.Fatalf("failed to minimally validate %q: %q", brokenUnicodeCombiningAcute, err)
	}
}
