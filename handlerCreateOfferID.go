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
	"owid"

	"github.com/google/uuid"
)

// Take the incoming parameters that map to the OfferID structure to create the
// OfferID. Then turn the OfferID into a byte array to be used as the payload
// for the OWID that is returned as a string.
func handlerCreateOfferID(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check caller can access
		if s.getAccessAllowed(w, r) == false {
			returnAPIError(&s.config, w,
				errors.New("Not authorized"),
				http.StatusUnauthorized)
			return
		}

		// Get the creator associated with this SWAN domain.
		c, err := s.owid.GetCreator(r.Host)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}
		if c == nil {
			err = fmt.Errorf(
				"No creator for '%s'. Use http[s]://%s/owid/register to setup "+
					"domain.",
				r.Host,
				r.Host)
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}

		// Create the offer ID.
		o, err := createOfferID(s, r, c)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
		}

		// Return the Offer ID as a byte array.
		b, err := o.TreeAsByteArray()
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Write(b)
	}
}

func createOfferID(
	s *services,
	r *http.Request,
	c *owid.Creator) (*owid.OWID, error) {
	of, err := getOfferID(s, r)
	if err != nil {
		return nil, err
	}
	b, err := of.AsByteArray()
	if err != nil {
		return nil, err
	}
	o := c.CreateOWID(b)
	err = c.Sign(o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func getValue(r *http.Request, k string) (string, error) {
	v := r.FormValue(k)
	if v == "" {
		return "", fmt.Errorf("missing '%s' parameter", k)
	}
	return v, nil
}

func getOWID(s *services, r *http.Request, k string) (*owid.OWID, error) {
	v, err := getValue(r, k)
	if err != nil {
		return nil, err
	}
	o, err := owid.TreeFromBase64(v)
	if err != nil {
		return nil, err
	}
	e, err := o.Verify(s.config.Scheme)
	if err != nil {
		return nil, err
	}
	if e == false {
		return nil, fmt.Errorf("'%s' not a valid OWID", k)
	}
	return o, nil
}

func getOfferID(s *services, r *http.Request) (*Offer, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	pl, err := getValue(r, "placement")
	if pl == "" {
		return nil, err
	}

	pu, err := getValue(r, "pubdomain")
	if pu == "" {
		return nil, err
	}

	cbid, err := getOWID(s, r, "cbid")
	if err != nil {
		return nil, err
	}

	sid, err := getOWID(s, r, "sid")
	if err != nil {
		return nil, err
	}

	pref, err := getOWID(s, r, "preferences")
	if err != nil {
		return nil, err
	}

	// Random one time data is used to ensure the Offer ID is unique for all
	// time.
	uuid, err := uuid.New().MarshalBinary()
	if err != nil {
		return nil, err
	}

	// Create the offer byte array.
	return &Offer{
		pl,
		pu,
		uuid,
		cbid.Payload,
		sid.Payload,
		pref.Payload}, nil
}
