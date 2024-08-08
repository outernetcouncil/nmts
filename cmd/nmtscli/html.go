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
	"html/template"
	"io"
	"slices"

	"github.com/urfave/cli/v2"
	er "outernetcouncil.org/nmts/lib/entity_relationship"
	npb "outernetcouncil.org/nmts/proto"
)

const orbHtmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>Orb | Render a simple graph</title>
  <script src="https://unpkg.com/@memgraph/orb@0.4.3/dist/browser/orb.js"></script>
  <style>
    html, body {
      height: 100%;
      margin: 0;
    }
    #graph {
      height: 100%;
    }
  </style>
</head>
<body>
  <div id="graph"></div>
  <script>
    const container = document.getElementById('graph');

    const nodes = [
{{ range .Nodes }}
      { {{ toNodeString . }} },
{{ end }}
    ];
    const edges = [
{{ range .Edges }}
      { {{ toEdgeString . }} },
{{ end }}
    ];

    const orb = new Orb.Orb(container);

    // Initialize nodes and edges
    orb.data.setup({ nodes, edges });

	orb.view.setSettings({
		simulation: {
		  isPhysicsEnabled: true,
		},
	});

    // Render and recenter the view
    orb.view.render(() => {
      orb.view.recenter();
    });
  </script>
</body>
</html>
`

func orbNodeString(e *npb.Entity) template.JS {
	// nosemgrep: unescaped-data-in-js,unsafe-template-type
	return template.JS(fmt.Sprintf("id: %q, label: %q", e.Id, e.Id))
}

func orbEdgeString(r er.Relationship) template.JS {
	makeId := func(r er.Relationship) string {
		return fmt.Sprintf("%s__%s__%s", r.A, r.Kind.String(), r.Z)
	}
	// nosemgrep: unescaped-data-in-js,unsafe-template-type
	return template.JS(fmt.Sprintf("id: %q, start: %q, end: %q, label: %q", makeId(r), r.A, r.Z, r.Kind.String()))
}

func exportHtml(appCtx *cli.Context) error {
	g, err := readGraph(appCtx)
	if err != nil {
		return err
	}

	outputTmpl := template.Must(template.New("html").Funcs(template.FuncMap{"toNodeString": orbNodeString, "toEdgeString": orbEdgeString}).Parse(orbHtmlTemplate))

	type templateInput struct {
		Nodes []*npb.Entity
		Edges []er.Relationship
	}
	input := &templateInput{
		Nodes: []*npb.Entity{},
		Edges: []er.Relationship{},
	}
	entKeys := make([]string, 0, len(g.Entities))
	for e := range g.Entities {
		entKeys = append(entKeys, e)
	}
	slices.Sort(entKeys)
	for _, key := range entKeys {
		e := g.Entities[key]
		input.Nodes = append(input.Nodes, e)
	}
	for _, rSet := range g.OutEdges {
		for r, _ := range rSet.Relations {
			input.Edges = append(input.Edges, r)
		}
	}
	slices.SortFunc(input.Edges, func(l, r er.Relationship) int {
		return cmp.Or(
			cmp.Compare(l.Kind, r.Kind),
			cmp.Compare(l.A, r.A),
			cmp.Compare(l.Z, r.Z),
		)
	})

	buf := &bytes.Buffer{}
	if err := outputTmpl.Execute(buf, input); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	_, err = io.Copy(appCtx.App.Writer, buf)
	return err
}
