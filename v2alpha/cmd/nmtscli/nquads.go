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
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/urfave/cli/v2"
	er "outernetcouncil.org/nmts/v2alpha/lib/entityrelationship"
	// npb "outernetcouncil.org/nmts/v2alpha/proto"
)

// TODO: choose sensible values for this
const (
	nquadSchema     = "http"
	nquadHost       = "example.org"
	nquadSchemaPath = "nmts/schema"
	nquadEntityPath = "nmts/entities"
)

// https://www.w3.org/TR/n-quads/
func exportNQuads(appCtx *cli.Context) error {
	g, err := readGraph(appCtx)
	if err != nil {
		return err
	}

	w := &bytes.Buffer{}
	for _, edges := range g.OutEdges {
		for edge := range edges.Relations {
			fmt.Fprintf(w, "<%s> <%s> <%s>.\n",
				keyToIRI(edge.A), edgeToIRI(edge), keyToIRI(edge.Z),
			)
		}
	}

	_, err = io.Copy(appCtx.App.Writer, w)
	return err
}

func keyToIRI(k string) string {
	return (&url.URL{
		Scheme:   nquadSchema,
		Host:     nquadHost,
		Path:     nquadEntityPath,
		Fragment: k,
	}).String()
}

func edgeToIRI(k er.Relationship) string {
	return (&url.URL{
		Scheme: nquadSchema,
		Host:   nquadHost,
		Path:   nquadSchemaPath,
	}).JoinPath(strings.ToLower(strings.TrimPrefix(k.Kind.String(), "RK_"))).String()
}
