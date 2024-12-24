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
	"slices"
	"strings"

	// nosemgrep: import-text-template
	"text/template"

	"github.com/ichiban/prolog"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/prototext"
	er "outernetcouncil.org/nmts/v0/lib/entityrelationship"
	npb "outernetcouncil.org/nmts/v0/proto"
	eklpb "outernetcouncil.org/nmts/v0/proto/ek/logical"
)

func exportDot(appCtx *cli.Context) error {
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
	toLabel := func(key string) string {
		return toRecordLabels(g.Entities[key])
	}
	subgraphTmpl := template.Must(template.New("subgraphs").Funcs(template.FuncMap{"toLabel": toLabel}).Parse(`
subgraph cluster_{{.N}} {
	label = "{{.Parent.Name}}";
	{{ printf "%q [label=%q, class=%q]" .Parent.Name (toLabel .Parent.Name) .Parent.Class }}
	{{ range .Children }}
	{{- printf "%q [label=%q, class=%q]" .Name (toLabel .Name) .Class }}
	{{ end }}
}
`))
	type graphEntity struct {
		Name, Class string
	}
	type subgraphInput struct {
		Parent   graphEntity
		Children []graphEntity
		N        int
	}

	buf := &bytes.Buffer{}
	buf.WriteString("digraph G {\n")
	buf.WriteString("node [shape=record];\n")

	keys := make([]prologEntity, 0, len(roots))
	for parent := range roots {
		keys = append(keys, parent)
	}
	slices.SortFunc(keys, compareEntity)

	// first, iterate the roots and create subgraphs
	for _, parent := range keys {
		children := roots[parent]
		slices.SortFunc(children, compareEntity)
		seen[parent] = struct{}{}
		for _, c := range children {
			seen[c] = struct{}{}
		}

		sgParent := graphEntity{
			Name:  parent.Name,
			Class: containerClass,
		}
		sgChildren := []graphEntity{}
		for _, c := range children {
			class, err := toNodeClass(p, c.ID)
			if err != nil {
				return err
			}
			sgChildren = append(sgChildren, graphEntity{Name: c.Name, Class: class})
		}
		if err := subgraphTmpl.Execute(buf, subgraphInput{
			N:        len(seen),
			Parent:   sgParent,
			Children: sgChildren,
		}); err != nil {
			return fmt.Errorf("executing template for parent %q: %w", parent, err)
		}
	}
	ents, err := collectQuery[prologEntity](p, `is_entity(ID), display_name(ID, Name).`)
	if err != nil {
		return err
	}
	slices.SortFunc(ents, compareEntity)
	for _, e := range ents {
		if _, ok := seen[e]; ok {
			continue
		}

		seen[e] = struct{}{}
		class, err := toNodeClass(p, e.ID)
		if err != nil {
			return err
		}
		fmt.Fprintf(buf, "%q [label=%q, class=%q]\n", e.Name, toLabel(e.Name), class)
	}

	type edge struct{ AID, AName, ZID, ZName, K string }
	rs, err := collectQuery[edge](p, `connected(AID, ZID, K), display_name(AID, AName), display_name(ZID, ZName).`)
	if err != nil {
		return err
	}
	slices.SortFunc(rs, func(l, r edge) int {
		return cmp.Or(
			cmp.Compare(l.AID, r.AID),
			cmp.Compare(l.ZID, r.ZID),
			cmp.Compare(l.AName, r.AName),
			cmp.Compare(l.ZName, r.ZName),
		)
	})
	for _, r := range rs {
		fmt.Fprintf(buf, "%q -> %q [label=%q]\n", r.AName, r.ZName, r.K)
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
