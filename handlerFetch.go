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
	"log"
	"net/http"
	"owid"
	"swift"
	"time"
)

// handlerFetch returns a URL that can be used in the browser primary navigation
// to retrieve the most current data from the SWAN network. If no data is
// available default values are returned.
func handlerFetch(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		// Check caller is authorized to access SWAN.
		if s.getAccessAllowed(w, r) == false {
			return
		}

		// Write out the input to the log if in debug mode.
		if s.config.Debug {
			log.Println(r.URL.String() + "?" + r.Form.Encode())
		}

		// Validate the set the return URL.
		err = swift.SetURL("returnUrl", "returnUrl", &r.Form)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// If the request includes data that is currently held by the caller
		// then configure the storage operation to use these values if they
		// relate to valid OWIDs.
		setDefaults(s, r)

		// Uses the SWIFT access node associated with this internet domain
		// to determine the URL to direct the browser to.
		u, err := createStorageOperationURL(s.swift, r, r.Form)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Write out the URL to the log if in debug mode.
		if s.config.Debug {
			log.Println(u)
		}

		// Return the response from the SWIFT layer.
		sendResponse(s, w, "text/plain; charset=utf-8", []byte(u))
	}
}

// setDefaults sets the values for the storage operation in SWIFT if there are
// no values in the network. SWID, preference OWIDs, and stop identifiers can be
// provided by the caller for this situation. If no SWID is provided then SWAN
// will assign a new random one.
func setDefaults(s *services, r *http.Request) {
	t := s.config.DeleteDate()
	q := &r.Form

	// Process any exist SWID, preference or stop data provided by the caller.
	setSWID(s, r, t)
	setPerf(s, r, t)
	setStop(s, r, t)

	// Get the email address either to return as the raw value, or to turn into
	// a SID once it's been fetched. Always favour the most recent email address
	// available across the network.
	q.Set("email>", "")

	// Delete any common parameters that might have been included in the request
	// that we do not need. Avoids SWIFT trying to process then as keys.
	q.Del("sid")
	q.Del("val")
}

// setStop uses the values provided and will add them to any other stop values
// already contained in the network.
func setStop(s *services, r *http.Request, t time.Time) {
	if r.Form.Get("stop") != "" {
		r.Form.Set(
			fmt.Sprintf("stop+%s", t.Format("2006-01-02")),
			r.Form.Get("stop"))
	} else {
		r.Form.Set("stop+", "")
	}
	r.Form.Del("stop")
}

// setPerf gets the value of perf from the request and verifies it's a valid
// OWID. If it is valid then use it as the default value if the SWAN network
// does not contain a value. If it is not valid then an empty value will be
// used to indicate that the user has not provided any preferences.
func setPerf(s *services, r *http.Request, t time.Time) {
	v := r.Form.Get("pref") // The value for the Perf. to use if one not found
	o, err := owid.FromBase64(v)
	if err != nil {
		logNonCriticalError(s, err)
		v = ""
	} else {

		// There is a valid OWID for the Perf. Does it meet the rules?
		b, err := o.Verify(s.config.Scheme)
		if err != nil {
			logNonCriticalError(s, err)
			v = ""
		} else if b {

			// Change the expiry time to be based on the Perf. creation date.
			t = o.Date.AddDate(0, 0, s.config.DeleteDays)

			// If the value has already expired then don't use it. If not then
			// use it as the value if the network does not currently contain a
			// value.
			if time.Now().UTC().After(t) {
				v = ""
			}
		} else {
			v = ""
		}
	}

	// Set the value in the SWIFT storage operation, and remove the perf from
	// the form.
	if v != "" {

		// There is an existing preference stored by the caller. Use this value
		// if the network does not currently contain a more recent version.
		r.Form.Set(fmt.Sprintf("pref>%s", t.Format("2006-01-02")), v)

	} else {

		// There is no existing preference available. Therefore retrieve the
		// newest value contained in the network.
		r.Form.Set("pref>", "")
	}

	// Remove pref key as this is not valid for a SWIFT operation.
	r.Form.Del("pref")
}

// setSWID gets the value of the SWID from the form associated with the request.
// If that SWID is a valid OWID, can be verified with the creators public key,
// and is SWAN access node that is known to this access node, and is finally
// still valid when checked against the delete date, then use this value in
// cases where the SWID does not exist in the SWAN network. This might be
// because the SWAN Operators nodes have had cookies removed due to tracking
// prevention methods, but the value that the caller has is still valid and can
// be used by the SWAN Operators.
// If none of the conditions are valid then a new SWID is created and used if
// the SWAN network does not contain any other values.
func setSWID(s *services, r *http.Request, t time.Time) {
	v := r.Form.Get("swid") // The value for the SWID to use if one not found
	o, err := owid.FromBase64(v)
	if err != nil {
		logNonCriticalError(s, err)
		v = ""
	} else {

		// There is a valid OWID for the SWID. Does it meet the rules?
		b, err := o.Verify(s.config.Scheme)
		if err != nil {
			logNonCriticalError(s, err)
			v = ""
		} else if b && isSWAN(s, o) {

			// Change the expiry time to be based on the SWID creation date.
			t = o.Date.AddDate(0, 0, s.config.DeleteDays)

			// If the value has already expired then don't use it. If not then
			// use it as the value if the network does not currently contain a
			// value.
			if time.Now().UTC().After(t) {
				v = ""
			}
		} else {
			v = ""
		}
	}

	// Set the value in the SWIFT storage operation, and remove the SWID from
	// the form.
	if v != "" {

		// There is an existing SWID stored by the caller. Use this value if the
		// network does not currently contain a more recent version.
		r.Form.Set(fmt.Sprintf("swid>%s", t.Format("2006-01-02")), v)

	} else {

		// There is no existing SWID available. Therefore retrieve the newest
		// value contained in the network.
		r.Form.Set("swid>", "")
	}

	// Remove swid key as this is not valid for a SWIFT operation.
	r.Form.Del("swid")
}

// isSWAN returns true if the OWID was created from an access node known to this
// SWAN access node.
func isSWAN(s *services, o *owid.OWID) bool {
	n, err := s.swift.GetAccessNodeForHost(o.Domain)
	if err != nil {
		logNonCriticalError(s, err)
		return false
	}
	return n != nil && n.Domain() == o.Domain
}

func logNonCriticalError(s *services, err error) {
	if s.config.Debug {
		log.Println(err)
	}
}
