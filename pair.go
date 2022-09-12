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
	"encoding"
	"encoding/base64"
	"net/http"
	"strings"
	"time"
)

// Prefix added to the key for any SWAN values stored by the caller as cookies.
const cookiePrefix = "swan-"

// CookieValue contains interfaces that values must implement to enable use
// with a swan.Pair and cookies.
type CookieValue interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

// Pair represents a key value pair stored in SWAN. The created and expiry times
// for the value are also available.
type Pair struct {
	Key     string    // The name of the key associated with the value
	Created time.Time // The UTC time when the value was created
	// The UTC time when the value will expire and should not be used
	Expires time.Time
	Value   CookieValue // The value associated with the key
}

// CookieName name for any cookie associated with the SWAN pair.
func (p *Pair) CookieName() string { return cookiePrefix + p.Key }

// IsSWANCookie returns true if a SWAN cookie.
func IsSWANCookie(c *http.Cookie) bool {
	return strings.HasPrefix(c.Name, cookiePrefix)
}

// NewPairFromCookie creates a new SWAN pair from the cookie.
// cookie as source for the pair.
// value instance to be assigned to the pair.
func NewPairFromCookie(cookie *http.Cookie, value CookieValue) (*Pair, error) {
	n := cookie.Name
	if IsSWANCookie(cookie) {
		n = cookie.Name[len(cookiePrefix):]
	}
	b, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}
	err = value.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}
	return &Pair{Key: n, Value: value}, nil
}

// AsCookie returns the pair as a cookie to be used in an HTTP response.
// host to use for the domain of the cookie.
// secure
func (p *Pair) AsCookie(host string, secure bool) (*http.Cookie, error) {
	v, err := p.Value.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:     p.CookieName(),
		Domain:   getDomain(host),                      // Specifically to this domain
		Value:    base64.StdEncoding.EncodeToString(v), // The value as a base 64 string
		SameSite: http.SameSiteLaxMode,                 // Available to all paths
		HttpOnly: false,
		Secure:   secure, // Secure if HTTPs, otherwise false.
		// Set the cookie expiry time to the same as the SWAN pair.
		Expires: p.Expires}, nil
}

// Remove any port information that may be included in the host as this is not
// used by cookies.
func getDomain(h string) string {
	s := strings.Split(h, ":")
	return s[0]
}
