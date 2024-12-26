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
	npb "outernetcouncil.org/nmts/v1alpha/proto"
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
