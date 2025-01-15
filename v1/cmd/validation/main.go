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

package main

import (
	"log"

	protovalidate "github.com/bufbuild/protovalidate-go"
	"outernetcouncil.org/nmts/v1/proto/types/ietf"
)

func main() {
	inetProtoInvalid := &ietf.IPv4Address{
		Str: "Test",
	}
	inetProtoValid := &ietf.IPv4Address{
		Str: "127.0.0.1",
	}
	if err := protovalidate.Validate(inetProtoInvalid); err != nil {
		log.Println("validation failed:", err)
	} else {
		log.Println("validation succeeded")
	}

	if err := protovalidate.Validate(inetProtoValid); err != nil {
		log.Println("validation failed:", err)
	} else {
		log.Println("validation succeeded")
	}

	return
}
