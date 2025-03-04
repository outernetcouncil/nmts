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
	"cmp"
	"fmt"
	"html"
	"io"
	"maps"
	"slices"
	"strings"

	// nosemgrep: import-text-template
	"text/template"

	"github.com/ichiban/prolog"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/prototext"
	er "outernetcouncil.org/nmts/v1/lib/entityrelationship"
	npb "outernetcouncil.org/nmts/v1/proto"
	eklpb "outernetcouncil.org/nmts/v1/proto/ek/logical"
)

func exportDot(appCtx *cli.Context) error {
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
	toLabel := func(key string) string { return toRecordLabels(g.Entities[key]) }
	subgraphTmpl := template.Must(template.New("subgraphs").Funcs(template.FuncMap{"toLabel": toLabel}).Parse(`
	subgraph cluster_{{.SubgraphID}} {
		label = "{{.Parent.ID}}";
		{{ printf "%q [label=%q, class=%q]" .Parent.ID (toLabel .Parent.ID) .Parent.Class }}
		{{ range .Children }}
		{{- printf "%q [label=%q, class=%q]" .ID (toLabel .ID) .Class }}
		{{ end }}
	}
`))
	type graphEntity struct {
		ID, Class string
	}
	type subgraphInput struct {
		Parent     graphEntity
		Children   []graphEntity
		SubgraphID int
	}

	buf := &bytes.Buffer{}
	buf.WriteString("digraph G {\n")
	buf.WriteString("\tnode [shape=record];\n")

	keys := slices.Sorted(maps.Keys(roots))

	// first, iterate the roots and create subgraphs
	for idx, parent := range keys {
		children := roots[parent]
		slices.Sort(children)

		seen[parent] = struct{}{}
		for _, c := range children {
			seen[c] = struct{}{}
		}

		sgParent := graphEntity{ID: parent, Class: containerClass}
		sgChildren := []graphEntity{}
		for _, c := range children {
			class, err := toNodeClass(p, c)
			if err != nil {
				return err
			}
			sgChildren = append(sgChildren, graphEntity{ID: c, Class: class})
		}
		if err := subgraphTmpl.Execute(buf, subgraphInput{
			SubgraphID: idx,
			Parent:     sgParent,
			Children:   sgChildren,
		}); err != nil {
			return fmt.Errorf("executing template for parent %q: %w", parent, err)
		}
	}
	remaining, err := collectQuery[prologEntity](p, `is_entity(EK, ID).`)
	if err != nil {
		return err
	}
	slices.SortFunc(remaining, compareEntity)
	for _, e := range remaining {
		if _, ok := seen[e.ID]; ok {
			continue
		}

		seen[e.ID] = struct{}{}
		class, err := toNodeClass(p, e.ID)
		if err != nil {
			return err
		}
		fmt.Fprintf(buf, "\t%q [label=%q, class=%q]\n", e.ID, toLabel(e.ID), class)
	}

	type edge struct{ A, Z, RK string }
	rs, err := collectQuery[edge](p, `edge(A, Z, RK).`)
	if err != nil {
		return err
	}
	slices.SortFunc(rs, func(l, r edge) int {
		return cmp.Or(cmp.Compare(l.A, r.A), cmp.Compare(l.Z, r.Z), cmp.Compare(l.RK, r.RK))
	})
	for _, r := range rs {
		fmt.Fprintf(buf, "\t%q -> %q [label=%q]\n", r.A, r.Z, r.RK)
	}

	buf.WriteString("}\n")
	_, err = io.Copy(appCtx.App.Writer, buf)
	return err
}

// Nodes of shape record can have multiple fields, which are provided as a
// "label" in the form of a sequence of pipe separated "port" / value pairs
// where the "port" is wrapped in angle brackets and the value is HTML escaped.
// The port is not rendered in the graph, but it can be used in conjunction
// with a node name to indicate where to attach a given edge to.
// https://graphviz.org/doc/info/shapes.html
func toRecordLabels(e *npb.Entity) string {
	fields := []string{
		fmt.Sprintf("<name> name = %s", htmlEscape(e.Id)),
		fmt.Sprintf("<kind> kind = %s", htmlEscape(er.EntityKindStringFromProto(e))),
	}

	switch v := e.GetKind().(type) {
	case *npb.Entity_EkInterface:
		switch layer := v.EkInterface.Layer.(type) {
		case *eklpb.Interface_Eth:
			fields = append(fields, fmt.Sprintf("<mac_addr> mac_addr = %s",
				htmlEscape(layer.Eth.GetMacAddr().GetStr())))
		case *eklpb.Interface_Cvlan:
		case *eklpb.Interface_Ip:
		}
	case *npb.Entity_EkPlatform:
		motionLabel := htmlEscape(prototext.MarshalOptions{}.Format(v.EkPlatform.GetMotion()))

		fields = append(fields, fmt.Sprintf("<motion> motion = %s", motionLabel))
	}

	return "{" + strings.Join(fields, "|") + "}"
}

func htmlEscape(s string) string {
	s = html.EscapeString(s)
	s = strings.ReplaceAll(s, "{", "&#123;")
	s = strings.ReplaceAll(s, "}", "&#125;")
	return s
}
