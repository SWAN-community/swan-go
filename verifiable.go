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

import "github.com/SWAN-community/owid-go"

// Verifiable interface used to identify any types that can be verified.
type Verifiable interface {

	// Returns the OWID associated with the instance. The returned OWID must
	// have the Target field set to the instance being verified.
	GetOWID() *owid.OWID

	// True if the entity has been signed, otherwise false.
	IsSigned() bool
}
