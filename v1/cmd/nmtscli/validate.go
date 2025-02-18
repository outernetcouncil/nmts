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
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
	"outernetcouncil.org/nmts/v1/lib/validation"
)

func validateGraph(appCtx *cli.Context) error {
	g, err := readGraphWithValidator(appCtx, validation.DefaultValidator{})
	if err != nil {
		return err
	}

	errs := []error{}
	if g.NumEntities() < 1 {
		errs = append(errs, fmt.Errorf("graph had no entities"))
	}
	if g.NumRelationships() < 1 {
		errs = append(errs, fmt.Errorf("graph had no relationships"))
	}

	return errors.Join(errs...)
}
