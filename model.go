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
	Value  []string
	Cookie *Cookie
}

type Entry struct {
	Cookie *Cookie
	OWID   *owid.OWID
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
	for _, v := range m.GetEntries() {
		ok, err := v.OWID.Verify(scheme)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("%s invalid", v.Cookie.Key)
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
	if s.Cookie == nil {
		s.Cookie = &Cookie{}
	}
	return s.Cookie.UnmarshalSwiftValidity(p)
}

func (m *ModelResponse) GetEntries() []*Entry {
	i := m.Model.GetEntries()
	if m.SID != nil {
		m.SID.GetCookie().Key = "sid"
		i = append(i, &Entry{
			OWID:   m.SID.GetOWID(),
			Cookie: m.SID.Cookie})
	}
	return i
}

func (m *Model) GetEntries() []*Entry {
	i := make([]*Entry, 0, 6)
	if m.Email != nil {
		m.Email.GetCookie().Key = "email"
		i = append(i, &Entry{
			OWID:   m.Email.GetOWID(),
			Cookie: m.Email.Cookie})
	}
	if m.Pref != nil {
		m.Pref.GetCookie().Key = "pref"
		i = append(i, &Entry{
			OWID:   m.Pref.GetOWID(),
			Cookie: m.Pref.Cookie})
	}
	if m.Salt != nil {
		m.Salt.GetCookie().Key = "salt"
		i = append(i, &Entry{
			OWID:   m.Salt.GetOWID(),
			Cookie: m.Salt.Cookie})
	}
	if m.RID != nil {
		m.RID.GetCookie().Key = "rid"
		i = append(i, &Entry{
			OWID:   m.RID.GetOWID(),
			Cookie: m.RID.Cookie})
	}
	if m.Stop != nil {
		i = append(i, &Entry{Cookie: m.Stop.Cookie})
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
		if v.Cookie.Expires.Before(m.Val.Expires) {
			m.Val.Expires = v.Cookie.Expires
		}
	}
	return nil
}
