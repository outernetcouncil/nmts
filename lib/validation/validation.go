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

	er "outernetcouncil.org/nmts/lib/entity_relationship"
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
	if "" == id {
		return fmt.Errorf("id must not be empty: '%s'", id)
	}

	if "" == er.EntityKindStringFromProto(entity) {
		return fmt.Errorf("Entity '%s' lacks an entity kind field", id)
	}

	return nil
}

func makeUnsupportedRelationshipError(rel er.Relationship) error {
	return fmt.Errorf("unsupport relationship between entites: '%v'", rel.String())
}

// Validate each relationship as it's loaded within the collection
// context assembled up to that point.
func (DefaultValidator) ValidateRelationship(coll *er.Collection, rel er.Relationship) error {
	switch rel.Kind {
	case npb.RK_RK_CONTAINS:
		return ValidateRkContains(coll, rel)
	case npb.RK_RK_CONTROLS:
		return ValidateRkControls(coll, rel)
	case npb.RK_RK_ORIGINATES:
		return ValidateRkOriginates(coll, rel)
	case npb.RK_RK_TERMINATES:
		return ValidateRkTerminates(coll, rel)
	case npb.RK_RK_TRAVERSES:
		return ValidateRkTraverses(coll, rel)
	default:
		// By default, Relationhip Kinds need to allowlist the entities
		// that are accepted. This may be by simple Entity Kind checks,
		// or by other contraints (e.g., in the graph structure).
		//
		// For example, an EK_PLATFORM may RK_CONTAINS an EK_NETWORK_NODE
		// but never the reverse.
		//
		// A more complex example might be a check that an EK_INTERFACE
		// RK_TRAVERSES only EK_PORTs that are RK_CONTAINS'd within the
		// same EK_PLATFORM (no traversing ports on another device).
		return makeUnsupportedRelationshipError(rel)
	}
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

type entityKindPair struct {
	a string
	z string
}

func validateSupportedRelationship(coll *er.Collection, rel er.Relationship, permitted []entityKindPair) error {
	for _, entry := range permitted {
		kindA := er.EntityKindStringFromProto(coll.Entities[rel.A])
		kindZ := er.EntityKindStringFromProto(coll.Entities[rel.Z])
		if entry.a == kindA && entry.z == kindZ {
			return nil
		}
	}

	return makeUnsupportedRelationshipError(rel)
}

func ValidateRkContains(coll *er.Collection, rel er.Relationship) error {
	permittedByEntityType := []entityKindPair{
		{"EK_PLATFORM", "EK_PLATFORM"},
		{"EK_PLATFORM", "EK_PORT"},
		{"EK_PLATFORM", "EK_NETWORK_NODE"},
		{"EK_NETWORK_NODE", "EK_INTERFACE"},
		{"EK_NETWORK_NODE", "EK_ROUTE_FN"},
		{"EK_NETWORK_NODE", "EK_SWITCH_FN"},
		{"EK_NETWORK_NODE", "EK_SDN_AGENT"},
	}

	return validateSupportedRelationship(coll, rel, permittedByEntityType)
}

func ValidateRkControls(coll *er.Collection, rel er.Relationship) error {
	permittedByEntityType := []entityKindPair{
		{"EK_SDN_AGENT", "EK_NETWORK_NODE"},
		{"EK_SDN_AGENT", "EK_PLATFORM"},
	}

	return validateSupportedRelationship(coll, rel, permittedByEntityType)
}

func ValidateRkOriginates(coll *er.Collection, rel er.Relationship) error {
	permittedByEntityType := []entityKindPair{
		{"EK_PORT", "EK_PHYSICAL_MEDIUM_LINK"},
		{"EK_INTERFACE", "EK_LOGICAL_PACKET_LINK"},
	}

	return validateSupportedRelationship(coll, rel, permittedByEntityType)
}

func ValidateRkTerminates(coll *er.Collection, rel er.Relationship) error {
	permittedByEntityType := []entityKindPair{
		{"EK_PORT", "EK_PHYSICAL_MEDIUM_LINK"},
		{"EK_INTERFACE", "EK_LOGICAL_PACKET_LINK"},
	}

	return validateSupportedRelationship(coll, rel, permittedByEntityType)
}

func ValidateRkTraverses(coll *er.Collection, rel er.Relationship) error {
	// TODO: validate same RK_CONTAINS EK_PLATFORM.
	permittedByEntityType := []entityKindPair{
		{"EK_INTERFACE", "EK_PORT"},
		{"EK_LOGICAL_PACKET_LINK", "EK_PHYSICAL_MEDIUM_LINK"},
	}

	return validateSupportedRelationship(coll, rel, permittedByEntityType)
}
