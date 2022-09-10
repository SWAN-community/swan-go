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
	"strings"

	"github.com/SWAN-community/common-go"
	"github.com/SWAN-community/owid-go"
	"github.com/google/uuid"
)

// Seed contains the information about the opportunity to advertise with a
// publisher. It is created and signed by the SWAN Root Party, typically the
// publisher or an agent acting on their behalf.
type Seed struct {
	Base
	PubDomain   string       `json:"pubDomain"`   // The domain that the advertisements will appear on
	UUID        uuid.UUID    `json:"uuid"`        // A unique identifier for this ID
	SWID        *Identifier  `json:"swid"`        // The Secure Web ID
	SID         *ByteArray   `json:"sid"`         // The Signed In ID
	Preferences *Preferences `json:"preferences"` // The privacy preferences
	Stopped     []string     `json:"stopped"`     // List of domains or advert IDs that should not be shown
}

// Returns a new swan.Seed with the correct version and a random uuid ready to
// have the other values added and then signed.
func NewSeed() (*Seed, error) {
	return &Seed{
		Base: Base{Version: swanVersion},
		UUID: uuid.New(),
	}, nil
}

func SeedFromJson(j []byte) (*Seed, error) {
	var s Seed
	err := json.Unmarshal(j, &s)
	if err != nil {
		return nil, err
	}
	s.OWID.Target = &s
	return &s, nil
}

func SeedFromBase64(value string) (*Seed, error) {
	var s Seed
	err := unmarshalString(&s, value)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Seed) ToBase64() (string, error) {
	b, err := s.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// Sign the seed including all the fields included.
func (s *Seed) Sign(signer *owid.Signer) error {
	var err error
	s.OWID, err = signer.CreateOWIDandSign(s)
	if err != nil {
		return err
	}
	return nil
}

// IsStopped returns true if the URL provided is stopped.
func (s *Seed) IsStopped(u string) bool {
	for _, i := range s.Stopped {
		if strings.EqualFold(u, i) {
			return true
		}
	}
	return false
}

func (s *Seed) MarshalOwid() ([]byte, error) {
	return s.marshalOwid(func(b *bytes.Buffer) error { return s.marshal(b) })
}

func (s *Seed) MarshalBinary() ([]byte, error) {
	return s.marshalBinary(func(b *bytes.Buffer) error { return s.marshal(b) })
}

func (s *Seed) marshal(b *bytes.Buffer) error {
	err := common.WriteString(b, s.PubDomain)
	if err != nil {
		return err
	}
	err = common.WriteMarshaller(b, s.UUID)
	if err != nil {
		return err
	}
	err = common.WriteMarshaller(b, s.SWID)
	if err != nil {
		return err
	}
	err = common.WriteMarshaller(b, s.Preferences)
	if err != nil {
		return err
	}
	err = common.WriteMarshaller(b, s.SID)
	if err != nil {
		return err
	}
	err = common.WriteStrings(b, s.Stopped)
	if err != nil {
		return err
	}
	return nil
}

func (s *Seed) UnmarshalBinary(data []byte) error {
	return s.unmarshalBinary(s, data, func(b *bytes.Buffer) error {
		var err error
		s.PubDomain, err = common.ReadString(b)
		if err != nil {
			return err
		}
		err = common.ReadMarshaller(b, &s.UUID)
		if err != nil {
			return err
		}
		if s.SWID == nil {
			s.SWID = &Identifier{}
		}
		err = common.ReadMarshaller(b, s.SWID)
		if err != nil {
			return err
		}
		if s.Preferences == nil {
			s.Preferences = &Preferences{}
		}
		err = common.ReadMarshaller(b, s.Preferences)
		if err != nil {
			return err
		}
		if s.SID == nil {
			s.SID = &ByteArray{}
		}
		err = common.ReadMarshaller(b, s.SID)
		if err != nil {
			return err
		}
		s.Stopped, err = common.ReadStrings(b)
		if err != nil {
			return err
		}
		return nil
	})
}
