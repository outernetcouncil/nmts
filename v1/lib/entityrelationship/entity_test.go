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
	"sync"
	"testing"

	"google.golang.org/protobuf/encoding/prototext"
	"outernetcouncil.org/nmts/v1/lib/entityrelationship"
	npb "outernetcouncil.org/nmts/v1/proto"
)

func mustUnmarshalEntity(t *testing.T, txtPb string) *npb.Entity {
	t.Helper()
	e := new(npb.Entity)
	if err := prototext.Unmarshal([]byte(txtPb), e); err != nil {
		t.Fatalf("failed to parse %s: %q", txtPb, err)
	}
	return e
}

var entityKindStringTestCases = []struct {
	txtPb string
	want  string
}{
	{txtPb: `ek_sdn_agent{}`, want: "EK_SDN_AGENT"},
	{txtPb: `ek_network_node{}`, want: "EK_NETWORK_NODE"},
	{txtPb: `ek_platform{}`, want: "EK_PLATFORM"},
	{txtPb: `ek_interface{}`, want: "EK_INTERFACE"},
	{txtPb: `ek_antenna{}`, want: "EK_ANTENNA"},
	{txtPb: `id: "no_kind"`, want: ""},
}

func TestEntityKindStringExamples(t *testing.T) {
	for _, tc := range entityKindStringTestCases {
		parsed := mustUnmarshalEntity(t, tc.txtPb)
		// Call twice: the result must be stable across repeated calls for the
		// same entity.
		for i := 0; i < 2; i++ {
			got := entityrelationship.EntityKindStringFromProto(parsed)
			if got != tc.want {
				t.Errorf("EntityKindStringFromProto(%s) call %d: wanted %q, got %q", tc.txtPb, i+1, tc.want, got)
			}
		}
	}
}

func TestEntityKindStringNilEntity(t *testing.T) {
	if got := entityrelationship.EntityKindStringFromProto(nil); got != "" {
		t.Errorf("EntityKindStringFromProto(nil): wanted %q, got %q", "", got)
	}
}

func TestEntityKindStringConcurrentCallers(t *testing.T) {
	parsed := make([]*npb.Entity, len(entityKindStringTestCases))
	for i, tc := range entityKindStringTestCases {
		parsed[i] = mustUnmarshalEntity(t, tc.txtPb)
	}
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i, tc := range entityKindStringTestCases {
				if got := entityrelationship.EntityKindStringFromProto(parsed[i]); got != tc.want {
					t.Errorf("EntityKindStringFromProto(%s): wanted %q, got %q", tc.txtPb, tc.want, got)
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkEntityKindStringFromProto(b *testing.B) {
	e := &npb.Entity{}
	if err := prototext.Unmarshal([]byte(`ek_network_node{}`), e); err != nil {
		b.Fatalf("failed to parse entity: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entityrelationship.EntityKindStringFromProto(e)
	}
}
