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
	"owid"

	"github.com/google/uuid"
)

// createCBID returns an OWID with a unique new CBID.
func createCBID(s *services, r *http.Request) (*owid.OWID, error) {
	c, err := s.owid.GetCreator(r.Host)
	if err != nil {
		return nil, err
	}
	u, err := uuid.New().MarshalBinary()
	if err != nil {
		return nil, err
	}
	o, err := c.CreateOWIDandSign(u)
	if err != nil {
		return nil, err
	}
	return o, err
}
