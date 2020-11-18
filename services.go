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
	"log"
	"net/http"
	"owid"
	"swift"
)

// Services references all the information needed for every method.
type services struct {
	config     Configuration
	swift      *swift.Services // Services used by the SWIFT network
	owid       *owid.Services  // Services used for OWID creation and verification
	accessNode string          // The access node for the SWAN network
	access     Access          // Instance of access service
}

// newServices a set of services to use with SWIFT. These provide defaults via
// the configuration parameter, and access to persistent storage via the store
// parameter.
func newServices(
	settingsFile string,
	swanAccess Access,
	swiftAccess swift.Access,
	owidAccess owid.Access) *services {
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
	swiftStore = configreSwiftStore(swiftConfig)

	// Link to the OWID storage.
	owidStore = configureOwidStore(owidConfig)

	// Get the default browser detector.
	b, err := swift.NewBrowserRegexes()
	if err != nil {
		panic(err)
	}

	// Create the swan configuration.
	c := newConfig(settingsFile)

	// Get the SWIFT access node for the SWAN network. Log any errors rather
	// than panic because it may be that a network has yet to be established
	// for SWAN in the storage tables.
	an, err := swiftStore.GetAccessNode(c.Network)
	if err != nil {
		log.Println(err.Error())
		log.Printf("Has a '%s' network been created?", c.Network)
	}

	// Return the services.
	return &services{
		c,
		swift.NewServices(swiftConfig, swiftStore, swiftAccess, b),
		owid.NewServices(owidConfig, owidStore, owidAccess),
		an,
		swanAccess}
}

func configreSwiftStore(swiftConfig swift.Configuration) swift.Store {
	var swiftStore swift.Store
	var err error

	if swiftConfig.AzureAccount != "" && swiftConfig.AzureAccessKey != "" {
		swiftStore, err = swift.NewAzure(
			swiftConfig.AzureAccount,
			swiftConfig.AzureAccessKey)
		if err != nil {
			panic(err)
		}
	} else if swiftConfig.UseDynamoDB {
		swiftStore, err = swift.NewAWS(swiftConfig.AWSRegion)
		if err != nil {
			panic(err)
		}
	}

	if swiftStore == nil {
		panic(errors.New("owid store not configured"))
	}

	return swiftStore
}

func configureOwidStore(owidConfig owid.Configuration) owid.Store {
	var owidStore owid.Store
	var err error

	if owidConfig.AzureAccount != "" && owidConfig.AzureAccessKey != "" {
		owidStore, err = owid.NewAzure(
			owidConfig.AzureAccount,
			owidConfig.AzureAccessKey)
		if err != nil {
			panic(err)
		}
	} else if owidConfig.UseDynamoDB {
		owidStore, err = owid.NewAWS(owidConfig.AWSRegion)
		if err != nil {
			panic(err)
		}
	}

	if owidStore == nil {
		panic(errors.New("owid store not configured"))
	}

	return owidStore
}

// Returns true if the request is allowed to access the handler, otherwise false.
// If false is returned then no further action is needed as the method will have
// responded to the request already.
func (s *services) getAccessAllowed(
	w http.ResponseWriter,
	r *http.Request) bool {
	err := r.ParseForm()
	if err != nil {
		returnAPIError(&s.config, w, err, http.StatusInternalServerError)
		return false
	}
	v, err := s.access.GetAllowed(r.FormValue("accessKey"))
	if v == false || err != nil {
		returnAPIError(&s.config, w,
			fmt.Errorf("Access denied"),
			http.StatusNetworkAuthenticationRequired)
		return false
	}
	return true
}
