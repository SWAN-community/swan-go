/* ****************************************************************************
 * Copyright 2022 51 Degrees Mobile Experts Limited (51degrees.com)
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
	"encoding/json"
	"net/http"

	"github.com/SWAN-community/common-go"
)

// Extension to Model with information needed in a request.
type ModelRequest struct {
	Model
}

// Verify the OWIDs in the model provided and handles any response to the
// caller. True is returned if the model is valid, otherwise false.
func (m *ModelRequest) Verify(w http.ResponseWriter, scheme string) bool {
	err := m.Model.Verify(scheme)
	if err != nil {
		common.ReturnApplicationError(w, &common.HttpError{
			Message: "invalid data",
			Error:   err,
			Code:    http.StatusBadRequest})
		return false
	}
	return true
}

// UnmarshalRequest populates the values of the model with those from the
// http request.
func (m *ModelRequest) UnmarshalRequest(r *http.Request) error {
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		return err
	}
	return nil
}
