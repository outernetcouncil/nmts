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

package validation

import (
	"fmt"
	"strings"

	"golang.org/x/text/unicode/norm"
	er "outernetcouncil.org/nmts/v1/lib/entityrelationship"
	"outernetcouncil.org/nmts/v1/lib/graph"
	npb "outernetcouncil.org/nmts/v1/proto"
	physicalpb "outernetcouncil.org/nmts/v1/proto/types/physical"
)

func IsEntityMinimallyWellFormed(entity *npb.Entity) error {
	if entity == nil {
		return fmt.Errorf("entity MUST NOT be nil")
	}
	id := entity.Id

	// In keeping with https://google.aip.dev/210#normalization
	// ensure Entity IDs are in Unicode Normal Form C.
	if norm.NFC.String(id) != id {
		return fmt.Errorf("entity ID MUST be in Unicode Normalization Form C")
	}

	// Do not permit extraneous whitespace; this likely indicates
	// some configuration or tooling error.
	if id != strings.TrimSpace(id) {
		return fmt.Errorf("id must not have lead nor trailing whitespace: '%s'", id)
	}

	// TODO: uuid.Validate(id)
	if id == "" {
		return fmt.Errorf("id must not be empty: %q", id)
	}

	if er.EntityKindStringFromProto(entity) == "" {
		return fmt.Errorf("Entity %q lacks an entity kind field", id)
	}

	return nil
}

func ValidateAntenna(entity *npb.Entity) error {
	antenna := entity.GetEkAntenna()
	if antenna == nil {
		return nil
	}

	// Checks that the model of the antenna's receive performance is valid.
	// Configurations that include G/T with either a receive antenna pattern or
	// sources of noise are invalid.
	hasGOverT := antenna.GetGOverTDbPerK() != 0
	hasReceiveAntennaPattern := len(antenna.GetAntennaPattern().GetReceiveFrequencyRangeToGainPatterns()) > 0
	hasNoiseTemperature := antenna.GetAntennaNoiseTemperatureK() != 0
	if hasGOverT && (hasReceiveAntennaPattern || hasNoiseTemperature) {
		return fmt.Errorf(
			"Antenna %q: configurations that include G/T with either a "+
				"receive antenna pattern or sources of noise are invalid",
			entity.GetId(),
		)
	}

	if err := validateEirpLimits(antenna.GetEirpLimits(), entity.GetId()); err != nil {
		return err
	}

	for i, mask := range antenna.GetEmissionEnvelope().GetEirpsdMasks() {
		if err := validateEirpsdMask(mask); err != nil {
			return fmt.Errorf("antenna %q: emission_envelope.eirpsd_masks[%d]: %w", entity.GetId(), i, err)
		}
	}

	return nil
}

// validateEirpLimits checks an antenna's eirp_limits and its nested EIRPSD masks.
// eirp_limits is optional, so a nil value is valid.
func validateEirpLimits(limits *physicalpb.EirpLimits, id string) error {
	if limits == nil {
		return nil
	}
	for i, mask := range limits.GetEirpsdMasks() {
		if err := validateEirpsdMask(mask); err != nil {
			return fmt.Errorf("Antenna %q: eirp_limits.eirpsd_masks[%d]: %w", id, i, err)
		}
	}
	return nil
}

func validateEirpsdMask(mask *physicalpb.EirpsdMask) error {
	psd := mask.GetPowerSpectralDensity()
	if psd == nil {
		return fmt.Errorf("power_spectral_density is required")
	}
	// frequency_range is optional (an unset range applies to all frequencies), but
	// when present it must be a non-empty half-open [min, max) interval.
	if fr := mask.GetFrequencyRange(); fr != nil {
		if fr.GetMinFrequencyHz() < 0 {
			return fmt.Errorf("frequency_range.min_frequency_hz must be non-negative: %d", fr.GetMinFrequencyHz())
		}
		if fr.GetMaxFrequencyHz() <= fr.GetMinFrequencyHz() {
			return fmt.Errorf(
				"frequency_range.max_frequency_hz (%d) must be greater than min_frequency_hz (%d)",
				fr.GetMaxFrequencyHz(), fr.GetMinFrequencyHz())
		}
	}
	return validatePowerSpectralDensity(psd)
}

func validatePowerSpectralDensity(psd *physicalpb.PowerSpectralDensity) error {
	if psd.GetReferenceBandwidthHz() <= 0.0 {
		return fmt.Errorf(
			"power_spectral_density.reference_bandwidth_hz must be positive: %v",
			psd.GetReferenceBandwidthHz())
	}
	switch t := psd.GetType().(type) {
	case *physicalpb.PowerSpectralDensity_Fixed:
		// Any finite power_dbw is valid; proto3 cannot distinguish an unset scalar from 0.
		return nil
	case *physicalpb.PowerSpectralDensity_OffAxis:
		return validateOffAxisPower(t.OffAxis)
	case nil:
		return fmt.Errorf("power_spectral_density.type must be set")
	default:
		return fmt.Errorf("power_spectral_density has unknown type: %T", t)
	}
}

func validateOffAxisPower(o *physicalpb.PowerSpectralDensity_OffAxisPower) error {
	points := o.GetControlPoints()
	if len(points) < 2 {
		return fmt.Errorf("off_axis.control_points must have at least two entries, got %d", len(points))
	}
	for i := 1; i < len(points); i++ {
		if points[i].GetAngleDeg() < points[i-1].GetAngleDeg() {
			return fmt.Errorf(
				"off_axis.control_points must be ordered by non-decreasing angle_deg: %v after %v",
				points[i].GetAngleDeg(), points[i-1].GetAngleDeg())
		}
	}
	return nil
}

type DefaultValidator struct{}

// Validate each entity as it's loaded within the collection context
// assembled up to that point.
func (DefaultValidator) ValidateEntity(coll *er.Collection, entity *npb.Entity) error {
	if err := IsEntityMinimallyWellFormed(entity); err != nil {
		return err
	}
	return ValidateAntenna(entity)
}

type allowedRelationship struct {
	a, z string
	rk   npb.RK
}

func (aRel allowedRelationship) String() string {
	return fmt.Sprintf("%s->%s->%s", aRel.a, aRel.rk, aRel.z)
}

var permittedRelationships = map[allowedRelationship]struct{}{
	{a: "EK_ANTENNA", rk: npb.RK_RK_ORIGINATES, z: "EK_PHYSICAL_MEDIUM_LINK"}: {},
	{a: "EK_ANTENNA", rk: npb.RK_RK_SIGNAL_TRANSITS, z: "EK_RECEIVER"}:        {},
	{a: "EK_ANTENNA", rk: npb.RK_RK_TERMINATES, z: "EK_PHYSICAL_MEDIUM_LINK"}: {},

	{a: "EK_INTERFACE", rk: npb.RK_RK_DATA_TRANSITS, z: "EK_INTERNAL_FABRIC"}:  {},
	{a: "EK_INTERFACE", rk: npb.RK_RK_ORIGINATES, z: "EK_LOGICAL_PACKET_LINK"}: {},
	{a: "EK_INTERFACE", rk: npb.RK_RK_TERMINATES, z: "EK_LOGICAL_PACKET_LINK"}: {},
	{a: "EK_INTERFACE", rk: npb.RK_RK_TRAVERSES, z: "EK_PORT"}:                 {},
	{a: "EK_INTERFACE", rk: npb.RK_RK_TRAVERSES, z: "EK_INTERFACE"}:            {},

	{a: "EK_INTERNAL_FABRIC", rk: npb.RK_RK_DATA_TRANSITS, z: "EK_INTERFACE"}:       {},
	{a: "EK_INTERNAL_FABRIC", rk: npb.RK_RK_DATA_TRANSITS, z: "EK_INTERNAL_FABRIC"}: {},

	{a: "EK_LOGICAL_PACKET_LINK", rk: npb.RK_RK_TRAVERSES, z: "EK_PHYSICAL_MEDIUM_LINK"}: {},
	{a: "EK_LOGICAL_PACKET_LINK", rk: npb.RK_RK_TRAVERSES, z: "EK_LOGICAL_PACKET_LINK"}:  {},

	{a: "EK_MODULATOR", rk: npb.RK_RK_SIGNAL_TRANSITS, z: "EK_SIGNAL_PROCESSING_CHAIN"}: {},

	{a: "EK_NETWORK_NODE", rk: npb.RK_RK_CONTAINS, z: "EK_ACCESS_FN"}:       {},
	{a: "EK_NETWORK_NODE", rk: npb.RK_RK_CONTAINS, z: "EK_BP_AGENT_FN"}:     {},
	{a: "EK_NETWORK_NODE", rk: npb.RK_RK_CONTAINS, z: "EK_INTERFACE"}:       {},
	{a: "EK_NETWORK_NODE", rk: npb.RK_RK_CONTAINS, z: "EK_INTERNAL_FABRIC"}: {},
	{a: "EK_NETWORK_NODE", rk: npb.RK_RK_CONTAINS, z: "EK_ROUTE_FN"}:        {},
	{a: "EK_NETWORK_NODE", rk: npb.RK_RK_CONTAINS, z: "EK_SWITCH_FN"}:       {},
	{a: "EK_NETWORK_NODE", rk: npb.RK_RK_CONTAINS, z: "EK_SDN_AGENT"}:       {},

	{a: "EK_PLATFORM", rk: npb.RK_RK_CONTAINS, z: "EK_ANTENNA"}:     {},
	{a: "EK_PLATFORM", rk: npb.RK_RK_CONTAINS, z: "EK_DEMODULATOR"}: {},
	{a: "EK_PLATFORM", rk: npb.RK_RK_CONTAINS, z: "EK_MODULATOR"}:   {},
	// It's theoretically possible to model a system as a platform of
	// platforms. This can, however, have complicating implications
	// for some code that might try to examine and enforce certain
	// relationships or constraints.
	//
	// Leave this commented out, as it will likely be revisited.
	//
	// {a: "EK_PLATFORM", rk: npb.RK_RK_CONTAINS, z: "EK_PLATFORM"}:                {},
	{a: "EK_PLATFORM", rk: npb.RK_RK_CONTAINS, z: "EK_PORT"}:                    {},
	{a: "EK_PLATFORM", rk: npb.RK_RK_CONTAINS, z: "EK_NETWORK_NODE"}:            {},
	{a: "EK_PLATFORM", rk: npb.RK_RK_CONTAINS, z: "EK_RECEIVER"}:                {},
	{a: "EK_PLATFORM", rk: npb.RK_RK_CONTAINS, z: "EK_SIGNAL_PROCESSING_CHAIN"}: {},
	{a: "EK_PLATFORM", rk: npb.RK_RK_CONTAINS, z: "EK_TRANSMITTER"}:             {},

	{a: "EK_PORT", rk: npb.RK_RK_ORIGINATES, z: "EK_MODULATOR"}:            {},
	{a: "EK_PORT", rk: npb.RK_RK_ORIGINATES, z: "EK_PHYSICAL_MEDIUM_LINK"}: {},
	{a: "EK_PORT", rk: npb.RK_RK_TERMINATES, z: "EK_DEMODULATOR"}:          {},
	{a: "EK_PORT", rk: npb.RK_RK_TERMINATES, z: "EK_PHYSICAL_MEDIUM_LINK"}: {},

	{a: "EK_RECEIVER", rk: npb.RK_RK_SIGNAL_TRANSITS, z: "EK_SIGNAL_PROCESSING_CHAIN"}: {},

	{a: "EK_ROUTE_FN", rk: npb.RK_RK_CONTROLS, z: "EK_NETWORK_NODE"}: {},

	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_ACCESS_FN"}:    {},
	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_ANTENNA"}:      {},
	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_BP_AGENT_FN"}:  {},
	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_DEMODULATOR"}:  {},
	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_MODULATOR"}:    {},
	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_NETWORK_NODE"}: {},
	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_PLATFORM"}:     {},
	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_ROUTE_FN"}:     {},

	{a: "EK_SIGNAL_PROCESSING_CHAIN", rk: npb.RK_RK_SIGNAL_TRANSITS, z: "EK_DEMODULATOR"}:             {},
	{a: "EK_SIGNAL_PROCESSING_CHAIN", rk: npb.RK_RK_SIGNAL_TRANSITS, z: "EK_SIGNAL_PROCESSING_CHAIN"}: {},
	{a: "EK_SIGNAL_PROCESSING_CHAIN", rk: npb.RK_RK_SIGNAL_TRANSITS, z: "EK_TRANSMITTER"}:             {},

	{a: "EK_TRANSMITTER", rk: npb.RK_RK_SIGNAL_TRANSITS, z: "EK_ANTENNA"}: {},

	{a: "EK_TRANSMITTER", rk: npb.RK_RK_SUPPORTS, z: "EK_CARRIER_CONFIGURATION"}: {},
	{a: "EK_RECEIVER", rk: npb.RK_RK_SUPPORTS, z: "EK_CARRIER_CONFIGURATION"}:    {},
}

// Validate each relationship as it's loaded within the collection
// context assembled up to that point.
func (DefaultValidator) ValidateRelationship(coll *er.Collection, rel er.Relationship) error {
	kindA := er.EntityKindStringFromProto(coll.Entities[rel.A])
	kindZ := er.EntityKindStringFromProto(coll.Entities[rel.Z])

	key := allowedRelationship{a: kindA, rk: rel.Kind, z: kindZ}
	if _, ok := permittedRelationships[key]; ok {
		return nil
	}

	// More detailed checks can be added here, after basic validity
	// has been checked and before the "default deny" error.

	return fmt.Errorf("unsupported relationship between entites: '%v' i.e. '%v'", rel.String(), key)
}

// Validate the complete collection.
func (DefaultValidator) ValidateCollection(coll *er.Collection) error {
	if coll.NumEntities() < 1 {
		return fmt.Errorf("found no entities")
	}
	if coll.NumRelationships() < 1 {
		return fmt.Errorf("found no relationships")
	}
	return nil
}

type DefaultGraphValidator struct{}

func (DefaultGraphValidator) ValidateRelationship(g *graph.Graph, rel er.Relationship) error {
	a := g.Node(rel.A)
	if a == nil {
		return fmt.Errorf("A entity %q does not exist", rel.A)
	}
	kindA := a.GetKind()

	z := g.Node(rel.Z)
	if z == nil {
		return fmt.Errorf("Z entity %q does not exist", rel.Z)
	}
	kindZ := z.GetKind()

	key := allowedRelationship{a: kindA, rk: rel.Kind, z: kindZ}
	if _, ok := permittedRelationships[key]; ok {
		return nil
	}

	// More detailed checks can be added here, after basic validity
	// has been checked and before the "default deny" error.

	return fmt.Errorf("unsupported relationship between entites: '%v' i.e. '%v'", rel.String(), key)
}
