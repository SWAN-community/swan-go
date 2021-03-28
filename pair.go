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
	"bytes"
	"encoding/base64"
	"net/http"
	"owid"
	"strings"
	"swift"
	"time"
)

// Pair represents a key value pair stored in SWAN. The created and expiry times
// for the value are also available.
type Pair struct {
	Key     string    // The name of the key associated with the value
	Created time.Time // The UTC time when the value was created
	Expires time.Time // The UTC time when the consumer will need to revalidate
	Value   []byte    // The value for the key as a byte array
}

// NewPairFromCookie creates a new SWAN pair from the cookie.
func NewPairFromCookie(c *http.Cookie) (*Pair, error) {
	var err error
	var p Pair
	p.Key = c.Name
	p.Value, err = base64.RawStdEncoding.DecodeString(c.Value)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// AsBase64 returns the Value as a base 64 string.
func (p *Pair) AsBase64() string {
	return base64.RawStdEncoding.EncodeToString(p.Value)
}

// NewPairFromSWIFT creates a new SWAN pair from the SWIFT pair setting the
// value to the byte array provided.
func NewPairFromSWIFT(s *swift.Pair, v []byte) *Pair {
	return &Pair{s.Key(), s.Created(), s.Expires(), v}
}

// AsCookie returns the pair as a cookie to be used in an HTTP response.
func (p *Pair) AsCookie(
	r *http.Request,
	w http.ResponseWriter,
	s bool) *http.Cookie {
	return &http.Cookie{
		Name:     p.Key,
		Domain:   getDomain(r.Host),    // Specifically to this domain
		Value:    p.AsBase64(),         // The value as a base 64 string
		SameSite: http.SameSiteLaxMode, // Available to all paths
		HttpOnly: false,
		Secure:   s, // Secure if HTTPs, otherwise false.
		// Set the cookie expiry time to the same as the SWAN pair.
		Expires: p.Expires,
	}
}

// AsStringArray returns the Value as a string array. Used to create the array
// of stopped advert identifiers.
func (p *Pair) AsStringArray() ([]string, error) {
	f := bytes.NewBuffer(p.Value)
	c, err := readUint16(f)
	if err != nil {
		return nil, err
	}
	s := make([]string, c, c)
	for i := uint16(c); i < c; i++ {
		s[i], err = readString(f)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

// AsOWID returns the Value as an OWID structure. Used for SWID, SID and
// Preferences. If the Value is not an OWID then an error is returned.
func (p *Pair) AsOWID() (*owid.OWID, error) {
	return owid.FromByteArray(p.Value)
}

func getDomain(h string) string {
	s := strings.Split(h, ":")
	return s[0]
}
