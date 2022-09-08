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
	"github.com/google/uuid"
)

// Seed contains the information about the opportunity to advertise with a
// publisher. It is created and signed by the SWAN Root Party, typically the
// publisher or an agent acting on their behalf.
type Seed struct {
	Base
	PubDomain   string       `json:"pubDomain"`   // The domain that the advertisements will appear on
	UUID        []byte       `json:"uuid"`        // A unique identifier for this ID
	SWID        *Identifier  `json:"swid"`        // The Secure Web ID
	SID         *Identifier  `json:"sid"`         // The Signed In ID
	Preferences *Preferences `json:"preferences"` // The privacy preferences
	Stopped     []string     `json:"stopped"`     // List of domains or advert IDs that should not be shown
}

// Returns a new swan.Seed with the correct version and a random uuid ready to
// have the other values added and then signed.
func NewSeed() (*Seed, error) {
	uuid, err := uuid.New().MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &Seed{
		Base: Base{Version: 1},
		UUID: uuid,
	}, nil
}

func (s *Seed) Sign(signer *owid.Signer) error {
	var err error
	s.Base.OWID, err = signer.CreateOWIDandSign(s)
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
	var b bytes.Buffer
	err := common.WriteMarshaller(&b, s.SWID.Value)
	if err != nil {
		return nil, err
	}
	err = common.WriteMarshaller(&b, *&s.Preferences)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
