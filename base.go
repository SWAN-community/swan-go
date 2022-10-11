/* ****************************************************************************
 * Copyright 2022 51 Degrees Mobile Experts Limited (51degrees.com)
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
	"fmt"

	"github.com/SWAN-community/common-go"
	"github.com/SWAN-community/owid-go"
	"github.com/SWAN-community/swift-go"
)

// Base used with any SWAN field.
type Base struct {
	Version byte       `json:"version"` // Used to indicate the version encoding of the type.
	OWID    *owid.OWID `json:"source"`  // OWID related to the structure
}

// GetVersion used to indicate the version encoding of the type.
func (b *Base) getVersion() byte {
	return b.Version
}

// Returns true if the structure is signed, otherwise false.
func (b *Base) IsSigned() bool {
	return b.OWID != nil && b.OWID.Signature != nil && len(b.OWID.Signature) > 0
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
	if b.Version < swanMinVersion || b.Version > swanMaxVersion {
		return fmt.Errorf("version '%d' not supported", b.Version)
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

// marshalBase64 encodes the marshalBinary result as a base64 string.
func (b *Base) marshalBase64(f func(*bytes.Buffer) error) ([]byte, error) {
	s, err := b.marshalBinary(f)
	if err != nil {
		return nil, err
	}
	return []byte(base64.StdEncoding.EncodeToString(s)), nil
}

// unmarshalBase64 uses the unmarshaler to read the byte array contained in the
// base64 encoded string.
func unmarshalBase64(b encoding.BinaryUnmarshaler, s []byte) error {
	d, err := base64.StdEncoding.DecodeString(string(s))
	if err != nil {
		return err
	}
	return b.UnmarshalBinary(d)
}

// validateSwiftPair checks that there is only one entry in the value array.
func validateSwiftPair(p *swift.Pair) error {
	if len(p.Values()) != 1 {
		return fmt.Errorf("%s pair value must have length 1", p.Key())
	}
	return nil
}
