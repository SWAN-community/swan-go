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
	"fmt"

	"github.com/SWAN-community/owid-go"
)

// Empty contains nothing. Used for most OWIDs that just sign the root and
// themselves.
type Empty struct {
	Response
}

func NewEmpty(s *owid.Signer) (*Empty, error) {
	var err error
	a := &Empty{}
	a.Version = swanVersion
	a.StructType = responseEmpty
	a.OWID, err = s.CreateOWIDandSign(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func EmptyFromJson(j []byte) (*Empty, error) {
	var a Empty
	err := json.Unmarshal(j, &a)
	if err != nil {
		return nil, err
	}
	a.OWID.Target = &a
	return &a, nil
}

func (a *Empty) ToBase64() (string, error) {
	b, err := a.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func EmptyFromBase64(value string) (*Empty, error) {
	var a Empty
	err := unmarshalString(&a, value)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *Empty) MarshalOwid() ([]byte, error) {
	return a.marshalOwid(func(b *bytes.Buffer) error { return a.marshal(b) })
}

func (a *Empty) MarshalBinary() ([]byte, error) {
	return a.marshalBinary(func(b *bytes.Buffer) error { return a.marshal(b) })
}

func (a *Empty) marshal(b *bytes.Buffer) error {
	return nil
}

func (a *Empty) UnmarshalBinary(data []byte) error {
	return a.unmarshalBinary(a, data, func(b *bytes.Buffer) error {
		if a.StructType != responseEmpty {
			return fmt.Errorf("struct type not failed '%d'", responseEmpty)
		}
		return nil
	})
}
