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
	"strings"

	"github.com/SWAN-community/common-go"
	"github.com/SWAN-community/owid-go"
)

// Seed contains the information about the opportunity to advertise with a
// publisher. It is created and signed by the SWAN Root Party, typically the
// publisher or an agent acting on their behalf.
type Seed struct {
	Base
	PubDomain      string       `json:"pubDomain"`      // The domain that the advertisements will appear on
	TransactionIds [][]byte     `json:"transactionIds"` // An array of transaction ids available in the containing request
	RID            *Identifier  `json:"rid"`            // The Random [browser] Id
	SID            *Identifier  `json:"sid"`            // The Signed in Id
	Preferences    *Preferences `json:"preferences"`    // The privacy preferences
	Stopped        []string     `json:"stopped"`        // List of domains or advert IDs that should not be shown
}

// Returns an OWID with the target populated, or nil of the Seed has not been
// signed.
func (s *Seed) GetOWID() *owid.OWID {
	if s.OWID == nil {
		return nil
	}
	if s.OWID.Target == nil {
		s.OWID.Target = s
	}
	return s.OWID
}

// Returns a new swan.Seed with the correct version and a random uuid ready to
// have the other values added and then signed.
// transactionIds associated with the new seed
func NewSeed(transactionIds [][]byte) (*Seed, error) {
	return &Seed{
		Base:           Base{Version: swanVersion},
		TransactionIds: transactionIds,
	}, nil
}

func SeedUnmarshalBase64(value []byte) (*Seed, error) {
	var s Seed
	err := unmarshalBase64(&s, value)
	if err != nil {
		return nil, err
	}
	s.OWID.Target = &s
	return &s, nil
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

func (s *Seed) MarshalBase64() ([]byte, error) {
	return s.marshalBase64(s.marshal)
}

func (s *Seed) MarshalOwid() ([]byte, error) {
	return s.marshalOwid(s.marshal)
}

func (s *Seed) MarshalBinary() ([]byte, error) {
	return s.marshalBinary(s.marshal)
}

func (s *Seed) marshal(b *bytes.Buffer) error {
	err := common.WriteString(b, s.PubDomain)
	if err != nil {
		return err
	}
	err = common.WriteByteArrayArray(b, s.TransactionIds)
	if err != nil {
		return err
	}
	err = common.WriteMarshaller(b, s.RID)
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
		s.TransactionIds, err = common.ReadByteArrayArray(b)
		if err != nil {
			return err
		}
		if s.RID == nil {
			s.RID = &Identifier{}
		}
		err = common.ReadMarshaller(b, s.RID)
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
			s.SID = &Identifier{}
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
