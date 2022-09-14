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
)

// Prefix added to the key for any SWAN values stored by the caller as cookies.
const cookiePrefix = "swan-"

// Pair represents a key value pair stored in SWAN. The created and expiry times
// for the value are also available.
type Pair struct {
	Key     string    // Key without any prefix for SWAN
	Created time.Time // The UTC time when the value was created
	// The UTC time when the value will expire and should not be used
	Expires time.Time
	Value   string // See https://tools.ietf.org/html/rfc6265 for details.
}

// CookieName name for any cookie associated with the SWAN pair.
func (p *Pair) CookieName() string { return cookiePrefix + p.Key }

// IsSWANCookie returns true if a SWAN cookie.
func IsSWANCookie(c *http.Cookie) bool {
	return strings.HasPrefix(c.Name, cookiePrefix)
}

// NewPairFromCookie creates a new SWAN pair from the key and field value.
func NewPairFromField(key string, value Field) (*Pair, error) {
	b, err := value.MarshalBase64()
	if err != nil {
		return nil, err
	}
	return &Pair{Key: key, Value: string(b)}, nil
}

// NewPairFromCookie creates a new SWAN pair from the cookie.
// cookie as source for the pair.
func NewPairFromCookie(cookie *http.Cookie) (*Pair, error) {
	n := cookie.Name
	if IsSWANCookie(cookie) {
		n = cookie.Name[len(cookiePrefix):]
	}
	return &Pair{Key: n, Value: cookie.Value}, nil
}

// AsCookie returns the pair as a cookie to be used in an HTTP response.
// host to use for the domain of the cookie.
// secure
func (p *Pair) AsCookie(host string, secure bool) (*http.Cookie, error) {
	return &http.Cookie{
		Name:     p.CookieName(),
		Domain:   getDomain(host), // Specifically to this domain
		Value:    p.Value,
		SameSite: http.SameSiteLaxMode, // Available to all paths
		HttpOnly: false,
		Secure:   secure, // Secure if HTTPs, otherwise false.
		// Set the cookie expiry time to the same as the SWAN pair.
		Expires: p.Expires}, nil
}

// UnmarshalBase64 the value string into the instance of field provided. Will
// fail if the value is either not base 64 or not of the same type as the field
// instance provided.
func (p *Pair) UnmarshalBase64(field Field) error {
	return field.UnmarshalBase64([]byte(p.Value))
}

// Remove any port information that may be included in the host as this is not
// used by cookies.
func getDomain(h string) string {
	s := strings.Split(h, ":")
	return s[0]
}
