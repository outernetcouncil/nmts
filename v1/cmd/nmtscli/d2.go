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

	// nosemgrep: import-text-template
	"text/template"

	"github.com/ichiban/prolog"
	"github.com/urfave/cli/v2"
)

func exportD2(appCtx *cli.Context) error {
	g, err := readGraph(appCtx)
	if err != nil {
		return err
	}

	p := prolog.New(nil, nil)
	if err := addGraphToPrologInterpreter(p, g); err != nil {
		return fmt.Errorf("loading graph: %w", err)
	}
	roots, err := queryForRootContainers(p)
	if err != nil {
		return err
	}
	seen := map[string]struct{}{}
	subgraphTmpl := template.Must(template.New("subgraphs").Parse(`
{{.Parent}} {
	label: "{{.Parent}}";
	{{ printf "%q" .Parent }}
	{{ range .Children }}
	{{- printf "%q" . }}
	{{ end }}
}
`))
	type subgraphInput struct {
		Parent   string
		Children []string
		N        int
	}

	buf := &bytes.Buffer{}

	prefixes := map[string]string{}

	// first, iterate the roots and create subgraphs
	for parent, children := range roots {
		prefixes[parent] = parent
		seen[parent] = struct{}{}
		for _, c := range children {
			prefixes[c] = parent
			seen[c] = struct{}{}
		}
	}
	ents, err := collectQuery[prologEntity](p, `is_entity(EK, ID).`)
	if err != nil {
		return err
	}
	for _, e := range ents {
		if _, ok := seen[e.ID]; ok {
			continue
		}
		seen[e.ID] = struct{}{}
		fmt.Fprintf(buf, "%q\n", e.ID)
	}

	rs, err := collectQuery[struct{ A, Z, RK string }](p, `edge(A, Z, RK).`)
	for _, r := range rs {
		aName := fmt.Sprintf("%q", r.A)
		if aPrefix, ok := prefixes[r.A]; ok {
			aName = fmt.Sprintf("%s.%q", aPrefix, r.A)
		}
		zName := fmt.Sprintf("%q", r.Z)
		if zPrefix, ok := prefixes[r.Z]; ok {
			zName = fmt.Sprintf("%s.%q", zPrefix, r.Z)
		}
		fmt.Fprintf(buf, "%s -> %s: %q\n", aName, zName, r.RK)
	}

	// finally, iterate the roots and create subgraphs
	for parent, children := range roots {
		if err := subgraphTmpl.Execute(buf, subgraphInput{
			N:        len(seen),
			Parent:   parent,
			Children: children,
		}); err != nil {
			return fmt.Errorf("executing template for parent %q: %w", parent, err)
		}
	}
	_, err = io.Copy(appCtx.App.Writer, buf)
	return err
}
