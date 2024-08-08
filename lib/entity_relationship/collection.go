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

package entity_relationship

import (
	"errors"
	"fmt"

	npb "outernetcouncil.org/nmts/proto"
)

type Collection struct {
	Entities map[string]*npb.Entity
	OutEdges map[string]*RelationshipSet
	InEdges  map[string]*RelationshipSet
}

func NewCollection() *Collection {
	return &Collection{
		Entities: make(map[string]*npb.Entity),
		OutEdges: make(map[string]*RelationshipSet),
		InEdges:  make(map[string]*RelationshipSet),
	}
}

func (erColl *Collection) NumEntities() int {
	return len(erColl.Entities)
}

func (erColl *Collection) NumRelationships() int {
	sum := 0
	for _, relSet := range erColl.OutEdges {
		sum += len(relSet.Relations)
	}
	return sum
}

func (erColl *Collection) EntityExists(key string) bool {
	_, exists := erColl.Entities[key]
	return exists
}

func (erColl *Collection) InsertEntity(entity *npb.Entity) error {
	if entity != nil {
		key := entity.Id
		if erColl.EntityExists(key) {
			return fmt.Errorf("entity already exists: '%v'", key)
		}
		erColl.Entities[key] = entity
	}

	return nil
}

func (erColl *Collection) InsertRelationshipProto(relationship *npb.Relationship) error {
	return erColl.InsertRelationship(RelationshipFromProto(relationship))
}

func (erColl *Collection) InsertRelationship(r Relationship) error {
	errs := []error{}

	if !erColl.EntityExists(r.A) {
		errs = append(errs, fmt.Errorf("relationship references non-existent entity: '%v'", r.A))
	}
	if !erColl.EntityExists(r.Z) {
		errs = append(errs, fmt.Errorf("relationship references non-existent entity: '%v'", r.Z))
	}

	addRelationship := func(k string, m map[string]*RelationshipSet) error {
		rs, exists := m[k]
		if !exists {
			rs = NewRelationshipSet()
			m[k] = rs
		}
		return rs.Insert(r)
	}

	if len(errs) == 0 {
		errs = append(errs,
			addRelationship(r.A, erColl.OutEdges),
			addRelationship(r.Z, erColl.InEdges),
		)
	}

	return errors.Join(errs...)
}
