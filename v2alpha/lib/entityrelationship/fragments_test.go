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

package entityrelationship_test

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"outernetcouncil.org/nmts/v2alpha/lib/entityrelationship"
	npb "outernetcouncil.org/nmts/v2alpha/proto"
)

func TestSplitBySize_EmptyFragment(t *testing.T) {
	if got := entityrelationship.SplitBySize(&npb.Fragment{}, 1024, 10); len(got) != 0 {
		t.Errorf("SplitBySize(empty) = %d batches, want 0", len(got))
	}
}

func TestSplitBySize_EntitiesBeforeRelationships(t *testing.T) {
	f := &npb.Fragment{
		Entity:       []*npb.Entity{{Id: "e1"}, {Id: "e2"}},
		Relationship: []*npb.Relationship{{A: "e1", Z: "e2", Kind: npb.RK_RK_CONTAINS}},
	}
	// Bounds large enough that everything fits, but entities and relationships
	// must never share a batch.
	got := entityrelationship.SplitBySize(f, 1<<20, 1000)
	if len(got) != 2 {
		t.Fatalf("SplitBySize = %d batches, want 2 (entities batch, then relationships batch)", len(got))
	}
	if len(got[0].GetEntity()) != 2 || len(got[0].GetRelationship()) != 0 {
		t.Errorf("batch 0 = %d entities / %d rels, want 2 / 0", len(got[0].GetEntity()), len(got[0].GetRelationship()))
	}
	if len(got[1].GetEntity()) != 0 || len(got[1].GetRelationship()) != 1 {
		t.Errorf("batch 1 = %d entities / %d rels, want 0 / 1", len(got[1].GetEntity()), len(got[1].GetRelationship()))
	}
}

func TestSplitBySize_ElementBoundSplitsBeforeBytes(t *testing.T) {
	f := &npb.Fragment{Entity: []*npb.Entity{
		{Id: "e1"}, {Id: "e2"}, {Id: "e3"}, {Id: "e4"}, {Id: "e5"},
	}}
	// Bytes are effectively unbounded; the element bound (2) is binding.
	got := entityrelationship.SplitBySize(f, 1<<20, 2)
	if len(got) != 3 {
		t.Fatalf("SplitBySize = %d batches, want 3 (2+2+1)", len(got))
	}
	for i, want := range []int{2, 2, 1} {
		if got := len(got[i].GetEntity()); got != want {
			t.Errorf("batch %d = %d entities, want %d", i, got, want)
		}
	}
}

func TestSplitBySize_ElementBoundSplitsRelationships(t *testing.T) {
	f := &npb.Fragment{Relationship: []*npb.Relationship{
		{A: "a", Z: "b", Kind: npb.RK_RK_CONTAINS},
		{A: "a", Z: "c", Kind: npb.RK_RK_CONTAINS},
		{A: "a", Z: "d", Kind: npb.RK_RK_CONTAINS},
	}}
	got := entityrelationship.SplitBySize(f, 1<<20, 2)
	if len(got) != 2 {
		t.Fatalf("SplitBySize = %d batches, want 2 (2+1)", len(got))
	}
	for i, want := range []int{2, 1} {
		if n := len(got[i].GetRelationship()); n != want {
			t.Errorf("batch %d = %d relationships, want %d", i, n, want)
		}
	}
}

func TestSplitBySize_ByteBound(t *testing.T) {
	// maxBytes == the serialized size of exactly one entity, so each batch holds
	// exactly one entity while the element bound is not binding.
	one := entityrelationship.BytesUpperBound(proto.Size(&npb.Entity{Id: "entity-with-id"}))
	f := &npb.Fragment{Entity: []*npb.Entity{
		{Id: "entity-with-id"}, {Id: "entity-with-id"}, {Id: "entity-with-id"},
	}}
	got := entityrelationship.SplitBySize(f, one, 1000)
	if len(got) != 3 {
		t.Fatalf("SplitBySize = %d batches, want 3", len(got))
	}
	for i := range got {
		if n := len(got[i].GetEntity()); n != 1 {
			t.Errorf("batch %d = %d entities, want 1", i, n)
		}
	}
}

func TestSplitBySize_SingleElementOverBoundGetsOwnBatch(t *testing.T) {
	f := &npb.Fragment{Entity: []*npb.Entity{{Id: "big"}}}
	got := entityrelationship.SplitBySize(f, 1, 1000) // maxBytes smaller than any element
	if len(got) != 1 || len(got[0].GetEntity()) != 1 {
		t.Fatalf("SplitBySize = %v, want a single 1-entity batch", got)
	}
}
