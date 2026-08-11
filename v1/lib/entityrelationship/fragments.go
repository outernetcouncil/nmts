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
	"os"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	npb "outernetcouncil.org/nmts/v1/proto"
)

func ReadFragmentFiles(fragmentFilenames []string) (*npb.Fragment, error) {
	g := &npb.Fragment{}

	for _, f := range fragmentFilenames {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("reading %q: %w", f, err)
		}
		subg := &npb.Fragment{}
		if err := prototext.Unmarshal(data, subg); err != nil {
			return nil, fmt.Errorf("parsing %q: %w", f, err)
		}

		for _, e := range subg.GetEntity() {
			g.Entity = append(g.Entity, e)
		}
		for _, r := range subg.GetRelationship() {
			g.Relationship = append(g.Relationship, r)
		}
	}

	return g, nil
}

type (
	BytesUpperBound        int
	ElementCountUpperBound int
)

// MaxFragmentUpsertBytes bounds the serialized size of each UpsertFragment
// request. It is set conservatively to stay well within typical transport and
// storage receive limits, leaving headroom for per-request overhead that
// SplitBySize does not measure and keeping per-transaction memory modest. For
// small elements the element-count cap, not this byte bound, is the binding
// constraint -- see MaxFragmentUpsertElements.
const MaxFragmentUpsertBytes = BytesUpperBound(9 * 1024 * 1024)

// MaxFragmentUpsertElements bounds the number of elements in each batch. It is
// set conservatively so that a batch stays within typical per-transaction limits
// on the number of underlying writes, leaving headroom for per-element overhead
// and future growth.
const MaxFragmentUpsertElements = ElementCountUpperBound(1000)

// SplitBySize splits f into a sequence of fragments that each stay under BOTH
// bounds: maxBytes (serialized size, to stay within transport receive limits and
// keep per-transaction memory modest) and maxElements (element count per batch,
// to stay within per-transaction write limits). All entities are emitted
// first, then all relationships, so every relationship batch can be applied after
// the entities it references already exist. A single element over a bound is
// placed in its own batch. An empty fragment yields no batches.
func SplitBySize(f *npb.Fragment, maxBytes BytesUpperBound, maxElements ElementCountUpperBound) []*npb.Fragment {
	var batches []*npb.Fragment
	var cur *npb.Fragment
	var curBytes, curCount int

	flush := func() {
		if cur != nil {
			batches = append(batches, cur)
			cur = nil
			curBytes = 0
			curCount = 0
		}
	}

	for _, e := range f.GetEntity() {
		sz := proto.Size(e)
		if cur != nil && (curBytes+sz > int(maxBytes) || curCount+1 > int(maxElements)) {
			flush()
		}
		if cur == nil {
			cur = &npb.Fragment{}
		}
		cur.Entity = append(cur.Entity, e)
		curBytes += sz
		curCount++
	}
	// Flush before relationships so entities and relationships never share a
	// batch: every relationship batch is applied strictly after all entities.
	flush()

	for _, r := range f.GetRelationship() {
		sz := proto.Size(r)
		if cur != nil && (curBytes+sz > int(maxBytes) || curCount+1 > int(maxElements)) {
			flush()
		}
		if cur == nil {
			cur = &npb.Fragment{}
		}
		cur.Relationship = append(cur.Relationship, r)
		curBytes += sz
		curCount++
	}
	flush()

	return batches
}
