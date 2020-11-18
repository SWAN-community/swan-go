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
)

// Take the incoming parameters that map to the OfferID structure to create the
// OfferID. Then turn the OfferID into a byte array to be used as the payload
// for the OWID that is returned as a string.
func handlerCreateOfferID(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusUnprocessableEntity)
			return
		}

		pl := r.FormValue("placement")
		if pl == "" {
			returnAPIError(&s.config, w, errors.New("missing placement paramter"), http.StatusBadRequest)
			return
		}
		pu := r.FormValue("pubdomain")
		if pu == "" {
			returnAPIError(&s.config, w, errors.New("missing pubdomain paramter"), http.StatusBadRequest)
			return
		}
		c := r.FormValue("cbid")
		if c == "" {
			returnAPIError(&s.config, w, errors.New("missing cbid paramter"), http.StatusBadRequest)
			return
		}
		si := r.FormValue("sid")
		if si == "" {
			returnAPIError(&s.config, w, errors.New("missing sid paramter"), http.StatusBadRequest)
			return
		}
		p := r.FormValue("preferences")
		if p == "" {
			returnAPIError(&s.config, w, errors.New("missing preferences paramter"), http.StatusBadRequest)
			return
		}

		o := OfferID{pl,
			pu,
			c,
			si,
			p}

		os, err := o.AsString()
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
