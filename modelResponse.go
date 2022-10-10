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
	"time"

	"github.com/SWAN-community/swift-go"
)

// Extension to Model with information needed in a response.
type ModelResponse struct {
	Model
	SID *Identifier `json:"sid,omitempty"`
	Val Time        `json:"val,omitempty"` // Validity of the data
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
			m.RID.GetCookie().Key = "rid"
		case "email":
			m.Email = &Email{}
			err := m.Email.UnmarshalSwift(v)
			if err != nil {
				return err
			}
			m.Email.GetCookie().Key = "email"
		case "salt":
			m.Salt = &Salt{}
			err := m.Salt.UnmarshalSwift(v)
			if err != nil {
				return err
			}
			m.Salt.GetCookie().Key = "salt"
		case "pref":
			m.Pref = &Preferences{}
			err := m.Pref.UnmarshalSwift(v)
			if err != nil {
				return err
			}
			m.Pref.GetCookie().Key = "pref"
		case "sid":
			m.SID = &Identifier{}
			err := m.SID.UnmarshalSwift(v)
			if err != nil {
				return err
			}
			m.SID.GetCookie().Key = "sid"
		case "stop":
			m.Stop = &StringArray{}
			err := m.Stop.UnmarshalSwift(v)
			if err != nil {
				return err
			}
			m.Stop.getCookie().Key = "stop"
		case "state":
			for _, i := range v.Values() {
				m.State = append(m.State, string(i))
			}
		}
	}
	return nil
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

// SetValidity sets the created and expires times for all the cookies in the
// response model. This is used by the caller to indicate when they should
// recheck the returned data with the network for updates.
func (m *ModelResponse) SetValidity(revalidateSeconds int) error {
	m.Val.Created = time.Now().UTC()
	m.Val.Expires = m.Val.Created.Add(
		time.Duration(revalidateSeconds) * time.Second)

	// If there are any entries that will expire before the validation expires
	// time then the validation expires time needs to be adjust to the earliest
	// of the associated entries.
	for _, v := range m.GetEntries() {
		c := v.getCookie()
		if !c.Expires.IsZero() && c.Expires.Before(m.Val.Expires) {
			m.Val.Expires = v.getCookie().Expires
		}
	}
	return nil
}
