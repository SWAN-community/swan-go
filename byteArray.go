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

// ByteArray used for general purpose data storage.
type ByteArray struct {
	Base
	Data []byte `json:"data"`
}

func NewByteArray(s *owid.Signer, data []byte) (*ByteArray, error) {
	var err error
	a := &ByteArray{Data: data}
	a.Version = swanVersion
	a.OWID, err = s.CreateOWIDandSign(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func ByteArrayFromJson(j []byte) (*ByteArray, error) {
	var a ByteArray
	err := json.Unmarshal(j, &a)
	if err != nil {
		return nil, err
	}
	a.OWID.Target = &a
	return &a, nil
}

func ByteArrayUnmarshalBase64(value []byte) (*ByteArray, error) {
	var a ByteArray
	err := a.UnmarshalBase64(value)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *ByteArray) UnmarshalBase64(value []byte) error {
	return unmarshalBase64(a, value)
}

func (a *ByteArray) MarshalBase64() ([]byte, error) {
	return a.marshalBase64(a.marshal)
}

func (a *ByteArray) MarshalOwid() ([]byte, error) {
	return a.marshalOwid(a.marshal)
}

func (a *ByteArray) MarshalBinary() ([]byte, error) {
	return a.marshalBinary(a.marshal)
}

func (a *ByteArray) marshal(b *bytes.Buffer) error {
	err := common.WriteByteArray(b, a.Data)
	if err != nil {
		return err
	}
	return nil
}

func (a *ByteArray) UnmarshalBinary(data []byte) error {
	return a.unmarshalBinary(a, data, func(b *bytes.Buffer) error {
		var err error
		a.Data, err = common.ReadByteArray(b)
		if err != nil {
			return err
		}
		return nil
	})
}
