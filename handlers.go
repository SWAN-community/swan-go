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
	"owid"
	"swift"
)

// AddHandlers adds swift and owid end points configured from the JSON file
// provided.
func AddHandlers(
	settingsFile string,
	malformedHandler func(w http.ResponseWriter, r *http.Request)) {

	// Create the new set of services.
	s := newServices(settingsFile)

	// Add the SWIFT handlers.
	swift.AddHandlers(s.swift, malformedHandler)

	// Add the OWID handlers.
	owid.AddHandlers(s.owid)

	// Add the SWAN handlers.
	http.HandleFunc("/swan/api/v1/fetch", handlerFetch(s))
	http.HandleFunc("/swan/api/v1/decrypt", handlerDecrypt(s))
	http.HandleFunc("/swan/api/v1/create-offer-id", handlerCreateOfferID(s))
	http.HandleFunc("/swan/preferences", handlerCapture(s))
}

func newResponseError(url string, resp *http.Response) error {
	in, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("API call '%s' returned '%d' and '%s'",
		url, resp.StatusCode, in)
}

func returnAPIError(
	c *Configuration,
	w http.ResponseWriter,
	err error,
	code int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(w, err.Error(), code)
	if c.Debug {
		println(err.Error())
	}
}

func returnServerError(c *Configuration, w http.ResponseWriter, err error) {
	w.Header().Set("Cache-Control", "no-cache")
	if c.Debug {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		http.Error(w, "", http.StatusInternalServerError)
	}
	if c.Debug {
		println(err.Error())
	}
}
