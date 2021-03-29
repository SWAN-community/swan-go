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
	"owid"
	"strings"
	"swift"
	"time"
)

// Prefix added to the key for any SWAN values stored by the caller as cookies.
const cookiePrefix = "swan-"

// Pair represents a key value pair stored in SWAN. The created and expiry times
// for the value are also available.
type Pair struct {
	Key     string    // The name of the key associated with the value
	Created time.Time // The UTC time when the value was created
	// The UTC time when the value will expire and should not be used
	Expires time.Time
	Value   string // The value for the key as a string
}

// CookieName name for any cookie associated with the SWAN pair.
func (p *Pair) CookieName() string { return cookiePrefix + p.Key }

// IsSWANCookie returns true if a SWAN cookie.
func IsSWANCookie(c *http.Cookie) bool {
	return strings.HasPrefix(c.Name, cookiePrefix)
}

// NewPairFromCookie creates a new SWAN pair from the cookie.
func NewPairFromCookie(c *http.Cookie) *Pair {
	n := c.Name
	if IsSWANCookie(c) {
		n = c.Name[len(cookiePrefix):]
	}
	return &Pair{
		Key:   n,
		Value: c.Value,
	}
}

// NewPairFromSWIFT creates a new SWAN pair from the SWIFT pair setting the
// value to the byte array provided.
func NewPairFromSWIFT(s *swift.Pair, v string) *Pair {
	return &Pair{
		Key:     s.Key(),
		Created: s.Created(),
		Expires: s.Expires(),
		Value:   v}
}

// AsCookie returns the pair as a cookie to be used in an HTTP response.
func (p *Pair) AsCookie(
	r *http.Request,
	w http.ResponseWriter,
	s bool) *http.Cookie {
	return &http.Cookie{
		Name:     p.CookieName(),
		Domain:   getDomain(r.Host),    // Specifically to this domain
		Value:    p.Value,              // The value as a base 64 string
		SameSite: http.SameSiteLaxMode, // Available to all paths
		HttpOnly: false,
		Secure:   s, // Secure if HTTPs, otherwise false.
		// Set the cookie expiry time to the same as the SWAN pair.
		Expires: p.Expires,
	}
}

// AsOWID returns the Value as an OWID structure. Used for SWID, SID and
// Preferences. If the Value is not an OWID then an error is returned.
func (p *Pair) AsOWID() (*owid.OWID, error) {
	return owid.FromBase64(p.Value)
}

func getDomain(h string) string {
	s := strings.Split(h, ":")
	return s[0]
}
