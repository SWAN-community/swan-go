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
	"net/http"
)

func handlerCreateSWID(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check caller is authorized to access SWAN.
		if s.getAccessAllowed(w, r) == false {
			return
		}

		// Create the SWID OWID for this SWAN Operator.
		c, err := createSWID(s, r)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
		}

		// Get the OWID as a byte array.
		b, err := c.AsByteArray()
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
		}

		// Return the URL as a byte array.
		sendResponse(s, w, "application/octet-stream", b)
	}
}
