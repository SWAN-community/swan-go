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

func EmptyUnmarshalBase64(value []byte) (*Empty, error) {
	var a Empty
	err := unmarshalBase64(&a, value)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *Empty) MarshalBase64() ([]byte, error) {
	return a.marshalBase64(a.marshal)
}

func (a *Empty) MarshalOwid() ([]byte, error) {
	return a.marshalOwid(a.marshal)
}

func (a *Empty) MarshalBinary() ([]byte, error) {
	return a.marshalBinary(a.marshal)
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
