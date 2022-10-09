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
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SWAN-community/common-go"
	"github.com/SWAN-community/swift-go"
)

type Entry interface {
	Verifiable

	// Returns the cookie representation of the SWAN entity.
	getCookie() *Cookie

	// Returns the version of the SWAN entity.
	getVersion() byte
}

// Model used when request or responding with SWAN data.
type Model struct {
	RID   *Identifier  `json:"rid,omitempty"`
	Pref  *Preferences `json:"pref,omitempty"`
	Email *Email       `json:"email,omitempty"`
	Salt  *Salt        `json:"salt,omitempty"`
	Stop  *StringArray `json:"stop,omitempty"`
	State []string     `json:"state,omitempty"`
}

// Extension to Model with information needed in a response.
type ModelResponse struct {
	Model
	SID *Identifier `json:"sid,omitempty"`
	Val Cookie      `json:"val,omitempty"`
}

// Extension to Model with information needed in a request.
type ModelRequest struct {
	Model
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

// UnmarshalSwift populates the values of the model with the SWIFT results.
func (m *ModelResponse) UnmarshalSwift(r *swift.Results) error {

	// Set the fields that are also fields in the SWIFT results.
	m.State = r.State

	// Unpack or copy the SWIFT key value pairs that the model knows about.
	for _, v := range r.Pairs {
		switch v.Key() {
		case "rid":
			m.RID = &Identifier{}
			err := m.RID.UnmarshalSwift(v)
			if err != nil {
				return err
			}
		case "email":
			m.Email = &Email{}
			err := m.Email.UnmarshalSwift(v)
			if err != nil {
				return err
			}
		case "salt":
			m.Salt = &Salt{}
			err := m.Salt.UnmarshalSwift(v)
			if err != nil {
				return err
			}
		case "pref":
			m.Pref = &Preferences{}
			err := m.Pref.UnmarshalSwift(v)
			if err != nil {
				return err
			}
		case "sid":
			m.SID = &Identifier{}
			err := m.SID.UnmarshalSwift(v)
			if err != nil {
				return err
			}
		case "stop":
			m.Stop = &StringArray{}
			err := m.Stop.UnmarshalSwift(v)
			if err != nil {
				return err
			}
		case "state":
			for _, i := range v.Values() {
				m.State = append(m.State, string(i))
			}
		}
	}
	return nil
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

// UnmarshalSwift to turn a SWIFT pair into a string array.
func (s *StringArray) UnmarshalSwift(p *swift.Pair) error {
	s.Value = make([]string, 0, len(p.Values()))
	for _, v := range p.Values() {
		if len(v) > 0 {
			s.Value = append(s.Value, string(v))
		}
	}
	if s.Cookie == nil {
		s.Cookie = &Cookie{}
	}
	return s.Cookie.UnmarshalSwiftValidity(p)
}

// GetEntries returns all the members of the model response as an array.
func (m *ModelResponse) GetEntries() []Entry {
	i := m.Model.GetEntries()
	if m.SID != nil && m.SID.OWID != nil {
		m.SID.GetCookie().Key = "sid"
		i = append(i, m.SID)
	}
	return i
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

// SetValidity sets the created and expires times. This is used by the caller to
// indicate when they should recheck the returned data with the SWAN network
// for updates.
func (m *ModelResponse) SetValidity(revalidateSeconds int) error {
	m.Val.Created = time.Now().UTC()
	m.Val.Expires = m.Val.Created.Add(
		time.Duration(revalidateSeconds) * time.Second)
	for _, v := range m.GetEntries() {
		if v.getCookie().Expires.Before(m.Val.Expires) {
			m.Val.Expires = v.getCookie().Expires
		}
	}
	return nil
}
