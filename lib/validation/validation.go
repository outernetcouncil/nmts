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

	er "outernetcouncil.org/nmts/lib/entityrelationship"
	npb "outernetcouncil.org/nmts/proto"
)

type DefaultValidator struct {
}

// Validate each entity as it's loaded within the collection context
// assembled up to that point.
func (DefaultValidator) ValidateEntity(coll *er.Collection, entity *npb.Entity) error {
	id := entity.Id

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

type allowedRelationship = struct {
	a, z string
	rk   npb.RK
}

var permittedRelationships = map[allowedRelationship]struct{}{
	{a: "EK_ANTENNA", rk: npb.RK_RK_SIGNAL_TRANSITS, z: "EK_RECEIVER"}: {},

	{a: "EK_INTERFACE", rk: npb.RK_RK_ORIGINATES, z: "EK_LOGICAL_PACKET_LINK"}: {},
	{a: "EK_INTERFACE", rk: npb.RK_RK_TERMINATES, z: "EK_LOGICAL_PACKET_LINK"}: {},
	{a: "EK_INTERFACE", rk: npb.RK_RK_TRAVERSES, z: "EK_PORT"}:                 {},
	{a: "EK_INTERFACE", rk: npb.RK_RK_TRAVERSES, z: "EK_INTERFACE"}:            {},

	{a: "EK_LOGICAL_PACKET_LINK", rk: npb.RK_RK_TRAVERSES, z: "EK_PHYSICAL_MEDIUM_LINK"}: {},
	{a: "EK_LOGICAL_PACKET_LINK", rk: npb.RK_RK_TRAVERSES, z: "EK_LOGICAL_PACKET_LINK"}:  {},

	{a: "EK_MODULATOR", rk: npb.RK_RK_SIGNAL_TRANSITS, z: "EK_SIGNAL_PROCESSING_CHAIN"}: {},

	{a: "EK_NETWORK_NODE", rk: npb.RK_RK_CONTAINS, z: "EK_INTERFACE"}: {},
	{a: "EK_NETWORK_NODE", rk: npb.RK_RK_CONTAINS, z: "EK_ROUTE_FN"}:  {},
	{a: "EK_NETWORK_NODE", rk: npb.RK_RK_CONTAINS, z: "EK_SWITCH_FN"}: {},
	{a: "EK_NETWORK_NODE", rk: npb.RK_RK_CONTAINS, z: "EK_SDN_AGENT"}: {},

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

	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_ANTENNA"}:      {},
	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_DEMODULATOR"}:  {},
	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_MODULATOR"}:    {},
	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_NETWORK_NODE"}: {},
	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_PLATFORM"}:     {},
	{a: "EK_SDN_AGENT", rk: npb.RK_RK_CONTROLS, z: "EK_ROUTE_FN"}:     {},

	{a: "EK_SIGNAL_PROCESSING_CHAIN", rk: npb.RK_RK_SIGNAL_TRANSITS, z: "EK_DEMODULATOR"}:             {},
	{a: "EK_SIGNAL_PROCESSING_CHAIN", rk: npb.RK_RK_SIGNAL_TRANSITS, z: "EK_SIGNAL_PROCESSING_CHAIN"}: {},
	{a: "EK_SIGNAL_PROCESSING_CHAIN", rk: npb.RK_RK_SIGNAL_TRANSITS, z: "EK_TRANSMITTER"}:             {},

	{a: "EK_TRANSMITTER", rk: npb.RK_RK_SIGNAL_TRANSITS, z: "EK_ANTENNA"}: {},
}

// Validate each relationship as it's loaded within the collection
// context assembled up to that point.
func (DefaultValidator) ValidateRelationship(coll *er.Collection, rel er.Relationship) error {
	kindA := er.EntityKindStringFromProto(coll.Entities[rel.A])
	kindZ := er.EntityKindStringFromProto(coll.Entities[rel.Z])

	if _, ok := permittedRelationships[allowedRelationship{a: kindA, rk: rel.Kind, z: kindZ}]; ok {
		return nil
	}

	// More detailed checks can be added here, after basic validity
	// has been checked and before the "default deny" error.

	return fmt.Errorf("unsupport relationship between entites: '%v'", rel.String())
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
