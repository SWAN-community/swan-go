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

	"github.com/SWAN-community/common-go"
	"github.com/SWAN-community/owid-go"
)

// Failed contains details about the request that was not signed by the
// recipient.
type Failed struct {
	Response
	Host  string `json:"host"`  // The domain that did not respond.
	Error string `json:"error"` // The error message to add to the tree.
}

func NewFailed(s *owid.Signer, host string, message string) (*Failed, error) {
	var err error
	f := &Failed{Host: host, Error: message}
	f.Version = swanVersion
	f.StructType = responseFailed
	f.OWID, err = s.CreateOWIDandSign(f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func FailedFromJson(j []byte) (*Failed, error) {
	var a Failed
	err := json.Unmarshal(j, &a)
	if err != nil {
		return nil, err
	}
	a.OWID.Target = &a
	return &a, nil
}

func (f *Failed) ToBase64() (string, error) {
	b, err := f.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func FailedFromBase64(value string) (*Failed, error) {
	var a Failed
	err := unmarshalString(&a, value)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (f *Failed) MarshalOwid() ([]byte, error) {
	return f.marshalOwid(func(b *bytes.Buffer) error { return f.marshal(b) })
}

func (f *Failed) MarshalBinary() ([]byte, error) {
	return f.marshalBinary(func(b *bytes.Buffer) error { return f.marshal(b) })
}

func (f *Failed) marshal(b *bytes.Buffer) error {
	err := common.WriteString(b, f.Host)
	if err != nil {
		return err
	}
	err = common.WriteString(b, f.Error)
	if err != nil {
		return err
	}
	return nil
}

func (f *Failed) UnmarshalBinary(data []byte) error {
	return f.unmarshalBinary(f, data, func(b *bytes.Buffer) error {
		var err error
		if f.StructType != responseFailed {
			return fmt.Errorf("struct type not failed '%d'", responseFailed)
		}
		f.Host, err = common.ReadString(b)
		if err != nil {
			return err
		}
		f.Error, err = common.ReadString(b)
		if err != nil {
			return err
		}
		return nil
	})
}
