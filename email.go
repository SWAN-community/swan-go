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

// Email used to represent an email address.
type Email struct {
	Base
	email string
}

func NewEmail(s *owid.Signer, email string) (*Email, error) {
	var err error
	e := &Email{email: email}
	e.Base.Version = swanVersion
	e.Base.OWID, err = s.CreateOWIDandSign(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func EmailFromBase64(value string) (*Email, error) {
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	var e Email
	err = e.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (e *Email) ToBase64() (string, error) {
	b, err := e.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (e *Email) Email() string { return e.email }

func (e *Email) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["version"] = e.Base.Version
	m["email"] = e.email
	m["source"] = e.Base.OWID
	return json.Marshal(m)
}

func (e *Email) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	if v, ok := m["version"].(float64); ok {
		e.Base.Version = byte(v)
	} else {
		return errorMissing("version")
	}
	if v, ok := m["email"].(string); ok {
		e.email = v
	} else {
		return errorMissing("email")
	}
	if o, ok := m["source"].(owid.OWID); ok {
		e.Base.OWID = &o
		o.Target = e
	} else {
		return errorMissing("source")
	}
	return nil
}

func (e *Email) marshal(b *bytes.Buffer) error {
	err := common.WriteByte(b, e.Base.Version)
	if err != nil {
		return err
	}
	err = common.WriteString(b, e.email)
	if err != nil {
		return err
	}
	return nil
}

func (e *Email) MarshalOwid() ([]byte, error) {
	var b bytes.Buffer
	err := e.marshal(&b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (e *Email) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	err := e.marshal(&b)
	if err != nil {
		return nil, err
	}
	err = e.Base.OWID.ToBuffer(&b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (e *Email) UnmarshalBinary(data []byte) error {
	var err error
	b := bytes.NewBuffer(data)
	e.Base.Version, err = common.ReadByte(b)
	if err != nil {
		return err
	}
	e.email, err = common.ReadString(b)
	if err != nil {
		return err
	}
	e.Base.OWID, err = owid.FromBuffer(b, e)
	if err != nil {
		return err
	}
	return nil
}
