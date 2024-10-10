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

package entityrelationship

import (
	"errors"

	npb "outernetcouncil.org/nmts/proto"
)

type Validator interface {
	// Validate each entity as it's loaded within the collection context
	// assembled up to that point.
	ValidateEntity(*Collection, *npb.Entity) error

	// Validate each relationship as it's loaded within the collection
	// context assembled up to that point.
	ValidateRelationship(*Collection, Relationship) error

	// Validate the complete collection.
	ValidateCollection(*Collection) error
}

type CollectionBuilder struct {
	erColl    *Collection
	validator Validator
}

func NewNonValidatingCollectionBuilder() *CollectionBuilder {
	return NewCollectionBuilder(nil)
}

func NewCollectionBuilder(v Validator) *CollectionBuilder {
	return &CollectionBuilder{
		erColl:    NewCollection(),
		validator: v,
	}
}

func (builder *CollectionBuilder) Build() (*Collection, error) {
	if builder.validator != nil {
		if err := builder.validator.ValidateCollection(builder.erColl); err != nil {
			return nil, err
		}
	}

	return builder.erColl, nil
}

func (builder *CollectionBuilder) InsertFragments(fragments ...*npb.Fragment) error {
	if fragments == nil {
		return nil
	}

	errs := []error{}

	for _, fragment := range fragments {
		if fragment.Entity != nil {
			for _, entity := range fragment.Entity {
				if builder.validator != nil {
					if err := builder.validator.ValidateEntity(builder.erColl, entity); err != nil {
						errs = append(errs, err)
						continue
					}
				}
				errs = append(errs, builder.erColl.InsertEntity(entity))
			}
		}
	}

	for _, fragment := range fragments {
		if fragment.Relationship != nil {
			for _, relationship := range fragment.Relationship {
				rel := RelationshipFromProto(relationship)
				if builder.validator != nil {
					if err := builder.validator.ValidateRelationship(builder.erColl, rel); err != nil {
						errs = append(errs, err)
						continue
					}
				}
				errs = append(errs, builder.erColl.InsertRelationship(rel))
			}
		}
	}

	return errors.Join(errs...)
}
