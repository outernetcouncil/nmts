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
	"errors"
	"fmt"
	"strings"

	"golang.org/x/text/unicode/norm"
	"google.golang.org/protobuf/reflect/protoreflect"
	er "outernetcouncil.org/nmts/v2alpha/lib/entityrelationship"
	"outernetcouncil.org/nmts/v2alpha/lib/graph"
	npb "outernetcouncil.org/nmts/v2alpha/proto"
)

func IsEntityMinimallyWellFormed(entity *npb.Entity) error {
	if entity == nil {
		return fmt.Errorf("entity MUST NOT be nil")
	}
	id := entity.Id

	if err := IsEntityIDMinimallyWellFormed(id); err != nil {
		return err
	}

	if er.EntityKindStringFromProto(entity) == "" {
		return fmt.Errorf("Entity %q lacks an entity kind field", id)
	}

	return nil
}

// IsEntityIDMinimallyWellFormed checks that an Entity ID is non-empty, in
// Unicode Normalization Form C, and free of leading and trailing whitespace.
func IsEntityIDMinimallyWellFormed(id string) error {
	// TODO: uuid.Validate(id)
	if id == "" {
		return fmt.Errorf("id must not be empty: %q", id)
	}

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

	return nil
}

// IsRelationshipMinimallyWellFormed checks the properties a Relationship can
// be held to on its own: both endpoints are minimally well formed and are not
// the same entity, and a relationship kind is set.
// Consistency with a surrounding model -- whether the endpoints refer to
// Entities that exist, and whether the kind is permitted between those two
// Entities' kinds -- requires a Collection or Graph and is checked elsewhere.
func IsRelationshipMinimallyWellFormed(relationship *npb.Relationship) error {
	if relationship == nil {
		return fmt.Errorf("relationship MUST NOT be nil")
	}
	if err := IsEntityIDMinimallyWellFormed(relationship.A); err != nil {
		return fmt.Errorf("relationship A endpoint: %w", err)
	}
	if err := IsEntityIDMinimallyWellFormed(relationship.Z); err != nil {
		return fmt.Errorf("relationship Z endpoint: %w", err)
	}
	if relationship.A == relationship.Z {
		return fmt.Errorf("relationship cannot be self-referential")
	}
	if relationship.Kind == npb.RK_RK_UNSPECIFIED {
		return fmt.Errorf("relationship between %q and %q lacks a relationship kind", relationship.A, relationship.Z)
	}
	return nil
}

// ValidateFragment checks the properties a Fragment can be held to on its own:
// every Entity and every Relationship is minimally well formed, no entity ID
// appears more than once, and no (a, kind, z) tuple appears more than once.
//
// A Fragment is a subgraph, so the Entities a Relationship refers to need not be
// carried within it. Endpoint existence and permitted entity-kind pairings
// therefore cannot be resolved here; they require the surrounding Collection or
// Graph. An empty Fragment is valid: whether one is meaningful is a matter for
// the caller.
//
// Validation does not stop at the first problem; every problem found is
// reported.
func ValidateFragment(fragment *npb.Fragment) error {
	var errs []error
	seenEntityIDs := make(map[string]struct{}, len(fragment.GetEntity()))
	for _, entity := range fragment.GetEntity() {
		if err := IsEntityMinimallyWellFormed(entity); err != nil {
			errs = append(errs, fmt.Errorf("invalid Entity: %w", err))
		}
		if _, ok := seenEntityIDs[entity.GetId()]; ok {
			errs = append(errs, fmt.Errorf("duplicate entity id %q", entity.GetId()))
		}
		seenEntityIDs[entity.GetId()] = struct{}{}
	}

	seenRelationships := make(map[er.Relationship]struct{}, len(fragment.GetRelationship()))
	for _, relationship := range fragment.GetRelationship() {
		if err := IsRelationshipMinimallyWellFormed(relationship); err != nil {
			errs = append(errs, fmt.Errorf("invalid Relationship: %w", err))
		}
		rel := er.RelationshipFromProto(relationship)
		if _, ok := seenRelationships[rel]; ok {
			errs = append(errs, fmt.Errorf("duplicate relationship %q", rel.String()))
		}
		seenRelationships[rel] = struct{}{}
	}

	return errors.Join(errs...)
}

type DefaultValidator struct{}

// Validate each entity as it's loaded within the collection context
// assembled up to that point.
func (DefaultValidator) ValidateEntity(coll *er.Collection, entity *npb.Entity) error {
	return IsEntityMinimallyWellFormed(entity)
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

// subKindOneofs names, for each Entity kind that has one, the oneof in the
// kind-specific message whose populated field gives the Entity's sub-kind:
// the layer of an EK_INTERFACE or the payload of an EK_LOGICAL_PACKET_LINK.
var subKindOneofs = map[string]protoreflect.Name{
	"EK_INTERFACE":           "layer",
	"EK_LOGICAL_PACKET_LINK": "payload",
}

// EntitySubKind returns the sub-kind of an Entity: the name of the field set
// in the oneof listed for its kind in subKindOneofs, e.g. "mpls" for an
// EK_INTERFACE with `mpls {}` or an EK_LOGICAL_PACKET_LINK with `mpls {}`.
//
// It returns "" if the Entity's kind has no sub-kinds or if none is set.
func EntitySubKind(entity *npb.Entity) string {
	oneofName, ok := subKindOneofs[er.EntityKindStringFromProto(entity)]
	if !ok {
		return ""
	}
	m := entity.ProtoReflect()
	kindField := m.WhichOneof(m.Descriptor().Oneofs().ByName("kind"))
	if kindField == nil || kindField.Message() == nil {
		return ""
	}
	kindMsg := m.Get(kindField).Message()
	oneof := kindMsg.Descriptor().Oneofs().ByName(oneofName)
	if oneof == nil {
		return ""
	}
	if set := kindMsg.WhichOneof(oneof); set != nil {
		return string(set.Name())
	}
	return ""
}

// sameSubKindRelationships lists relationships whose endpoints must share a
// sub-kind for the sub-kinds given: if either endpoint has one of the listed
// sub-kinds, the other endpoint must have the same one.
//
// For example, an EK_LOGICAL_PACKET_LINK at the MPLS layer can only be
// originated or terminated by an EK_INTERFACE also at the MPLS layer, and an
// MPLS-layer EK_INTERFACE can only originate or terminate MPLS links.
//
// An endpoint whose sub-kind is not set is not checked, so partially
// specified models keep validating.
var sameSubKindRelationships = map[allowedRelationship][]string{
	{a: "EK_INTERFACE", rk: npb.RK_RK_ORIGINATES, z: "EK_LOGICAL_PACKET_LINK"}: {"mpls"},
	{a: "EK_INTERFACE", rk: npb.RK_RK_TERMINATES, z: "EK_LOGICAL_PACKET_LINK"}: {"mpls"},
}

// checkSubKinds applies sameSubKindRelationships to a relationship whose
// entity kinds are already known to be permitted.
func checkSubKinds(key allowedRelationship, rel er.Relationship, entityA, entityZ *npb.Entity) error {
	subKinds, ok := sameSubKindRelationships[key]
	if !ok {
		return nil
	}
	subA := EntitySubKind(entityA)
	subZ := EntitySubKind(entityZ)
	if subA == "" || subZ == "" || subA == subZ {
		return nil
	}
	for _, sub := range subKinds {
		if subA == sub || subZ == sub {
			return fmt.Errorf(
				"unsupported relationship between entities: '%v': %s %q (%s) and %s %q (%s) must both be %q",
				rel.String(), key.a, rel.A, subA, key.z, rel.Z, subZ, sub)
		}
	}
	return nil
}

// Validate each relationship as it's loaded within the collection
// context assembled up to that point.
func (DefaultValidator) ValidateRelationship(coll *er.Collection, rel er.Relationship) error {
	entityA := coll.Entities[rel.A]
	entityZ := coll.Entities[rel.Z]
	kindA := er.EntityKindStringFromProto(entityA)
	kindZ := er.EntityKindStringFromProto(entityZ)

	key := allowedRelationship{a: kindA, rk: rel.Kind, z: kindZ}
	if _, ok := permittedRelationships[key]; ok {
		// More detailed checks, after basic validity has been checked.
		return checkSubKinds(key, rel, entityA, entityZ)
	}

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
		// More detailed checks, after basic validity has been checked.
		return checkSubKinds(key, rel, a.GetEntity(), z.GetEntity())
	}

	return fmt.Errorf("unsupported relationship between entites: '%v' i.e. '%v'", rel.String(), key)
}
