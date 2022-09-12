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
	"strconv"

	"github.com/SWAN-community/common-go"
	"github.com/SWAN-community/owid-go"
)

// Salt used to store the integer used as salt when hashing the email address
// to form the Signed in Id (SID).
type Salt struct {
	Base
	Salt uint32 `json:"salt"`
}

func (s *Salt) AsByteArray() ([]byte, error) {
	var b bytes.Buffer
	err := common.WriteUint32(&b, s.Salt)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func NewSaltFromString(s *owid.Signer, data string) (*Salt, error) {
	i, err := strconv.Atoi(data)
	if err != nil {
		return nil, err
	}
	return NewSalt(s, uint32(i))
}

func NewSalt(s *owid.Signer, data uint32) (*Salt, error) {
	var err error
	a := &Salt{Salt: data}
	a.Version = swanVersion
	a.OWID, err = s.CreateOWIDandSign(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func SaltFromJson(j []byte) (*Salt, error) {
	var a Salt
	err := json.Unmarshal(j, &a)
	if err != nil {
		return nil, err
	}
	a.OWID.Target = &a
	return &a, nil
}

func (a *Salt) ToBase64() (string, error) {
	b, err := a.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func SaltFromBase64(value string) (*Salt, error) {
	var a Salt
	err := unmarshalString(&a, value)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *Salt) MarshalOwid() ([]byte, error) {
	return a.marshalOwid(func(b *bytes.Buffer) error { return a.marshal(b) })
}

func (a *Salt) MarshalBinary() ([]byte, error) {
	return a.marshalBinary(func(b *bytes.Buffer) error { return a.marshal(b) })
}

func (a *Salt) marshal(b *bytes.Buffer) error {
	err := common.WriteUint32(b, a.Salt)
	if err != nil {
		return err
	}
	return nil
}

func (a *Salt) UnmarshalBinary(data []byte) error {
	return a.unmarshalBinary(a, data, func(b *bytes.Buffer) error {
		var err error
		a.Salt, err = common.ReadUint32(b)
		if err != nil {
			return err
		}
		return nil
	})
}
