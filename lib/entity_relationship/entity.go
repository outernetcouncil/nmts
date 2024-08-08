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
	"reflect"
	"strings"

	npb "outernetcouncil.org/nmts/proto"
)

func EntityKindStringFromProto(e *npb.Entity) string {
	if e == nil {
		return ""
	}

	kind := e.GetKind()
	if kind == nil {
		return ""
	}

	tag := reflect.TypeOf(kind).Elem().Field(0).Tag.Get("protobuf")
	for _, element := range strings.Split(tag, ",") {
		if strings.HasPrefix(element, "name=") {
			return strings.ToUpper(strings.TrimPrefix(element, "name="))
		}
	}

	return ""
}
