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
	"encoding/base64"
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

func NewPreferences(s *owid.Signer, data bool) (*Preferences, error) {
	var err error
	p := &Preferences{Data: PreferencesData{data}}
	p.Base.Version = swanVersion
	p.Base.OWID, err = s.CreateOWIDandSign(p)
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
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	var p Preferences
	err = p.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Preferences) ToBase64() (string, error) {
	b, err := p.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (p *Preferences) marshal(b *bytes.Buffer) error {
	err := common.WriteByte(b, p.Base.Version)
	if err != nil {
		return err
	}
	err = common.WriteMarshaller(b, &p.Data)
	if err != nil {
		return err
	}
	return nil
}

func (p *Preferences) MarshalOwid() ([]byte, error) {
	var b bytes.Buffer
	err := p.marshal(&b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (p *Preferences) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	err := p.marshal(&b)
	if err != nil {
		return nil, err
	}
	err = p.Base.OWID.ToBuffer(&b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (p *Preferences) UnmarshalBinary(data []byte) error {
	var err error
	b := bytes.NewBuffer(data)
	p.Base.Version, err = common.ReadByte(b)
	if err != nil {
		return err
	}
	err = common.ReadMarshaller(b, &p.Data)
	if err != nil {
		return err
	}
	p.Base.OWID, err = owid.FromBuffer(b, p)
	if err != nil {
		return err
	}
	return nil
}
