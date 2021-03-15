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
	"fmt"
	"net/http"
	"time"
)

func handlerStop(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check caller is authorized to access SWAN.
		if s.getAccessAllowed(w, r) == false {
			returnAPIError(&s.config, w,
				errors.New("Not authorized"),
				http.StatusUnauthorized)
			return
		}

		// Validate the host parameter is present.
		if r.Form.Get("host") == "" {
			returnAPIError(
				&s.config,
				w,
				fmt.Errorf("'host' must be provided"),
				http.StatusBadRequest)
			return
		}

		// Validate the set the return URL.
		err := setURL("returnUrl", "returnUrl", &r.Form)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Create the URL with the parameters provided by the publisher.
		t := time.Now().UTC().AddDate(0, 3, 0).Format("2006-01-02")
		r.Form.Set(fmt.Sprintf("stop+%s", t), r.Form.Get("host"))
		r.Form.Set("message", fmt.Sprintf(
			"Bye, bye %s. Thanks for telling the world.",
			r.Form.Get("host")))
		r.Form.Del("host")
		u, err := createStorageOperationURL(s.swift, r, r.Form)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}

		// Return the URL as a text string.
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
