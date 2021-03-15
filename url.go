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
	"fmt"
	"net/url"
)

// setURL take the value of key s, validates the value is a URL, and the sets
// the value of key d to the validated value. If the value is not a URL then
// an error is returned.
func setURL(s string, d string, q *url.Values) error {
	u, err := validateURL(s, q.Get(s))
	if err != nil {
		return err
	}
	q.Set(d, u.String())
	return nil
}

// validateURL confirms that the parameter is a valid URL and then returns the
// URL ready for use with SWAN if valid. The method checks that the SWAN
// encrypted data can be appended to the end of the string as an identifiable
// segment.
func validateURL(n string, v string) (*url.URL, error) {
	if v == "" {
		return nil, fmt.Errorf("%s must be a valid URL", n)
	}
	u, err := url.Parse(v)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" {
		return nil, fmt.Errorf("%s '%s' must include a scheme", n, v)
	}
	if u.Host == "" {
		return nil, fmt.Errorf("%s '%s' must include a host", n, v)
	}

	// If the last character of the path is not a forward slash then append one.
	// The result of the storage operation is always appended to the path before
	// the query string.
	if u.Path[len(u.Path)-1:] != "/" {
		u.Path = u.Path + "/"
	}
	return u, nil
}
