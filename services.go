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
	"owid"
	"swift"
)

// HTTP headers that if present indicate a request is probably from a web
// browser.
var invalidHTTPHeaders = []string{
	"Accept",
	"Accept-Language",
	"Cookie"}

// Services references all the information needed for every method.
type services struct {
	config Configuration
	swift  *swift.Services // Services used by the SWIFT network
	owid   *owid.Services  // Services used for OWID creation and verification
	access Access          // Instance of access service
}

// newServices a set of services to use with SWAN. These provide defaults via
// the configuration parameter, and access to persistent storage via the store
// parameter.
func newServices(settingsFile string, swanAccess Access) *services {
	var swiftStore swift.Store
	var owidStore owid.Store

	// Use the file provided to get the SWIFT settings.
	swiftConfig := swift.NewConfig(settingsFile)
	err := swiftConfig.Validate()
	if err != nil {
		panic(err)
	}

	// Use the file provided to get the OWID settings.
	owidConfig := owid.NewConfig(settingsFile)
	err = owidConfig.Validate()
	if err != nil {
		panic(err)
	}

	// Link to the SWIFT storage.
	swiftStore = swift.NewStore(swiftConfig)

	// Link to the OWID storage.
	owidStore = owid.NewStore(owidConfig)

	// Get the default browser detector.
	b, err := swift.NewBrowserRegexes()
	if err != nil {
		panic(err)
	}

	// Create the swan configuration.
	c := newConfig(settingsFile)

	// Return the services.
	return &services{
		c,
		swift.NewServices(swiftConfig, swiftStore, swanAccess, b),
		owid.NewServices(owidConfig, owidStore, swanAccess),
		swanAccess}
}

// Returns true if the request is allowed to access the handler, otherwise
// false. Removes the accessKey parameter from the form to prevent it being
// used by other methods.  If false is returned then no further action is
// needed as the method will have responded to the request already.
func (s *services) getAccessAllowed(
	w http.ResponseWriter,
	r *http.Request) bool {

	// Check that there are no HTTP headers that are usually sent by browsers.
	// SWAN can only be used from server side environments to ensure that the
	// accessKey does not become publicly available.
	if s.config.Debug == false {
		for _, h := range invalidHTTPHeaders {
			if r.Header.Get(h) != "" {
				returnAPIError(&s.config, w,
					fmt.Errorf(
						"'%s' header must not be present in SWAN API requests as "+
							"this indicates that the request is coming from a web "+
							"browser and therefore the access key might be "+
							"compromised if this configuration where to be made "+
							"publicly available",
						h),
					http.StatusNetworkAuthenticationRequired)
				return false
			}
		}
	}

	// Check that the domain for this request relates to a valid access node.
	a, err := s.swift.GetAccessNodeForHost(r.Host)
	if err != nil {
		returnAPIError(
			&s.config,
			w,
			err,
			http.StatusBadRequest)
		return false
	}
	if a == nil {
		returnAPIError(
			&s.config,
			w,
			fmt.Errorf("'%s' not a valid SWAN access node", r.Host),
			http.StatusBadRequest)
		return false
	}

	// Validate that the access key provided is valid in the access provider.
	err = r.ParseForm()
	if err != nil {
		returnAPIError(&s.config, w, err, http.StatusInternalServerError)
		return false
	}
	v, err := s.access.GetAllowed(r.FormValue("accessKey"))
	if v == false || err != nil {
		returnAPIError(&s.config, w,
			fmt.Errorf("Access denied. Verify parameter accessKey"),
			http.StatusNetworkAuthenticationRequired)
		return false
	}

	r.Form.Del("accessKey")

	return true
}
