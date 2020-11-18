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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

func handlerUpdate(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check caller can access
		if s.getAccessAllowed(w, r) == false {
			returnAPIError(&s.config, w,
				errors.New("Not authorized"),
				http.StatusUnauthorized)
			return
		}

		// Get the form values from the input request.
		err := r.ParseForm()
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}

		// Copy the incoming parameters into the outgoing ones.
		q, err := url.ParseQuery(r.Form.Encode())
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}

		// Create the URL with the parameters provided by the publisher.
		u, err := createStorageOperationURL(
			s,
			&q,
			func(q *url.Values) {
				t := time.Now().UTC().AddDate(0, 3, 0).Format("2006-01-02")
				q.Set(fmt.Sprintf("cbid<%s", t), uuid.New().String())
				q.Set(fmt.Sprintf("email<%s", t), "")
				q.Set(fmt.Sprintf("allow<%s", t), "")

				// Store the return URL from the publisher in the state for the
				// the storage operation in SWIFT. Then replace the return URL
				// with the SWAN preferences page URL.
				q.Set("state", q.Get("returnUrl"))
				q.Set("returnUrl",
					s.config.Scheme+"://"+r.Host+"/swan/preferences/")
			})
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusUnprocessableEntity)
			return
		}

		// Return the URL as a text string.
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Write([]byte(u))
	}
}
