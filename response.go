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
	"fmt"

	"github.com/SWAN-community/common-go"
	"github.com/SWAN-community/owid-go"
)

// First byte of the data structure will be the type of response.
const (
	responseBid byte = iota + 1
	responseFailed
	responseEmpty
)

// Response from an OpenRTB transation.
type Response struct {
	Base
	StructType byte  `json:"type"` // The type of structure the response relates to
	Seed       *Seed `json:"-"`    // The seed for the transmission
}

// writeData writes the base and type before calling the function.
func (r *Response) writeData(u *bytes.Buffer, f func(*bytes.Buffer) error) error {
	return r.Base.writeData(u, func(b *bytes.Buffer) error {
		err := common.WriteByte(b, r.StructType)
		if err != nil {
			return err
		}
		err = f(u)
		if err != nil {
			return err
		}
		return nil
	})
}

// marshalBase64 encodes the marshalBinary result as a base64 string.
func (r *Response) marshalBase64(f func(*bytes.Buffer) error) ([]byte, error) {
	s, err := r.marshalBinary(f)
	if err != nil {
		return nil, err
	}
	return []byte(base64.StdEncoding.EncodeToString(s)), nil
}

// marshalOwid returns a byte array of all the data that forms the response
// instance AND the the seed.
func (r *Response) marshalOwid(f func(*bytes.Buffer) error) ([]byte, error) {
	var u bytes.Buffer
	err := r.writeData(&u, f)
	if err != nil {
		return nil, err
	}
	if r.Seed == nil {
		return nil, fmt.Errorf("missing seed")
	}
	err = common.WriteMarshaller(&u, r.Seed)
	if err != nil {
		return nil, err
	}
	return u.Bytes(), nil
}

// marshalBinary marshals the version, calls the function to add more data, and
// finishes by adding the OWID before returning the byte array.
func (r *Response) marshalBinary(f func(*bytes.Buffer) error) ([]byte, error) {
	var u bytes.Buffer
	err := r.writeData(&u, f)
	if err != nil {
		return nil, err
	}
	err = r.OWID.ToBuffer(&u)
	if err != nil {
		return nil, err
	}
	return u.Bytes(), nil
}

// unmarshalBinary handles converting a byte array into all the fields of a
// structure that inherits from Response.
// m the marshaler for the OWID
// d the byte array with the data
// f function to add the content from the caller
func (r *Response) unmarshalBinary(
	m owid.Marshaler,
	d []byte,
	f func(*bytes.Buffer) error) error {
	return r.Base.unmarshalBinary(m, d, func(b *bytes.Buffer) error {
		var err error
		r.StructType, err = common.ReadByte(b)
		if err != nil {
			return err
		}
		err = f(b)
		if err != nil {
			return err
		}
		return nil
	})
}

// ResponseFromByteArray turns the byte array into an instance of a structure
// that includes swan.Response. Either Bid, Failed, or Empty.
// Intended to be used to pass individual responses as string parameters.
func ResponseFromByteArray(data []byte) (Signed, error) {
	var r Response
	err := r.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}
	return r.UnmarshalBinaryToType(data)
}

// ResponseFromBase64 turns the base64 string into an instance of a structure
// that includes swan.Response. Either Bid, Failed, or Empty.
// Intended to be used to pass individual responses as string parameters.
func ResponseFromBase64(data []byte) (Signed, error) {
	b, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	return ResponseFromByteArray(b)
}

// UnmarshalBinary reads the response version and structure type and ignores
// the rest of the data. Used to determine the type of response.
func (r *Response) UnmarshalBinary(data []byte) error {
	var err error
	b := bytes.NewBuffer(data)
	r.Version, err = common.ReadByte(b)
	if err != nil {
		return err
	}
	r.StructType, err = common.ReadByte(b)
	if err != nil {
		return err
	}
	return nil
}

// UnmarshalBinaryToType
func (r *Response) UnmarshalBinaryToType(data []byte) (Signed, error) {
	var err error
	var i Signed
	switch r.StructType {
	case responseBid:
		var b Bid
		err = b.UnmarshalBinary(data)
		b.OWID.Target = &b
		i = &b
	case responseEmpty:
		var e Empty
		err = e.UnmarshalBinary(data)
		e.OWID.Target = &e
		i = &e
	case responseFailed:
		var f Failed
		err = f.UnmarshalBinary(data)
		f.OWID.Target = &f
		i = &f
	default:
		return nil, fmt.Errorf("struct '%d' unknown", r.StructType)
	}
	if err != nil {
		return nil, err
	}
	return i, nil
}
