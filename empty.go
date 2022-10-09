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
	"fmt"

	"github.com/SWAN-community/owid-go"
)

// Empty contains nothing. Used for most OWIDs that just sign the root and
// themselves.
type Empty struct {
	Response
}

// Returns an OWID with the target populated.
func (e *Empty) GetOWID() *owid.OWID {
	if e.OWID == nil {
		e.OWID = &owid.OWID{}
	}
	if e.OWID.Target == nil {
		e.OWID.Target = e
	}
	return e.OWID
}

func NewEmpty(signer *owid.Signer, seed *Seed) (*Empty, error) {
	var err error
	e := &Empty{}
	e.Version = swanVersion
	e.StructType = responseEmpty
	e.Seed = seed
	e.OWID, err = signer.CreateOWIDandSign(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func EmptyUnmarshalBase64(value []byte) (*Empty, error) {
	var e Empty
	err := unmarshalBase64(&e, value)
	if err != nil {
		return nil, err
	}
	e.OWID.Target = &e
	return &e, nil
}

func (e *Empty) MarshalBase64() ([]byte, error) {
	return e.Response.marshalBase64(e.marshal)
}

func (e *Empty) MarshalOwid() ([]byte, error) {
	return e.Response.marshalOwid(e.marshal)
}

func (e *Empty) MarshalBinary() ([]byte, error) {
	return e.Response.marshalBinary(e.marshal)
}

func (e *Empty) UnmarshalBinary(data []byte) error {
	return e.Response.unmarshalBinary(e, data, func(b *bytes.Buffer) error {
		if e.StructType != responseEmpty {
			return fmt.Errorf("struct type not failed '%d'", responseEmpty)
		}
		return nil
	})
}

func (e *Empty) marshal(b *bytes.Buffer) error {
	return nil
}
