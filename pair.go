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
	"owid"
	"time"
)

// Pair represents a key value pair stored in SWAN. The created and expiry times
// for the value are also available.
type Pair struct {
	Key     string    // The name of the key associated with the value
	Created time.Time // The UTC time when the value was created
	Expires time.Time // The UTC time when the consumer will need to revalidate
	Value   string    // The OWID value as a base 64 string
}

// AsOWID returns the Value as an OWID structure.
func (p *Pair) AsOWID() (*owid.OWID, error) {
	return owid.FromBase64(p.Value)
}

// AsBase64 returns the Value OWID as a base 64 string.
func (p *Pair) AsBase64() (string, error) {
	o, err := p.AsOWID()
	if err != nil {
		return "", err
	}
	return o.AsBase64()
}
