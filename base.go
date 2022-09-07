/* ****************************************************************************
 * Copyright 2020 51 Degrees Mobile Experts Limited (51degrees.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 * ***************************************************************************/

package swan

import (
	"fmt"

	"github.com/SWAN-community/owid-go"
)

// Base used with any SWAN field.
type Base struct {
	Version byte       `json:"version"` // Used to indicate the version encoding of the type.
	OWID    *owid.OWID `json:"source"`  // OWID related to the structure
}

// errorMissing function to create error messages for missing JSON keys.
func errorMissing(name string) error {
	return fmt.Errorf("'%s' missing", name)
}

// errorInvalid function to create error messages for invalid JSON keys.
func errorInvalid(name string, typeName string) error {
	return fmt.Errorf("'%s' invalid for type '%s'", name, typeName)
}
