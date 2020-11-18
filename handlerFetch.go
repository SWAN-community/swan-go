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
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

func handlerFetch(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create the URL with the parameters provided by the publisher.
		u, err := createStorageOperationURL(
			s,
			r.URL.RawQuery,
			func(q *url.Values) {
				t := time.Now().UTC().AddDate(0, 3, 0).Format("2006-01-02")
				q.Set(fmt.Sprintf("cbid<%s", t), uuid.New().String())
				q.Set(fmt.Sprintf("uuid<%s", t), "")
				q.Set(fmt.Sprintf("allow<%s", t), "")
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

func createStorageOperationURL(
	s *services,
	p string,
	fn func(q *url.Values)) (string, error) {

	// Check that an access node exists for SWAN. If not try to update the
	// access node before erroring.
	if s.accessNode == "" {
		an, err := s.swift.GetAccessNode(s.config.Network)
		if err != nil && an == "" {
			return "", fmt.Errorf("An access node has not been created for the"+
				" '%s' network. Use http[s]://[domain]/swift/register to start"+
				" the network.",
				s.config.Network)
		}
		s.accessNode = an
	}

	// Build a new URL to request the first storage operation URL.
	u, err := url.Parse(s.config.Scheme + "://" + s.accessNode)
	if err != nil {
		return "", err
	}
	u.Path = "/swift/api/v1/create"

	// Use the function passed to the method to add any additional query
	// parameters.
	q, err := url.ParseQuery(p)
	if err != nil {
		return "", err
	}
	fn(&q)

	// Copy the query string parameters exactly from those provided by the
	// publisher.
	u.RawQuery = q.Encode()

	// Get the first storage URL from the access node.
	res, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", newResponseError(u.String(), res)
	}

	// Read the response as a string.
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
