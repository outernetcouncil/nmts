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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	// nosemgrep: import-text-template
	"text/template"

	"google.golang.org/protobuf/encoding/prototext"
	npb "outernetcouncil.org/nmts/v0/proto"
)

var marshalOptions = prototext.MarshalOptions{
	Multiline: true,
	Indent:    "  ",
}

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout io.Writer) error {
	// Define a flag to hold the Fragment txtpb template filename.
	fs := flag.NewFlagSet("template2txtpb", flag.ExitOnError)
	tmplFilenames := []string{}
	fs.Func("tmpl_filename", "Fragment txtpb template filename", func(v string) error {
		tmplFilenames = append(tmplFilenames, v)
		return nil
	})
	outFile := fs.String("output", "", "Output file (stdout if blank)")
	inFile := fs.String("input", "", "Input file (stdin if blank)")

	skipProtoValidation := fs.Bool("skip_proto_validation", false, "Skip validating the template output as a protobuf message and emit it immediately")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parsing flags: %w", err)
	}

	// Read in the template file.
	t, err := template.New(filepath.Base(tmplFilenames[0])).
		Option("missingkey=error").
		Funcs(map[string]any{
			"mk_slice": func(v ...any) []any { return v },
		}).
		ParseFiles(tmplFilenames...)
	if err != nil {
		return fmt.Errorf("parsing Fragment txtpb template file: %w", err)
	}

	var input io.Reader
	if *inFile == "" {
		input = stdin
	} else {
		inf, err := os.Open(*inFile)
		if err != nil {
			return err
		}
		defer inf.Close()
		input = inf
	}

	var output io.Writer
	if *outFile == "" {
		output = stdout
	} else {
		outf, err := os.Create(*outFile)
		if err != nil {
			return err
		}
		defer outf.Close()
		output = outf
	}

	// Read in JSON data from stdin.
	jsonData, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("reading JSON input file: %w", err)
	}

	// Parse JSON data.
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("parsing JSON input file: %w", err)
	}

	// Execute the template with the JSON data into a local buffer,
	// and assert no missing template substitutions.
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, data); err != nil {
		return fmt.Errorf("executing template with supplied JSON data: %w", err)
	}

	if *skipProtoValidation {
		fmt.Fprintln(output, buffer.String())
		return nil
	}

	// Parse the local buffer as a Fragment message.
	fragmentMsg := &npb.Fragment{}
	if err := prototext.Unmarshal(buffer.Bytes(), fragmentMsg); err != nil {
		return fmt.Errorf("parsing executed template as Fragment: %w", err)
	}

	// Prettyprint the Fragment message to output
	text, err := marshalOptions.Marshal(fragmentMsg)
	if err != nil {
		return fmt.Errorf("prettyprinting Fragment: %w", err)
	}
	fmt.Fprintln(output, string(text))
	return nil
}
