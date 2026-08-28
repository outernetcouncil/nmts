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
	"errors"
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
		entityZ: `id: "internal_fabricZ" ek_internal_fabric{}`,
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
		{
			name: "eirp_limits with fixed psd and frequency range is valid",
			entity: `id: "my-antenna"
				ek_antenna {
					eirp_limits {
						eirpsd_masks {
							frequency_range { min_frequency_hz: 0 max_frequency_hz: 300000000000 }
							power_spectral_density {
								reference_bandwidth_hz: 1000000
								fixed { power_dbw: 20 }
							}
						}
					}
				}`,
			wantErr: false,
		},
		{
			name: "eirp_limits with ordered off_axis control points is valid",
			entity: `id: "my-antenna"
				ek_antenna {
					eirp_limits {
						eirpsd_masks {
							power_spectral_density {
								reference_bandwidth_hz: 1000000
								off_axis {
									control_points { angle_deg: 0 power_dbw: 20 }
									control_points { angle_deg: 5 power_dbw: 10 }
								}
							}
						}
					}
				}`,
			wantErr: false,
		},
		{
			name: "eirp_limits mask with zero reference_bandwidth_hz is invalid",
			entity: `id: "my-antenna"
				ek_antenna {
					eirp_limits {
						eirpsd_masks {
							power_spectral_density {
								reference_bandwidth_hz: 0
								fixed { power_dbw: 20 }
							}
						}
					}
				}`,
			wantErr: true,
		},
		{
			name: "eirp_limits mask without power_spectral_density is invalid",
			entity: `id: "my-antenna"
				ek_antenna {
					eirp_limits {
						eirpsd_masks {
							frequency_range { min_frequency_hz: 0 max_frequency_hz: 1000 }
						}
					}
				}`,
			wantErr: true,
		},
		{
			name: "eirp_limits off_axis with one control point is invalid",
			entity: `id: "my-antenna"
				ek_antenna {
					eirp_limits {
						eirpsd_masks {
							power_spectral_density {
								reference_bandwidth_hz: 1000000
								off_axis {
									control_points { angle_deg: 0 power_dbw: 20 }
								}
							}
						}
					}
				}`,
			wantErr: true,
		},
		{
			name: "eirp_limits off_axis with out-of-order angles is invalid",
			entity: `id: "my-antenna"
				ek_antenna {
					eirp_limits {
						eirpsd_masks {
							power_spectral_density {
								reference_bandwidth_hz: 1000000
								off_axis {
									control_points { angle_deg: 5 power_dbw: 10 }
									control_points { angle_deg: 0 power_dbw: 20 }
								}
							}
						}
					}
				}`,
			wantErr: true,
		},
		{
			name: "eirp_limits frequency_range with max not greater than min is invalid",
			entity: `id: "my-antenna"
				ek_antenna {
					eirp_limits {
						eirpsd_masks {
							frequency_range { min_frequency_hz: 1000 max_frequency_hz: 1000 }
							power_spectral_density {
								reference_bandwidth_hz: 1000000
								fixed { power_dbw: 20 }
							}
						}
					}
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

func TestIsEntityIDMinimallyWellFormed(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name: "ascii id is valid",
			id:   "cafe",
		},
		{
			name: "precomposed unicode id is valid",
			id:   "café",
		},
		{
			name:    "empty id is invalid",
			id:      "",
			wantErr: true,
		},
		{
			name:    "id with leading whitespace is invalid",
			id:      " cafe",
			wantErr: true,
		},
		{
			name:    "id with trailing whitespace is invalid",
			id:      "cafe ",
			wantErr: true,
		},
		{
			name:    "id not in Unicode Normalization Form C is invalid",
			id:      "cafe\u0301",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validation.IsEntityIDMinimallyWellFormed(tc.id)
			if (err != nil) != tc.wantErr {
				t.Errorf("IsEntityIDMinimallyWellFormed() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestIsRelationshipMinimallyWellFormed(t *testing.T) {
	tests := []struct {
		name         string
		relationship string
		wantErr      bool
	}{
		{
			name:         "fully specified relationship is valid",
			relationship: `a: "platform" kind: RK_CONTAINS z: "port"`,
		},
		{
			name:         "empty a endpoint is invalid",
			relationship: `kind: RK_CONTAINS z: "port"`,
			wantErr:      true,
		},
		{
			name:         "empty z endpoint is invalid",
			relationship: `a: "platform" kind: RK_CONTAINS`,
			wantErr:      true,
		},
		{
			name:         "a endpoint with surrounding whitespace is invalid",
			relationship: `a: " platform " kind: RK_CONTAINS z: "port"`,
			wantErr:      true,
		},
		{
			name:         "z endpoint not in Unicode Normalization Form C is invalid",
			relationship: `a: "platform" kind: RK_CONTAINS z: "cafe\u0301"`,
			wantErr:      true,
		},
		{
			name:         "self-referential relationship is invalid",
			relationship: `a: "platform" kind: RK_CONTAINS z: "platform"`,
			wantErr:      true,
		},
		{
			name:         "unspecified kind is invalid",
			relationship: `a: "platform" z: "port"`,
			wantErr:      true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			relationship := new(npb.Relationship)
			if err := prototext.Unmarshal([]byte(tc.relationship), relationship); err != nil {
				t.Fatalf("failed to parse relationship: %v", err)
			}
			err := validation.IsRelationshipMinimallyWellFormed(relationship)
			if (err != nil) != tc.wantErr {
				t.Errorf("IsRelationshipMinimallyWellFormed() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestIsRelationshipMinimallyWellFormedFailsNil(t *testing.T) {
	if err := validation.IsRelationshipMinimallyWellFormed(nil); err == nil {
		t.Fatalf("failed to minimally invalidate nil")
	}
}

// countErrors reports how many problems err holds, unwrapping the error
// errors.Join returns when more than one problem was found.
func countErrors(err error) int {
	if err == nil {
		return 0
	}
	var joined interface{ Unwrap() []error }
	if errors.As(err, &joined) {
		return len(joined.Unwrap())
	}
	return 1
}

func TestValidateFragment(t *testing.T) {
	tests := []struct {
		name     string
		fragment string
		wantErrs int
	}{
		{
			name: "fragment of well formed, distinct elements is valid",
			fragment: `
				entity { id: "platform" ek_platform{} }
				entity { id: "port" ek_port{} }
				relationship { a: "platform" kind: RK_CONTAINS z: "port" }`,
		},
		{
			name:     "empty fragment is valid",
			fragment: ``,
		},
		{
			name:     "entity lacking a kind is invalid",
			fragment: `entity { id: "platform" }`,
			wantErrs: 1,
		},
		{
			name:     "entity with an empty id is invalid",
			fragment: `entity { ek_platform{} }`,
			wantErrs: 1,
		},
		{
			name: "repeated entity id is invalid",
			fragment: `
				entity { id: "platform" ek_platform{} }
				entity { id: "platform" ek_platform{} }`,
			wantErrs: 1,
		},
		{
			name: "relationship lacking a kind is invalid",
			fragment: `
				entity { id: "platform" ek_platform{} }
				entity { id: "port" ek_port{} }
				relationship { a: "platform" z: "port" }`,
			wantErrs: 1,
		},
		{
			name: "self-referential relationship is invalid",
			fragment: `
				entity { id: "fabric" ek_internal_fabric{} }
				relationship { a: "fabric" kind: RK_DATA_TRANSITS z: "fabric" }`,
			wantErrs: 1,
		},
		{
			name: "repeated relationship tuple is invalid",
			fragment: `
				entity { id: "platform" ek_platform{} }
				entity { id: "port" ek_port{} }
				relationship { a: "platform" kind: RK_CONTAINS z: "port" }
				relationship { a: "platform" kind: RK_CONTAINS z: "port" }`,
			wantErrs: 1,
		},
		{
			// Validation does not stop at the first problem; every Entity and
			// every Relationship is checked and all errors are reported.
			name: "every problem in the fragment is reported",
			fragment: `
				entity { id: "platform" }
				entity { id: "platform" ek_platform{} }
				relationship { a: "platform" kind: RK_CONTAINS z: "platform" }`,
			wantErrs: 3,
		},
		{
			name: "relationships between the same pair differing in kind are valid",
			fragment: `
				entity { id: "port" ek_port{} }
				entity { id: "link" ek_physical_medium_link{} }
				relationship { a: "port" kind: RK_ORIGINATES z: "link" }
				relationship { a: "port" kind: RK_TERMINATES z: "link" }`,
		},
		{
			// A Fragment is a subgraph: its relationships may refer to Entities
			// held elsewhere in the model, so endpoints are not resolved here.
			name:     "relationship referring to entities outside the fragment is valid",
			fragment: `relationship { a: "platform" kind: RK_CONTAINS z: "port" }`,
		},
		{
			// Whether a kind is permitted between two Entities' kinds requires
			// the surrounding Collection or Graph; DefaultValidator rejects this
			// pairing, ValidateFragment does not.
			name: "relationship between an impermissible pair of entity kinds is valid",
			fragment: `
				entity { id: "platformA" ek_platform{} }
				entity { id: "platformZ" ek_platform{} }
				relationship { a: "platformA" kind: RK_CONTAINS z: "platformZ" }`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fragment := new(npb.Fragment)
			if err := prototext.Unmarshal([]byte(tc.fragment), fragment); err != nil {
				t.Fatalf("failed to parse fragment: %v", err)
			}
			err := validation.ValidateFragment(fragment)
			if got := countErrors(err); got != tc.wantErrs {
				t.Errorf("ValidateFragment() error = %v, want %d errors", err, tc.wantErrs)
			}
		})
	}
}

func TestValidateFragmentAcceptsNil(t *testing.T) {
	if err := validation.ValidateFragment(nil); err != nil {
		t.Errorf("ValidateFragment(nil) = %v, want no error", err)
	}
}
