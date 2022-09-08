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

// Salt used to represent an Salt address.
type Salt struct {
	Base
	Salt []byte `json:"salt"`
}

func NewSalt(s *owid.Signer, salt []byte) (*Salt, error) {
	var err error
	a := &Salt{Salt: salt}
	a.Base.OWID, err = s.CreateOWIDandSign(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func SaltFromJson(j []byte) (*Salt, error) {
	var s Salt
	err := json.Unmarshal(j, &s)
	if err != nil {
		return nil, err
	}
	s.OWID.Target = &s
	return &s, nil
}

func SaltFromBase64(value string) (*Salt, error) {
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	var s Salt
	err = s.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Salt) ToBase64() (string, error) {
	b, err := s.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (s *Salt) marshal(b *bytes.Buffer) error {
	err := common.WriteByte(b, s.Base.Version)
	if err != nil {
		return err
	}
	err = common.WriteByteArray(b, s.Salt)
	if err != nil {
		return err
	}
	return nil
}

func (s *Salt) MarshalOwid() ([]byte, error) {
	var b bytes.Buffer
	err := s.marshal(&b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (s *Salt) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	err := s.marshal(&b)
	if err != nil {
		return nil, err
	}
	err = s.Base.OWID.ToBuffer(&b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (s *Salt) UnmarshalBinary(data []byte) error {
	var err error
	b := bytes.NewBuffer(data)
	s.Base.Version, err = common.ReadByte(b)
	if err != nil {
		return err
	}
	s.Salt, err = common.ReadByteArray(b)
	if err != nil {
		return err
	}
	s.Base.OWID, err = owid.FromBuffer(b, s)
	if err != nil {
		return err
	}
	return nil
}
