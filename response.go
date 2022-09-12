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

	"github.com/SWAN-community/common-go"
	"github.com/SWAN-community/owid-go"
)

// First byte of the data structure will be the type of response.
const (
	responseBid byte = iota + 1
	responseID
	responseFailed
	responseEmpty
)

// Response from an OpenRTB transation.
type Response struct {
	Base
	StructType byte `json:"type"` // The type of structure the response relates to
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

// marshalOwid returns a byte array of all the data needed by an OWID.
func (r *Response) marshalOwid(f func(*bytes.Buffer) error) ([]byte, error) {
	var u bytes.Buffer
	err := r.writeData(&u, f)
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

// ResponseFromJSON returns an instance of Bid, Failed, or Empty for the JSON
// provided, or an error if the JSON can not be unmarshalled to a response.
func ResponseFromJSON(j []byte) (interface{}, error) {
	var r Response
	err := json.Unmarshal(j, &r)
	if err != nil {
		return nil, err
	}
	var i interface{ owid.Marshaler }
	var b *Base
	switch r.StructType {
	case responseBid:
		var n Bid
		b = &n.Base
		i = &n
	case responseEmpty:
		var n Empty
		b = &n.Base
		i = &n
	case responseFailed:
		var n Failed
		b = &n.Base
		i = &n
	default:
		return nil, fmt.Errorf("type '%d' unknown", r.StructType)
	}
	json.Unmarshal(j, i)
	if err != nil {
		return nil, err
	}
	b.OWID.Target = i
	return i, nil
}
