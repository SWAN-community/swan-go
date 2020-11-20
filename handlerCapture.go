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
	"net/url"
	"strings"
	"swift"
	"time"

	"github.com/google/uuid"
)

type model struct {
	url.Values
}

func (m *model) Title() string           { return m.Get("title") }
func (m *model) CBID() string            { return m.Get("cbid") }
func (m *model) Email() string           { return m.Get("email") }
func (m *model) Allow() string           { return m.Get("allow") }
func (m *model) BackgroundColor() string { return m.Get("backgroundColor") }

func handlerCapture(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the model from the URL.
		var m model
		m.Values = make(url.Values)

		// The last path segment is the data.
		l := strings.LastIndex(r.URL.Path, "/")
		if l < 0 {
			returnRequestError(&s.config, w, nil, http.StatusBadRequest)
			return
		}

		// Decrypt the data. If not possible return a bad request error.
		in, err := decrypt(s, r.URL.Path[l+1:])
		if err != nil {
			returnRequestError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Get the results.
		res, err := swift.DecodeResults(in)
		if err != nil {
			returnRequestError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Build the form parameters from the data received from SWIFT.
		m.Set("title", res.HTML.Title)
		m.Set("backgroundColor", res.HTML.BackgroundColor)
		m.Set("messageColor", res.HTML.MessageColor)
		m.Set("progressColor", res.HTML.ProgressColor)
		m.Set("message", res.HTML.Message)
		m.Set("returnUrl", res.State)
		m.Set("cbid", res.Get("cbid").Value)
		m.Set("email", res.Get("email").Value)
		m.Set("allow", res.Get("allow").Value)

		// Respond based on the method used.
		switch r.Method {
		case "GET":
			handlerCaptureGet(s, w, r, &m)
		case "POST":
			handlerCapturePost(s, w, r, &m)
		}
	}
}

func handlerCaptureGet(
	s *services,
	w http.ResponseWriter,
	r *http.Request,
	m *model) {

	// Display the user interface with the data provided.
	err := captureTemplate.Execute(w, m)
	if err != nil {
		returnServerError(&s.config, w, err)
		return
	}
}

func handlerCapturePost(
	s *services,
	w http.ResponseWriter,
	r *http.Request,
	m *model) {

	// Get the data provided in the post back.
	err := r.ParseForm()
	if err != nil {
		returnAPIError(&s.config, w, err, http.StatusUnprocessableEntity)
		return
	}

	// Update the model with the values from the form.
	m.Set("email", r.Form.Get("email"))
	m.Set("allow", r.Form.Get("allow"))
	m.Set("cbid", r.Form.Get("cbid"))

	// Check to see if the post is as a result of the CBID reset.
	if r.Form.Get("reset-cbid") != "" {

		// Replace the CBID with a new random value.
		m.Set("cbid", uuid.New().String())

		// Display the template again.
		err = captureTemplate.Execute(w, &m)
		if err != nil {
			returnServerError(&s.config, w, err)
		}
		return
	}

	// Check to see if the post is as a result for all data.
	if r.Form.Get("reset-all") != "" {

		// Replace the data.
		m.Set("email", "")
		m.Set("allow", "")
		m.Set("cbid", uuid.New().String())

		// Display the template again.
		err = captureTemplate.Execute(w, &m)
		if err != nil {
			returnServerError(&s.config, w, err)
		}
		return
	}

	// Create the URL to use to write the values to Swift.
	u, err := createStorageOperationURL(s, &m.Values,
		func(q *url.Values) {

			// Add any parameters from the form being posted back with a common
			// date for expiry in 3 months.
			t := time.Now().UTC().AddDate(0, 3, 0).Format("2006-01-02")

			// Add the Common Browser ID replacing any existing values if
			// present in the network.
			q.Set(fmt.Sprintf("cbid>%s", t), r.PostForm.Get("cbid"))

			// Add the email so that it can verified as some time in the future
			// if necessary by the SWAN provider.
			q.Set(fmt.Sprintf("email>%s", t), r.PostForm.Get("email"))

			if r.PostForm.Get("allow") == "" {
				q.Set(fmt.Sprintf("allow>%s", t), "off")
			} else {
				q.Set(fmt.Sprintf("allow>%s", t), r.PostForm.Get("allow"))
			}

			// Delete the keys that were provided from the publisher so that the
			// conflict resolution policy and date can be applied.
			q.Del("cbid")
			q.Del("email")
			q.Del("allow")

			// Add the access key for the SWIFT network.
			q.Set("accessKey", s.config.AccessKey)
		})
	if err != nil {
		returnServerError(&s.config, w, err)
		return
	}

	// Redirect the browser window to start the write process.
	http.Redirect(w, r, u, 303)
}
