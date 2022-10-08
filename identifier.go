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
	"encoding/hex"
	"net/http"

	"github.com/SWAN-community/common-go"
	"github.com/SWAN-community/owid-go"
	"github.com/SWAN-community/swift-go"
)

// Identifier represents a OneKey compatible random identifier.
// https://github.com/OneKey-Network/addressability-framework/blob/main/mvp-spec/model/identifier.md
type Identifier struct {
	Writeable
	IdType string `json:"type"` // Type of identifier
	// Byte array for the identifier, usual a 16 byte UUID or the result of a hash operation
	Value []byte `json:"value"`
}

func (i *Identifier) GetOWID() *owid.OWID {
	if i.OWID.Target == nil {
		i.OWID.Target = i
	}
	return i.OWID
}

func (i *Identifier) GetCookie() *Cookie {
	if i.Cookie == nil {
		i.Cookie = &Cookie{Created: i.GetOWID().TimeStamp}
	}
	return i.Cookie
}

func (i *Identifier) AsPrintable() string {
	if i.Value != nil {
		return hex.EncodeToString(i.Value)
	}
	return ""
}

func (i *Identifier) AsHttpCookie(
	host string,
	secure bool) (*http.Cookie, error) {
	return i.GetCookie().asHttpCookie(host, secure, i)
}

func NewIdentifier(
	s *owid.Signer,
	idType string,
	source encoding.BinaryMarshaler) (*Identifier, error) {
	v, err := source.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return NewIdentifierFromByteArray(s, idType, v)
}

func NewIdentifierFromByteArray(
	s *owid.Signer,
	idType string,
	source []byte) (*Identifier, error) {
	var err error
	i := &Identifier{IdType: idType, Value: source}
	i.Version = swanVersion
	i.OWID, err = s.CreateOWIDandSign(i)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func IdentifierUnmarshalBase64(value []byte) (*Identifier, error) {
	var i Identifier
	err := i.UnmarshalBase64(value)
	if err != nil {
		return nil, err
	}
	i.OWID.Target = &i
	return &i, nil
}

func (i *Identifier) UnmarshalSwift(p *swift.Pair) error {
	if len(p.Values()) == 0 {
		return nil
	}
	err := i.UnmarshalBase64(p.Values()[0])
	if err != nil {
		return err
	}
	i.Cookie = &Cookie{}
	return i.Cookie.UnmarshalSwiftValidity(p)
}

func (i *Identifier) UnmarshalBase64(value []byte) error {
	return unmarshalBase64(i, value)
}

func (i *Identifier) MarshalBase64() ([]byte, error) {
	return i.marshalBase64(i.marshal)
}

func (i *Identifier) MarshalOwid() ([]byte, error) {
	return i.marshalOwid(i.marshal)
}

func (i *Identifier) MarshalBinary() ([]byte, error) {
	return i.marshalBinary(i.marshal)
}

func (i *Identifier) marshal(b *bytes.Buffer) error {
	err := common.WriteString(b, i.IdType)
	if err != nil {
		return err
	}
	err = common.WriteByteArray(b, i.Value)
	if err != nil {
		return err
	}
	return nil
}

func (i *Identifier) UnmarshalBinary(data []byte) error {
	return i.unmarshalBinary(i, data, func(b *bytes.Buffer) error {
		var err error
		i.IdType, err = common.ReadString(b)
		if err != nil {
			return err
		}
		i.Value, err = common.ReadByteArray(b)
		if err != nil {
			return err
		}
		return nil
	})
}
