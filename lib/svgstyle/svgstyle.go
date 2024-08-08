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

// Package svgstyle provides a tool to embed a CSS stylesheet into an SVG
// document.
package svgstyle

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"errors"
	"io"
	"slices"
)

//go:embed graphviz-style.css
var stylesheet []byte

type styleTag struct {
	XMLName xml.Name `xml:"style"`
	Body    []byte   `xml:",cdata"`
}

func Embed(data []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	dec := xml.NewDecoder(bytes.NewReader(data))
	enc := xml.NewEncoder(buf)
decodeLoop:
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break decodeLoop
		} else if err != nil {
			return nil, err
		}

		if v, ok := tok.(xml.StartElement); ok && v.Name.Local == "svg" {
			newTok := xml.CopyToken(tok).(xml.StartElement)
			// something about this process causes a duplicate `xmlns`
			// attribute, so we remove that. The resulting svg has exactly one
			// `xmlns="http://www.w3.org/2000/svg"` attribute.
			newTok.Attr = slices.DeleteFunc(newTok.Attr, func(attr xml.Attr) bool { return attr.Name.Local == "xmlns" })

			if err := enc.EncodeToken(newTok); err != nil {
				return nil, err
			}
			if err := enc.Encode(styleTag{Body: stylesheet}); err != nil {
				return nil, err
			}
		} else {
			err = enc.EncodeToken(xml.CopyToken(tok))
		}
		if err != nil {
			return nil, err
		}
	}
	err := enc.Flush()
	return buf.Bytes(), err
}
