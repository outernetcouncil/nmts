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

// Package main provides a general purpose CLI tool for working with NMTS
// graphs.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
	er "outernetcouncil.org/nmts/v1alpha/lib/entityrelationship"
)

const appName = "nmtscli"

func main() {
	if err := App(os.Stdin, os.Stdout, os.Stderr).Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %v\n", err)
		os.Exit(1)
	}
}

func App(stdin io.Reader, stdout, stderr io.Writer) *cli.App {
	return &cli.App{
		Name: appName,
		Commands: []*cli.Command{
			{
				Name: "export",
				Subcommands: []*cli.Command{
					{
						Name:   "dot",
						Action: exportDot,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "rankdir",
							},
						},
					},
					{
						Name:   "d2",
						Action: exportD2,
					},
					{
						Name:   "html",
						Action: exportHtml,
					},
					{
						Name:   "nquads",
						Action: exportNQuads,
					},
					{
						Name:   "prolog",
						Action: exportProlog,
					},
				},
			},
			{
				Name:   "validate",
				Action: validateGraph,
			},
		},
	}
}

func readGraph(appCtx *cli.Context) (*er.Collection, error) {
	return readGraphWithValidator(appCtx, nil)
}

func readGraphWithValidator(appCtx *cli.Context, v er.Validator) (*er.Collection, error) {
	srcs := appCtx.Args().Slice()
	if len(srcs) == 0 {
		return nil, fmt.Errorf("missing input files")
	}

	g, err := er.ReadFragmentFiles(srcs)
	if err != nil {
		return nil, err
	}

	collBldr := er.NewCollectionBuilder(v)
	if err := collBldr.InsertFragments(g); err != nil {
		return nil, err
	}
	return collBldr.Build()
}
