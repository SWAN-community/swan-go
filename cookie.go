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
	"net/http"
	"strings"
	"time"

	"github.com/SWAN-community/swift-go"
)

// Prefix added to the key for any SWAN values stored by the caller as cookies.
const cookiePrefix = "swan-"

type Validity struct {
	Key     string    `json:"key"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

func (v *Validity) UnmarshalSwiftValidity(p *swift.Pair) error {
	v.Created = p.Created()
	v.Expires = p.Expires()
	return nil
}

// CookieName name for any cookie associated with the SWAN pair.
func (v *Validity) CookieName() string { return cookiePrefix + v.Key }

// AsHttpCookie creates a HTTP cookie that needs to have the Value field set
// to the base 64 data associated with the SWAN entity.
func (p *Validity) AsHttpCookie(host string, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     p.CookieName(),
		Domain:   getDomain(host),      // Specifically to this domain
		SameSite: http.SameSiteLaxMode, // Available to all paths
		HttpOnly: false,
		Secure:   secure, // Secure if HTTPs, otherwise false.
		// Set the cookie expiry time to the same as the SWAN pair.
		Expires: p.Expires}
}

// Remove any port information that may be included in the host as this is not
// used by cookies.
func getDomain(h string) string {
	s := strings.Split(h, ":")
	return s[0]
}

// IsSWANCookie returns true if a SWAN cookie.
func IsSWANCookie(c *http.Cookie) bool {
	return strings.HasPrefix(c.Name, cookiePrefix)
}
