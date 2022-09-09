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
	"encoding"
	"encoding/base64"

	"github.com/SWAN-community/common-go"
	"github.com/SWAN-community/owid-go"
)

// Base used with any SWAN field.
type Base struct {
	Version byte       `json:"version"` // Used to indicate the version encoding of the type.
	OWID    *owid.OWID `json:"source"`  // OWID related to the structure
}

type base interface {
	marshal(*bytes.Buffer) error
}

// writeData writes the version before calling the function.
func (b *Base) writeData(u *bytes.Buffer, f func(*bytes.Buffer) error) error {
	err := common.WriteByte(u, b.Version)
	if err != nil {
		return err
	}
	err = f(u)
	if err != nil {
		return err
	}
	return nil
}

// marshalOwid returns a byte array of all the data needed by an OWID.
func (b *Base) marshalOwid(f func(*bytes.Buffer) error) ([]byte, error) {
	var u bytes.Buffer
	err := b.writeData(&u, f)
	if err != nil {
		return nil, err
	}
	return u.Bytes(), nil
}

// marshalBinary marshals the version, calls the function to add more data, and
// finishes by adding the OWID before returning the byte array.
func (b *Base) marshalBinary(f func(*bytes.Buffer) error) ([]byte, error) {
	var u bytes.Buffer
	err := b.writeData(&u, f)
	if err != nil {
		return nil, err
	}
	err = b.OWID.ToBuffer(&u)
	if err != nil {
		return nil, err
	}
	return u.Bytes(), nil
}

// toBase64 encodes the marshalBinary result as a base64 string.
func (b *Base) toBase64(f func(*bytes.Buffer) error) (string, error) {
	d, err := b.marshalBinary(f)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(d), nil
}

// unmarshalBinary handles converting a byte array into all the fields of a
// structure that inherits from Base.
// m the marshaler for the OWID
// d the byte array with the data
// f function to add the content from the caller
func (b *Base) unmarshalBinary(
	m owid.Marshaler,
	d []byte,
	f func(*bytes.Buffer) error) error {
	var err error
	u := bytes.NewBuffer(d)

	// Read the version first.
	b.Version, err = common.ReadByte(u)
	if err != nil {
		return err
	}

	// Call the provided function to read the fields for the calling type.
	err = f(u)
	if err != nil {
		return err
	}

	// Finally read the OWID data passing
	b.OWID, err = owid.FromBuffer(u, m)
	if err != nil {
		return err
	}
	return nil
}

// unmarshalString uses the unmarshaler to read the byte array contained in the
// base64 encoded string.
func unmarshalString(b encoding.BinaryUnmarshaler, s string) error {
	d, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	return b.UnmarshalBinary(d)
}
