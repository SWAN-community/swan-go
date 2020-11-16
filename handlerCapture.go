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
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"owid"
	"time"
)

type model struct {
	request *http.Request // The HTTP request used to generate the page.
	title   string        // The title for the data capture form
}

func (m *model) Title() string { return m.title }
func (m *model) CBID() string  { return getOWIDValue(m.request.Form.Get("cbid")) }
func (m *model) Allow() string { return getOWIDValue(m.request.Form.Get("allow")) }
func (m *model) BackgroundColor() string {
	return m.request.Form.Get("backgroundColor")
}

func getOWIDValue(v string) string {
	o, err := owid.DecodeFromBase64(v)
	if err == nil {
		return string(o.Payload)
	}
	return ""
}

func handlerCapture(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlerCaptureGet(s, w, r)
		case "POST":
			handlerCapturePost(s, w, r)
		}
	}
}

func handlerCaptureGet(s *services, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		returnAPIError(&s.config, w, err, http.StatusUnprocessableEntity)
		return
	}
	err = captureTemplate.Execute(w, &model{r, r.Form.Get("title")})
	if err != nil {
		returnServerError(&s.config, w, err)
		return
	}
}

func handlerCapturePost(s *services, w http.ResponseWriter, r *http.Request) {

	// Get the data provided in the post back.
	err := r.ParseForm()
	if err != nil {
		returnAPIError(&s.config, w, err, http.StatusUnprocessableEntity)
		return
	}

	// Create the URL to use to write the values to Swift.
	u, err := createStorageOperationURL(s, r.URL.RawQuery,
		func(q *url.Values) {

			// Add any parameters from the form being posted back with a common
			// date for expiry in 3 months.
			t := time.Now().UTC().AddDate(0, 3, 0).Format("2006-01-02")

			// Add the Common Browser ID prefering any existing values if they
			// exist already.
			q.Set(fmt.Sprintf("cbid<%s", t), r.PostForm.Get("cbid"))

			// Add the email hash as a Universally Unique ID.
			q.Set(fmt.Sprintf("uuid>%s", t), createUUID(r.PostForm.Get("email")))

			if r.PostForm.Get("allow") == "" {
				q.Set(fmt.Sprintf("allow>%s", t), "off")
			} else {
				q.Set(fmt.Sprintf("allow>%s", t), r.PostForm.Get("allow"))
			}

			// Delete the keys that were provided from the publisher so that the
			// conflict resolution policy and date can be applied.
			q.Del("cbid")
			q.Del("uuid")
			q.Del("email")
			q.Del("allow")
		})
	if err != nil {
		returnServerError(&s.config, w, err)
		return
	}

	// Redirect the browser window to start the write process.
	http.Redirect(w, r, u, 303)
}

// TODO : What hashing algorithm do we want to use to turn email address into
// hashes?
func createUUID(v string) string {
	hasher := sha1.New()
	hasher.Write([]byte(v))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
