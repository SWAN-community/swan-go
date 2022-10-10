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
	"fmt"
	"net/http"

	"github.com/SWAN-community/common-go"
)

// Model used when request or responding with SWAN data.
type Model struct {
	RID   *Identifier  `json:"rid,omitempty"`
	Pref  *Preferences `json:"pref,omitempty"`
	Email *Email       `json:"email,omitempty"`
	Salt  *Salt        `json:"salt,omitempty"`
	Stop  *StringArray `json:"stop,omitempty"`
	State []string     `json:"state,omitempty"`
}

// ModelRequestFromHttpRequest turns the request instance into a SWAN model.
// Any errors are responded to by the method. The caller can assume any problems
// have been dealt with if there is a nil reponse.
func ModelRequestFromHttpRequest(
	r *http.Request,
	w http.ResponseWriter) *ModelRequest {
	m := &ModelRequest{}
	err := m.UnmarshalRequest(r)
	if err != nil {
		common.ReturnApplicationError(w, &common.HttpError{
			Message: "bad data structure",
			Error:   err,
			Code:    http.StatusBadRequest})
		return nil
	}
	return m
}

// Verify confirms all the entries in the model have OWIDs that pass
// verification and have versions that are supported. The first one that does
// not pass will result in an error being returned. If no errors are returned
// then the model is fully verified.
func (m *Model) Verify(scheme string) error {
	for _, v := range m.GetEntries() {
		if v.GetOWID() != nil {
			ok, err := v.GetOWID().Verify(scheme)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("%s invalid", v.getCookie().Key)
			}
		}
		n := v.getVersion()
		if n < swanMinVersion || n > swanMaxVersion {
			return fmt.Errorf("version '%d' not supported", n)
		}
	}
	return nil
}

// GetEntries returns all the members of the model as an array.
func (m *Model) GetEntries() []Entry {
	i := make([]Entry, 0, 6)
	if m.Email != nil {
		m.Email.GetCookie().Key = "email"
		i = append(i, m.Email)
	}
	if m.Pref != nil {
		m.Pref.GetCookie().Key = "pref"
		i = append(i, m.Pref)
	}
	if m.Salt != nil {
		m.Salt.GetCookie().Key = "salt"
		i = append(i, m.Salt)
	}
	if m.RID != nil {
		m.RID.GetCookie().Key = "rid"
		i = append(i, m.RID)
	}
	if m.Stop != nil {
		i = append(i, m.Stop)
	}
	return i
}
