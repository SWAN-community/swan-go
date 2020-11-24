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

		err := r.ParseForm()
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusUnprocessableEntity)
			return
		}

		// TODO : Verify that the OWIDs are valid before creating the offer ID.
		// Do this using the call to the domain in the OWID and the standardized
		// path. This is because the OWIDs could have come from anywhere.
		pl := r.FormValue("placement")
		if pl == "" {
			returnAPIError(&s.config, w,
				errors.New("missing placement parameter"),
				http.StatusUnprocessableEntity)
			return
		}

		pu := r.FormValue("pubdomain")
		if pu == "" {
			returnAPIError(&s.config, w,
				errors.New("missing pubdomain parameter"),
				http.StatusBadRequest)
			return
		}

		cbid, err := owid.DecodeFromBase64(r.FormValue("cbid"))
		if cbid == nil || err != nil {
			returnAPIError(&s.config, w,
				errors.New("missing cbid parameter"),
				http.StatusBadRequest)
			return
		}

		p, err := owid.DecodeFromBase64(r.FormValue("preferences"))
		if p == nil || err != nil {
			returnAPIError(&s.config, w,
				errors.New("missing preferences parameter"),
				http.StatusBadRequest)
			return
		}

		uuid, err := uuid.New().MarshalBinary()
		if err != nil {
			returnServerError(&s.config, w, err)
			return
		}

		oid := OfferID{
			pl,
			pu,
			uuid,
			cbid.PayloadAsString(),
			p.PayloadAsString()}

		os, err := oid.AsByteArray()
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusUnprocessableEntity)
			return
		}

		owid, err := encodeAsOWID(s, r, os)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusUnprocessableEntity)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Write([]byte(owid))
	}
}
