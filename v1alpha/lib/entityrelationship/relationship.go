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
	"fmt"

	npb "outernetcouncil.org/nmts/v1alpha/proto"
)

type Relationship struct {
	A    string
	Z    string
	Kind npb.RK
}

func RelationshipFromProto(r *npb.Relationship) Relationship {
	return Relationship{
		A:    r.A,
		Z:    r.Z,
		Kind: r.Kind,
	}
}

func (r *Relationship) ToProto() *npb.Relationship {
	return &npb.Relationship{
		A:    r.A,
		Z:    r.Z,
		Kind: r.Kind,
	}
}

func (r *Relationship) String() string {
	return fmt.Sprintf("%v->%s->%v", r.A, r.Kind.String(), r.Z)
}

type RelationshipSet struct {
	Relations map[Relationship]struct{}
}

func NewRelationshipSet() *RelationshipSet {
	return &RelationshipSet{
		Relations: make(map[Relationship]struct{}),
	}
}

func (rset *RelationshipSet) Insert(r Relationship) error {
	if _, exists := rset.Relations[r]; exists {
		return fmt.Errorf("relationship already exists in set: '%v'", r)
	}
	rset.Relations[r] = struct{}{}
	return nil
}
