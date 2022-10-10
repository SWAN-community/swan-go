/* ****************************************************************************
 * Copyright 2022 51 Degrees Mobile Experts Limited (51degrees.com)
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
	"net/http"
	"strings"
	"time"

	"github.com/SWAN-community/swift-go"
)

// Prefix added to the key for any SWAN values stored by the caller as cookies.
const cookiePrefix = "swan-"

type Cookie struct {
	Key     string    `json:"key,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Expires time.Time `json:"expires,omitempty"`
}

// asHttpCookie returns the validity instance as a cookie.
// f is an instance of a structure that implements Field.
func (c *Cookie) asHttpCookie(h string, s bool, f Field) (*http.Cookie, error) {
	if c == nil {
		return nil, fmt.Errorf("nil cookie")
	}
	b := c.AsHttpCookie(h, s)
	d, err := f.MarshalBase64()
	if err != nil {
		return nil, err
	}
	b.Value = string(d)
	return b, nil
}

func (c *Cookie) UnmarshalSwiftValidity(p *swift.Pair) error {
	c.Created = p.Created()
	c.Expires = p.Expires()
	return nil
}

// CookieName name for any cookie associated with the SWAN pair.
func (c *Cookie) CookieName() string { return cookiePrefix + c.Key }

// AsHttpCookie creates a HTTP cookie that needs to have the Value field set
// to the base 64 data associated with the SWAN entity.
func (c *Cookie) AsHttpCookie(host string, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     c.CookieName(),
		Domain:   getDomain(host),      // Specifically to this domain
		SameSite: http.SameSiteLaxMode, // Available to all paths
		HttpOnly: false,
		Secure:   secure, // Secure if HTTPs, otherwise false.
		// Set the cookie expiry time to the same as the SWAN pair.
		Expires: c.Expires}
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
