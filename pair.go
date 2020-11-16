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
	Value   string    // The value as a byte array
}

// CreatorURL returns a URL to get the creator details.
func (p *Pair) CreatorURL() string {
	o, _ := p.AsOWID()
	return "//" + o.Domain + "/owid.json"
}

// AsOWID returns the Value as an OWID structure.
func (p *Pair) AsOWID() (*owid.OWID, error) {
	return owid.DecodeFromBase64(p.Value)
}

// VerifyURL returns a URL that can be called to verify the OWID.
func (p *Pair) VerifyURL() string {
	o, _ := p.AsOWID()
	return "//" + o.Domain + "/owid/api/v1/verify?owid=" + p.Value
}

// DecodeAndVerifyURL returns a URL that can be called to decode the OWID.
func (p *Pair) DecodeAndVerifyURL() string {
	o, _ := p.AsOWID()
	return "//" + o.Domain + "/owid/api/v1/decode-and-verify?owid=" + p.Value
}
