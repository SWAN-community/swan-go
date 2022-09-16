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

	"github.com/SWAN-community/owid-go"
	"github.com/SWAN-community/swift-go"
)

type StringArray struct {
	Value []string
	Cookie
}

type entry struct {
	key      string
	validity *Cookie
	owid     *owid.OWID
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
	SID *ByteArray `json:"sid,omitempty"`
	Val Cookie     `json:"val,omitempty"`
}

// Extension to Model with information needed in a request.
type ModelRequest struct {
	Model
}

func (m *ModelRequest) UnmarshalRequest(r *http.Request) error {
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		return err
	}
	return nil
}

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
			m.SID = &ByteArray{}
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

func (m *Model) Verify(scheme string) error {
	for _, v := range m.getEntries() {
		ok, err := v.owid.Verify(scheme)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("%s invalid", v.key)
		}
	}
	return nil
}

func (s *StringArray) UnmarshalSwift(p *swift.Pair) error {
	s.Value = make([]string, 0, len(p.Values()))
	for _, v := range p.Values() {
		if len(v) > 0 {
			s.Value = append(s.Value, string(v))
		}
	}
	return s.UnmarshalSwiftValidity(p)
}

func (m *ModelResponse) getEntries() []*entry {
	i := m.Model.getEntries()
	if m.SID != nil {
		i = append(i, &entry{
			key:      "sid",
			owid:     m.SID.GetOWID(),
			validity: m.SID.Cookie})
	}
	return i
}

func (m *Model) getEntries() []*entry {
	i := make([]*entry, 0, 5)
	if m.Email != nil {
		i = append(i, &entry{
			key:      "email",
			owid:     m.Email.GetOWID(),
			validity: m.Email.Cookie})
	}
	if m.Pref != nil {
		i = append(i, &entry{
			key:      "pref",
			owid:     m.Pref.GetOWID(),
			validity: m.Pref.Cookie})
	}
	if m.Salt != nil {
		i = append(i, &entry{
			key:      "salt",
			owid:     m.Salt.GetOWID(),
			validity: m.Salt.Cookie})
	}
	if m.RID != nil {
		i = append(i, &entry{
			key:      "rid",
			owid:     m.RID.GetOWID(),
			validity: m.RID.Cookie})
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
	for _, v := range m.getEntries() {
		if v.validity.Expires.Before(m.Val.Expires) {
			m.Val.Expires = v.validity.Expires
		}
	}
	return nil
}
