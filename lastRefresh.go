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

import "net/http"

// LastRefresh of the SWAN data. Used by the publisher to determine if a refresh
// of SWAN data is required. The created date is the content of the cookie, and
// the expiry date which is used by the browser to clear the cookie is when the
// SWAN data should be refreshed. The presence of this cookie indicates the data
// can be used.
type LastRefresh struct {
	Cookie
}

// GetCookie returns the time cookie with the key set to "ref".
func (l *LastRefresh) GetCookie() *Cookie {
	l.Key = "ref"
	return &l.Cookie
}

// AsHttpCookie sets the value of the cookie to the created date as a string.
func (l *LastRefresh) AsHttpCookie(host string, secure bool) *http.Cookie {
	c := l.GetCookie().AsHttpCookie(host, secure)
	c.Value = l.Created.String()
	return c
}
