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
	rules := graphToPrologFacts(g)

	p := prolog.New(nil, nil)
	if err = p.Exec(rules); err != nil {
		return err
	}

	roots, err := queryForRootContainers(p)
	if err != nil {
		return err
	}
	seen := map[prologEntity]struct{}{}
	subgraphTmpl := template.Must(template.New("subgraphs").Parse(`
{{.Parent.ID}} {
	label: "{{.Parent.Name}}";
	{{ printf "%q" .Parent.Name }}
	{{ range .Children }}
	{{- printf "%q" .Name }}
	{{ end }}
}
`))
	type subgraphInput struct {
		Parent   prologEntity
		Children []prologEntity
		N        int
	}

	buf := &bytes.Buffer{}

	prefixes := map[string]string{}

	// first, iterate the roots and create subgraphs
	for parent, children := range roots {
		prefixes[parent.ID] = parent.ID
		seen[parent] = struct{}{}
		for _, c := range children {
			prefixes[c.ID] = parent.ID
			seen[c] = struct{}{}
		}
	}
	ents, err := collectQuery[prologEntity](p, `is_entity(ID), display_name(ID, Name).`)
	if err != nil {
		return err
	}
	for _, e := range ents {
		if _, ok := seen[e]; ok {
			continue
		}
		seen[e] = struct{}{}
		fmt.Fprintf(buf, "%q\n", e.Name)
	}

	rs, err := collectQuery[struct {
		AID, AName, ZID, ZName, K string
	}](p, `connected(AID, ZID, K), display_name(AID, AName), display_name(ZID, ZName).`)
	for _, r := range rs {
		aName := fmt.Sprintf("%q", r.AName)
		if aPrefix, ok := prefixes[r.AID]; ok {
			aName = fmt.Sprintf("%s.%q", aPrefix, r.AName)
		}
		zName := fmt.Sprintf("%q", r.ZName)
		if zPrefix, ok := prefixes[r.ZID]; ok {
			zName = fmt.Sprintf("%s.%q", zPrefix, r.ZName)
		}
		fmt.Fprintf(buf, "%s -> %s: %q\n", aName, zName, r.K)
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
