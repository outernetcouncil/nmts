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
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/ichiban/prolog"
	"github.com/urfave/cli/v2"
	er "outernetcouncil.org/nmts/v1/lib/entityrelationship"
)

const (
	physicalLayerClass = "deeppurple"
	linkLayerClass     = "orange"
	networkLayerClass  = "teal"
	containerClass     = "googleblue"
)

func graphToPrologFacts(g *er.Collection) []string {
	facts := []string{"set_prolog_flag(double_quotes, atom)."}
	for _, e := range g.Entities {
		facts = append(facts, fmt.Sprintf(`is_entity(%q, %q).`, strings.ToLower(er.EntityKindStringFromProto(e)), e.Id))
	}

	for _, rs := range g.OutEdges {
		for r := range rs.Relations {
			facts = append(facts, fmt.Sprintf(`edge(%q, %q, %q).`, r.A, r.Z, strings.ToLower(r.Kind.String())))
		}
	}

	facts = append(facts, `reachable_via(Left, Right, E) :- edge(Left, Right, E).`)
	facts = append(facts, `reachable_via(Left, Right, E) :- edge(Left, Middle, E), edge(Middle, Right, E).`)
	facts = append(facts, `is_physical(E) :- is_entity("ek_platform", E) ; is_entity("ek_physical_medium_link", E).`)
	facts = append(facts, `is_logical(E) :- is_entity(K, E), \+ is_physical(E).`)
	// Entity E is a root iff it's the LHS of an rk_contains relationship with
	// child C *and* it is not the RHS of an rk_contains relationship with
	// parent P
	facts = append(facts, `is_root(E) :- reachable_via(E, C, "rk_contains"), \+ reachable_via(P, E, "rk_contains"), E @< C.`)

	return facts
}

func addGraphToPrologInterpreter(p *prolog.Interpreter, g *er.Collection) error {
	return p.Exec(strings.Join(graphToPrologFacts(g), "\n"))
}

func exportProlog(appCtx *cli.Context) error {
	g, err := readGraph(appCtx)
	if err != nil {
		return err
	}

	fmt.Fprintln(appCtx.App.Writer, strings.Join(graphToPrologFacts(g), "\n"))
	return nil
}

// collectQuery is a generic shortcut for executing a query and collecting its
// results.
func collectQuery[T any](p *prolog.Interpreter, q string, args ...interface{}) (results []T, resErr error) {
	sols, err := p.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := sols.Close()
		resErr = errors.Join(resErr, sols.Err(), closeErr)
	}()

	for sols.Next() {
		var v T
		if err := sols.Scan(&v); err != nil {
			return nil, err
		}
		results = append(results, v)
	}
	return results, nil
}

type prologEntity struct{ ID string }

func compareEntity(l, r prologEntity) int { return cmp.Compare(l.ID, r.ID) }

func toNodeClass(p *prolog.Interpreter, id string) (string, error) {
	layerOne, err := collectQuery[struct{ B bool }](p, fmt.Sprintf(`edge(E, %q, "rk_traverses").`, id))
	if err != nil {
		return "", err
	}
	if len(layerOne) > 0 {
		return physicalLayerClass, nil
	}

	layerTwo, err := collectQuery[struct{ B bool }](p, fmt.Sprintf(`edge(%q, ID, "rk_traverses").`, id))
	if err != nil {
		return "", err
	}
	if len(layerTwo) > 0 {
		return linkLayerClass, nil
	}

	return networkLayerClass, nil
}

func queryForRootContainers(p *prolog.Interpreter) (map[string][]string, error) {
	type root struct {
		Parent   string
		Children []string
	}
	results, err := collectQuery[root](p, `setof(C, (is_root(Parent), reachable_via(Parent, C, "rk_contains")), Children).`)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(results, func(l, r root) int {
		return cmp.Or(cmp.Compare(l.Parent, r.Parent), slices.Compare(l.Children, r.Children))
	})

	roots := map[string][]string{}
	for _, r := range results {
		roots[r.Parent] = r.Children
	}
	return roots, nil
}
