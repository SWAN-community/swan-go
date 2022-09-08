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
	"github.com/google/uuid"
)

// Identifier represents a OneKey compatible random identifier.
// https://github.com/OneKey-Network/addressability-framework/blob/main/mvp-spec/model/identifier.md
type Identifier struct {
	Base
	IdType    string    `json:"type"`  // Type of identifier
	Value     uuid.UUID `json:"value"` // In practice the value is a UUID so store it as one
	Persisted bool      // True if the value has been stored.
}

func NewIdentifier(
	s *owid.Signer,
	idType string,
	value uuid.UUID) (*Identifier, error) {
	var err error
	i := &Identifier{IdType: idType, Value: value}
	i.Base.Version = swanVersion
	i.Base.OWID, err = s.CreateOWIDandSign(i)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func IdentifierFromJson(j []byte) (*Identifier, error) {
	var i Identifier
	err := json.Unmarshal(j, &i)
	if err != nil {
		return nil, err
	}
	i.OWID.Target = &i
	return &i, nil
}

func IdentifierFromBase64(value string) (*Identifier, error) {
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	var i Identifier
	err = i.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (i *Identifier) ToBase64() (string, error) {
	b, err := i.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (i *Identifier) marshal(b *bytes.Buffer) error {
	err := common.WriteByte(b, i.Base.Version)
	if err != nil {
		return err
	}
	err = common.WriteString(b, i.IdType)
	if err != nil {
		return err
	}
	err = common.WriteMarshaller(b, i.Value)
	if err != nil {
		return err
	}
	return nil
}

func (i *Identifier) MarshalOwid() ([]byte, error) {
	var b bytes.Buffer
	err := i.marshal(&b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (i *Identifier) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	err := i.marshal(&b)
	if err != nil {
		return nil, err
	}
	err = i.Base.OWID.ToBuffer(&b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (i *Identifier) UnmarshalBinary(data []byte) error {
	var err error
	b := bytes.NewBuffer(data)
	i.Base.Version, err = common.ReadByte(b)
	if err != nil {
		return err
	}
	i.IdType, err = common.ReadString(b)
	if err != nil {
		return err
	}
	u, err := common.ReadByteArray(b)
	if err != nil {
		return err
	}
	err = i.Value.UnmarshalBinary(u)
	if err != nil {
		return err
	}
	i.Base.OWID, err = owid.FromBuffer(b, i)
	if err != nil {
		return err
	}
	return nil
}
