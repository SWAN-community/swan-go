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
	"compress/gzip"
	"errors"
	"net/http"
)

// handlerFetch returns a URL that can be used in the browser primary navigation
// to retrieve the most current data from the SWAN network. If no data is
// available default values are returned.
func handlerFetch(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		// Check caller is authorized to access SWAN.
		if s.getAccessAllowed(w, r) == false {
			returnAPIError(&s.config, w,
				errors.New("Not authorized"),
				http.StatusUnauthorized)
			return
		}

		// Validate the set the return URL.
		err = setURL("returnUrl", "returnUrl", &r.Form)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Set the SWAN fields to empty values that will only be used if no
		// other value is returned.
		setDefaults(&r.Form)

		// Uses the SWIFT access node associated with this internet domain
		// to determine the URL to direct the browser to.
		u, err := createStorageOperationURL(s.swift, r, r.Form)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}

		// Return the response from the SWIFT layer.
		g := gzip.NewWriter(w)
		defer g.Close()
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		_, err = g.Write([]byte(u))
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}
	}
}
