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
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/ichiban/prolog"
	"github.com/urfave/cli/v2"
	er "outernetcouncil.org/nmts/v0/lib/entityrelationship"
	npb "outernetcouncil.org/nmts/v0/proto"
)

const (
	physicalLayerClass = "deeppurple"
	linkLayerClass     = "orange"
	networkLayerClass  = "teal"
	containerClass     = "googleblue"
)

func graphToPrologFacts(g *er.Collection) string {
	w := &strings.Builder{}

	kinds := map[string][]string{}
	for _, e := range g.Entities {
		name := strings.ToLower(e.Id)
		ename := entityNameToAtom(name)
		kname := strings.ToLower(er.EntityKindStringFromProto(e))
		kinds[kname] = append(kinds[kname], ename)

		fmt.Fprintf(w, "is_entity(%s).\n", ename)
	}
	for _, e := range g.Entities {
		name := strings.ToLower(e.Id)
		ename := entityNameToAtom(name)
		fmt.Fprintf(w, "display_name(%s, %q).\n", ename, e.Id)
	}
	for k, subjs := range kinds {
		for _, s := range subjs {
			fmt.Fprintf(w, "%s(%s).\n", k, s)
		}
	}

	edges := map[npb.RK][][]string{}

	for _, rs := range g.OutEdges {
		for r := range rs.Relations {
			a, z := entityNameToAtom(r.A), entityNameToAtom(r.Z)
			edges[r.Kind] = append(edges[r.Kind], []string{a, z})
		}
	}

	for rk, subjs := range edges {
		rname := strings.ToLower(rk.String())
		for _, s := range subjs {
			fmt.Fprintf(w, "%s(%s, %s).\n", rname, s[0], s[1])
		}
	}

	for rk := range edges {
		rname := strings.ToLower(rk.String())
		fmt.Fprintf(w, `
connected(Left, Right, %[1]s) :- %[1]s(Left, Right).
`, rname)
	}

	for rk := range edges {
		rname := strings.ToLower(rk.String())
		fmt.Fprintf(w, `
reachable_via(Left, Right, %[1]s) :- %[1]s(Left, Right).
reachable_via(Left, Right, %[1]s) :- %[1]s(Left, Middle), %[1]s(Middle, Right).
`, rname)
	}

	fmt.Fprintln(w, `is_physical(E) :- ek_platform(E) ; ek_physical_medium_link(E).`)
	fmt.Fprintln(w, `is_logical(E) :- is_entity(E), \+ is_physical(E).`)

	// Entity E is a root iff it's the LHS of an rk_contains relationship with
	// child C *and* it is not the RHS of an rk_contains relationship with
	// parent P
	fmt.Fprintln(w, `is_root(E) :- rk_contains(E, C), \+ rk_contains(P, E), E @< C.`)

	return w.String()
}

func exportProlog(appCtx *cli.Context) error {
	g, err := readGraph(appCtx)
	if err != nil {
		return err
	}

	rules := graphToPrologFacts(g)

	_, err = fmt.Fprintln(appCtx.App.Writer, rules)
	return err
}

func entityNameToAtom(s string) string {
	s = strings.ReplaceAll(s, "->", "__")
	s = strings.ReplaceAll(s, "(", "_")
	s = strings.ReplaceAll(s, ")", "_")
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ToLower(s)
	return s
}

// collectQuery is a generic shortcut for executing a query and collecting its
// results.
func collectQuery[T any](p *prolog.Interpreter, q string, args ...interface{}) ([]T, error) {
	results := []T{}
	sols, err := p.Query(q, args...)
	if err != nil {
		return nil, err
	}
	for sols.Next() {
		var v T
		if err := sols.Scan(&v); err != nil {
			return nil, err
		}
		results = append(results, v)
	}
	return results, nil
}

type prologEntity struct {
	ID   string
	Name string
}

func compareEntity(l, r prologEntity) int {
	return cmp.Or(cmp.Compare(l.ID, r.ID), cmp.Compare(l.Name, r.Name))
}

func displayName(p *prolog.Interpreter, id string) (string, error) {
	sol := p.QuerySolution(fmt.Sprintf(`display_name(%s, A).`, id))
	var name struct{ A string }
	if err := sol.Scan(&name); err != nil {
		return "", err
	}
	return name.A, nil
}

func toNodeClass(p *prolog.Interpreter, id string) (string, error) {
	layerOne, err := collectQuery[struct{ B bool }](p, fmt.Sprintf(`rk_traverses(E, %[1]s).`, id))
	if err != nil {
		return "", err
	}
	if len(layerOne) > 0 {
		return physicalLayerClass, nil
	}

	layerTwo, err := collectQuery[struct{ B bool }](p, fmt.Sprintf(`rk_traverses(%[1]s, ID).`, id))
	if err != nil {
		return "", err
	}
	if len(layerTwo) > 0 {
		return linkLayerClass, nil
	}

	return networkLayerClass, nil
}

func queryForRootContainers(p *prolog.Interpreter) (map[prologEntity][]prologEntity, error) {
	type root struct {
		Parent   string
		Children []string
	}
	results, err := collectQuery[root](p, `bagof(C, (is_root(Parent), reachable_via(Parent, C, rk_contains)), Children).`)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(results, func(l, r root) int {
		return cmp.Or(cmp.Compare(l.Parent, r.Parent), slices.Compare(l.Children, r.Children))
	})

	roots := map[prologEntity][]prologEntity{}
	for _, r := range results {
		pn, err := displayName(p, r.Parent)
		if err != nil {
			return nil, err
		}

		rEnt := prologEntity{Name: pn, ID: r.Parent}
		children := make([]prologEntity, 0, len(r.Children))
		for _, c := range r.Children {
			cn, err := displayName(p, c)
			if err != nil {
				return nil, err
			}
			children = append(children, prologEntity{Name: cn, ID: c})
		}
		roots[rEnt] = children
	}
	return roots, nil
}
