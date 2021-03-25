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
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// handlerOperationAsJSON returns the original data held in the the operation.
// Used by user interfaces to get the operations details for dispaly, or to
// continue a storage operation after time has passed waiting for the user.
// This method should never be used for passing for purposes other than for
// users editing their data.
func handlerOperationAsJSON(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check caller can access.
		if s.getAccessAllowed(w, r) == false {
			returnAPIError(&s.config, w,
				errors.New("Not authorized"),
				http.StatusUnauthorized)
			return
		}

		// Decode the query string to form the byte array.
		d, err := base64.RawURLEncoding.DecodeString(r.Form.Get("data"))
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Decrypt the string with the access node.
		v, err := decryptAndDecode(s.swift, r.Host, d)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Modify the expiry time.
		for _, i := range v.Values {
			i.Expires = time.Now().UTC().Add(time.Second * s.config.ValueTimeout)
		}

		// Turn the Results into a JSON string.
		j, err := json.Marshal(v)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}

		// Send the JSON string.
		sendGzipJSON(s, w, j)
	}
}

// handlerValuesAsJSON turns the the "data" parameter into an array of key value
// pairs where the value is encoded as an OWID using the credentials of the SWAN
// Operator.
// If the timestamp of the data provided has expired then an error is returned.
// The Email value is converted to a hashed version before being returned.
func handlerValuesAsJSON(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check caller can access.
		if s.getAccessAllowed(w, r) == false {
			returnAPIError(&s.config, w,
				errors.New("Not authorized"),
				http.StatusUnauthorized)
			return
		}

		// Decode the query string to form the byte array.
		d, err := base64.RawURLEncoding.DecodeString(r.Form.Get("data"))
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Decrypt the string with the access node.
		v, err := decryptAndDecode(s.swift, r.Host, d)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Validate that the timestamp has not expired.
		if v.IsTimeStampValid() == false {
			returnAPIError(
				&s.config,
				w,
				fmt.Errorf("data expired"),
				http.StatusBadRequest)
			return
		}

		// Change the values to OWIDs.
		for _, p := range v.Values {
			if p.Key == "email" {
				p.Key = "sid"
				p.Value, err = encodeAsOWID(s, r, createSID(p.Value))
			} else {
				p.Value, err = encodeAsOWID(s, r, []byte(p.Value))
			}
			if err != nil {
				returnAPIError(&s.config, w, err, http.StatusInternalServerError)
				return
			}
		}

		// Modify the expiry time.
		for _, i := range v.Values {
			i.Expires = time.Now().UTC().Add(time.Second * s.config.ValueTimeout)
		}

		// Turn the Results into a JSON string.
		j, err := json.Marshal(v.Values)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}

		// Send the JSON string.
		sendGzipJSON(s, w, j)
	}
}

func sendGzipJSON(s *services, w http.ResponseWriter, j []byte) {
	g := gzip.NewWriter(w)
	defer g.Close()
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	_, err := g.Write(j)
	if err != nil {
		returnAPIError(&s.config, w, err, http.StatusInternalServerError)
	}
}

func encodeAsOWID(s *services, r *http.Request, v []byte) (string, error) {

	// Get the creator associated with this SWAN domain.
	c, err := s.owid.GetCreator(r.Host)
	if err != nil {
		return "", err
	}
	if c == nil {
		return "", fmt.Errorf(
			"No creator for '%s'. Use http[s]://%s/owid/register to setup "+
				"domain.",
			r.Host,
			r.Host)
	}

	// Create and sign the OWID.
	o := c.CreateOWID(v)
	err = c.Sign(o)
	if err != nil {
		return "", err
	}

	return o.AsBase64()
}

// TODO : What hashing algorithm do we want to use to turn email address into
// hashes?
func createSID(s string) []byte {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	return hasher.Sum(nil)
}
