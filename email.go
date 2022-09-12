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

// Email used to represent an email address.
type Email struct {
	Base
	Email string `json:"email"`
}

func NewEmail(s *owid.Signer, email string) (*Email, error) {
	var err error
	e := &Email{Email: email}
	e.Version = swanVersion
	e.OWID, err = s.CreateOWIDandSign(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func EmailFromJson(j []byte) (*Email, error) {
	var e Email
	err := json.Unmarshal(j, &e)
	if err != nil {
		return nil, err
	}
	e.OWID.Target = &e
	return &e, nil
}

func EmailUnmarshalBase64(value []byte) (*Email, error) {
	var e Email
	err := e.UnmarshalBase64(value)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (e *Email) UnmarshalBase64(value []byte) error {
	return unmarshalBase64(e, value)
}

func (e *Email) MarshalBase64() ([]byte, error) {
	return e.marshalBase64(e.marshal)
}

func (e *Email) MarshalOwid() ([]byte, error) {
	return e.marshalOwid(e.marshal)
}

func (e *Email) MarshalBinary() ([]byte, error) {
	return e.marshalBinary(e.marshal)
}

func (e *Email) marshal(b *bytes.Buffer) error {
	err := common.WriteString(b, e.Email)
	if err != nil {
		return err
	}
	return nil
}

func (e *Email) UnmarshalBinary(data []byte) error {
	return e.unmarshalBinary(e, data, func(b *bytes.Buffer) error {
		var err error
		e.Email, err = common.ReadString(b)
		if err != nil {
			return err
		}
		return nil
	})
}
