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

package testing

import (
	"errors"

	"google.golang.org/protobuf/encoding/prototext"
	er "outernetcouncil.org/nmts/v1/lib/entityrelationship"
	"outernetcouncil.org/nmts/v1/lib/graph"
	"outernetcouncil.org/nmts/v1/lib/validation"
	nmtspb "outernetcouncil.org/nmts/v1/proto"
)

// Consider wrapping in https://github.com/samber/lo lo.Must().
func EntityFrom(txtpb string) (*nmtspb.Entity, error) {
	e := &nmtspb.Entity{}
	if err := prototext.Unmarshal([]byte(txtpb), e); err != nil {
		return nil, err
	}
	return e, nil
}

// Consider wrapping in https://github.com/samber/lo lo.Must().
func FragmentFrom(txtpb string) (*nmtspb.Fragment, error) {
	e := &nmtspb.Fragment{}
	if err := prototext.Unmarshal([]byte(txtpb), e); err != nil {
		return nil, err
	}
	return e, nil
}

// Consider wrapping in https://github.com/samber/lo lo.Must().
func GraphFromFragments(fragments ...*nmtspb.Fragment) (*graph.Graph, error) {
	return UpdateGraphWithFragments(graph.New(), fragments...)
}

// Consider wrapping in https://github.com/samber/lo lo.Must().
func UpdateGraphWithFragments(g *graph.Graph, fragments ...*nmtspb.Fragment) (*graph.Graph, error) {
	errs := []error{}

	// Add entities before relationships, because relationship
	// insertion tests for a: and z: existence.
	for _, fragment := range fragments {
		for _, entity := range fragment.Entity {
			if _, err := g.UpsertEntity(entity); err != nil {
				errs = append(errs, err)
			}
		}
	}

	validator := validation.DefaultGraphValidator{}
	for _, fragment := range fragments {
		for _, relationship := range fragment.Relationship {
			if err := validator.ValidateRelationship(g, er.RelationshipFromProto(relationship)); err != nil {
				errs = append(errs, err)
				continue
			}
			if _, err := g.AddRelationship(relationship); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return g, errors.Join(errs...)
}
