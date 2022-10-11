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
	"time"

	"github.com/SWAN-community/swift-go"
)

// Extension to Model with information needed in a response.
type ModelResponse struct {
	Model
	SID *Identifier `json:"sid,omitempty"`
	Ref LastRefresh `json:"val,omitempty"` // Validity of the data
}

// UnmarshalSwift populates the values of the model with the SWIFT results.
func (m *ModelResponse) UnmarshalSwift(r *swift.Results) error {

	// Set the fields that are also fields in the SWIFT results.
	m.State = r.State

	// Unpack or copy the SWIFT key value pairs that the model knows about and
	// which have at least one value. Use a scoped variable n as the instance
	// which will only be used in the model if no errors are found.
	for _, v := range r.Pairs {
		if len(v.Values()) > 0 {
			switch v.Key() {
			case "rid":
				n := &Identifier{}
				err := n.UnmarshalSwift(v)
				if err != nil {
					return fmt.Errorf("rid invalid: %w", err)
				}
				m.RID = n
				m.RID.GetCookie().Key = "rid"
			case "email":
				n := &Email{}
				err := n.UnmarshalSwift(v)
				if err != nil {
					return fmt.Errorf("email invalid: %w", err)
				}
				m.Email = n
				m.Email.GetCookie().Key = "email"
			case "salt":
				n := &Salt{}
				err := n.UnmarshalSwift(v)
				if err != nil {
					return fmt.Errorf("salt invalid: %w", err)
				}
				m.Salt = n
				m.Salt.GetCookie().Key = "salt"
			case "pref":
				n := &Preferences{}
				err := n.UnmarshalSwift(v)
				if err != nil {
					return fmt.Errorf("pref invalid: %w", err)
				}
				m.Pref = n
				m.Pref.GetCookie().Key = "pref"
			case "sid":
				n := &Identifier{}
				err := n.UnmarshalSwift(v)
				if err != nil {
					return fmt.Errorf("sid invalid: %w", err)
				}
				m.SID = n
				m.SID.GetCookie().Key = "sid"
			case "stop":
				n := &StringArray{}
				err := n.UnmarshalSwift(v)
				if err != nil {
					return err
				}
				m.Stop = n
				m.Stop.getCookie().Key = "stop"
			case "state":
				for _, i := range v.Values() {
					m.State = append(m.State, string(i))
				}
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
	m.Ref.Created = time.Now().UTC()
	m.Ref.Expires = m.Ref.Created.Add(
		time.Duration(revalidateSeconds) * time.Second)

	// If there are any entries that will expire before the validation expires
	// time then the validation expires time needs to be adjust to the earliest
	// of the associated entries.
	for _, v := range m.GetEntries() {
		c := v.getCookie()
		if !c.Expires.IsZero() && c.Expires.Before(m.Ref.Expires) {
			m.Ref.Expires = v.getCookie().Expires
		}
	}

	// Now set the cookie validity of all cookies to the validity period
	// expiration date. This ensures the browser will clear the cookies all at
	// the same time forcing a refresh.
	for _, v := range m.GetEntries() {
		v.getCookie().Expires = m.Ref.Expires
	}

	return nil
}
