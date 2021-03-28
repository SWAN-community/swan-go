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
	"swift"
)

// The methods defined here assume that the internet domain associated with the
// SWAN Operator has also been registered as a SWIFT access node in the SWIFT
// services. If this is not the case then an error is returned.

// createStorageOperationURL returns a URL to redirect the web browser to that
// will perform the storage operation requested.
// s instance of the SWIFT services that are mapped to this SWAN Operator
// r the current HTTP request from the web browser
// q key value pairs to include in the query, updated by the method
func createStorageOperationURL(
	s *swift.Services,
	r *http.Request,
	q url.Values) (string, error) {

	// Add the HTTP headers that will impact the home node calculation.
	swift.SetHomeNodeHeaders(r, &q)

	// Set the table to SWAN overriding any current value.
	q.Set("table", "swan")

	return swift.Create(s, r.Host, q)
}

// decrypt uses the SWIFT service for this access node to decrypt the data as
// a byte array.
// s instance of the SWIFT services that are mapped to this SWAN Operator
// h the internet domain associated with the access node
// d byte array to be decrypted and decoded
func decryptAndDecode(
	s *swift.Services,
	h string,
	d []byte) (*swift.Results, error) {

	// Get the node associated with the request.
	n, err := s.GetAccessNodeForHost(h)
	if err != nil {
		return nil, err
	}
	return n.DecryptAndDecode(d)
}

// setDefaults sets the empty values for the storage operation in SWIFT. If
// values exist then these are used rather than the defaults. If not CBID exists
// the a new random value is used.
func setDefaults(s *services, r *http.Request) error {
	t := s.config.DeleteDate().Format("2006-01-02")
	q := &r.Form
	c, err := createCBID(s, r)
	if err != nil {
		return err
	}
	q.Set(fmt.Sprintf("cbid<%s", t), c.AsString())
	q.Set(fmt.Sprintf("email<%s", t), "")
	q.Set(fmt.Sprintf("allow<%s", t), "")
	q.Set(fmt.Sprintf("stop<%s", t), "")
	return nil
}
