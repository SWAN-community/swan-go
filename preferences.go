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
	"bytes"
	"encoding/json"

	"github.com/SWAN-community/common-go"
	"github.com/SWAN-community/owid-go"
)

// Preferences
// https://github.com/OneKey-Network/addressability-framework/blob/main/mvp-spec/model/preferences.md
type Preferences struct {
	Base
	Data PreferencesData `json:"data"`
}

func NewPreferences(s *owid.Signer, personalizedMarketing bool) (*Preferences, error) {
	var err error
	p := &Preferences{Data: PreferencesData{personalizedMarketing}}
	p.Version = swanVersion
	p.OWID, err = s.CreateOWIDandSign(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func PreferencesFromJson(j []byte) (*Preferences, error) {
	var p Preferences
	err := json.Unmarshal(j, &p)
	if err != nil {
		return nil, err
	}
	p.OWID.Target = &p
	return &p, nil
}

func PreferencesFromBase64(value string) (*Preferences, error) {
	var p Preferences
	err := unmarshalString(&p, value)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Preferences) ToBase64() (string, error) {
	return p.toBase64(func(b *bytes.Buffer) error { return p.marshal(b) })
}

func (p *Preferences) MarshalOwid() ([]byte, error) {
	return p.marshalOwid(func(b *bytes.Buffer) error { return p.marshal(b) })
}

func (p *Preferences) MarshalBinary() ([]byte, error) {
	return p.marshalBinary(func(b *bytes.Buffer) error { return p.marshal(b) })
}

func (p *Preferences) marshal(b *bytes.Buffer) error {
	err := common.WriteMarshaller(b, &p.Data)
	if err != nil {
		return err
	}
	return nil
}

func (p *Preferences) UnmarshalBinary(data []byte) error {
	return p.unmarshalBinary(p, data, func(b *bytes.Buffer) error {
		err := common.ReadMarshaller(b, &p.Data)
		if err != nil {
			return err
		}
		return nil
	})
}
