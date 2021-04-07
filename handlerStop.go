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
	"net/http"
)

func handlerStop(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check caller is authorized to access SWAN.
		if s.getAccessAllowed(w, r) == false {
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
		t := s.config.DeleteDate().Format("2006-01-02")
		r.Form.Set(fmt.Sprintf("stop+%s", t), r.Form.Get("host"))
		r.Form.Set("message", fmt.Sprintf(
			"Bye, bye %s. Thanks for telling the world.",
			r.Form.Get("host")))
		r.Form.Del("host")

		// Uses the SWIFT access node associated with this internet domain
		// to determine the URL to direct the browser to.
		u, err := createStorageOperationURL(s.swift, r, r.Form)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Return the URL as a text string.
		sendResponse(s, w, "text/plain; charset=utf-8", []byte(u))
	}
}
