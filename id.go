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
	"owid"
	"strings"

	"github.com/google/uuid"
)

// Used to separate stopped advert IDs in a single string.
const idStoppedSeparator = "\r"

// ID contains the information about the opportunity to advertise with a
// publisher. It is created and signed by the SWAN Root Party, typically the
// publisher or an agent acting on their behalf.
type ID struct {
	base
	PubDomain   string     // The domain that the advertisements will appear on
	UUID        []byte     // A unique identifier for this ID
	SWID        *owid.OWID // The Secure Web ID as an OWID
	SID         *owid.OWID // The Signed In ID as an OWID
	Preferences *owid.OWID // The privacy preferences as an OWID
	Stopped     []string   // List of domains or advert IDs that should not be shown
}

// Returns a new swan.ID with the correct version and type set as well as random
// data to ensure unique for all time.
func NewID() (*ID, error) {
	uuid, err := uuid.New().MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &ID{
		base: base{typeVersion, typeID},
		UUID: uuid,
	}, nil
}

// SWIDAsString as a base 64 string.
func (o *ID) SWIDAsString() string {
	u, err := uuid.FromBytes(o.SWID.Payload)
	if err != nil {
		return o.SWID.PayloadAsPrintable()
	}
	return u.String()
}

// SIDAsString as a base 64 string.
func (o *ID) SIDAsString() string {
	return o.SID.PayloadAsPrintable()
}

// PreferencesAsString as a base 64 string.
func (o *ID) PreferencesAsString() string {
	return o.Preferences.PayloadAsString()
}

// StoppedAsArray returns an array of domains that should not be included in
// bids.
func (o *ID) StoppedAsArray() []string {
	return o.Stopped
}

// IsStopped returns true if the URL provided is stopped.
func (o *ID) IsStopped(u string) bool {
	for _, i := range o.StoppedAsArray() {
		if strings.EqualFold(u, i) {
			return true
		}
	}
	return false
}

// IDFromOWID returns an ID created from the OWID payload.
func IDFromOWID(i *owid.OWID) (*ID, error) {
	var o ID
	buf := bytes.NewBuffer(i.Payload)
	err := o.setFromBuffer(buf)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// IDFromNode returns an ID created from the Node payload.
func IDFromNode(n *owid.Node) (*ID, error) {
	var o ID
	w, err := n.GetOWID()
	if err != nil {
		return nil, err
	}
	f := bytes.NewBuffer(w.Payload)
	err = o.setFromBuffer(f)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// AsByteArray returns the ID as a byte array.
func (o *ID) AsByteArray() ([]byte, error) {
	var buf bytes.Buffer
	o.writeToBuffer(&buf)
	return buf.Bytes(), nil
}

// AsString returns the ID as a string.
func (o *ID) AsString() (string, error) {
	b, err := o.AsByteArray()
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(b), nil
}

func (o *ID) writeToBuffer(f *bytes.Buffer) error {
	o.base.version = typeVersion
	o.base.structType = typeID
	err := o.base.writeToBuffer(f)
	if err != nil {
		return err
	}
	err = writeString(f, o.PubDomain)
	if err != nil {
		return err
	}
	err = writeByteArray(f, o.UUID)
	if err != nil {
		return err
	}
	err = o.SWID.ToBuffer(f)
	if err != nil {
		return err
	}
	err = o.SID.ToBuffer(f)
	if err != nil {
		return err
	}
	err = o.Preferences.ToBuffer(f)
	if err != nil {
		return err
	}
	err = writeString(f, strings.Join(o.Stopped, idStoppedSeparator))
	if err != nil {
		return err
	}
	return nil
}

func (o *ID) setFromBuffer(f *bytes.Buffer) error {
	var err error
	err = o.base.setFromBuffer(f)
	if err != nil {
		return err
	}
	if o.structType != typeID {
		return fmt.Errorf(
			"type %s not valid for %s",
			typeAsString(o.structType),
			typeAsString(typeID))
	}
	switch o.base.version {
	case byte(1):
		err = o.setFromBufferVersion1(f)
		if err != nil {
			return err
		}
	default:
		err = fmt.Errorf("version '%d' not supported", o.base.version)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ID) setFromBufferVersion1(f *bytes.Buffer) error {
	var err error
	o.PubDomain, err = readString(f)
	if err != nil {
		return err
	}
	o.UUID, err = readByteArray(f)
	if err != nil {
		return err
	}
	o.SWID, err = owid.FromBuffer(f)
	if err != nil {
		return err
	}
	o.SID, err = owid.FromBuffer(f)
	if err != nil {
		return err
	}
	o.Preferences, err = owid.FromBuffer(f)
	if err != nil {
		return err
	}
	s, err := readString(f)
	if err != nil {
		return err
	}
	o.Stopped = strings.Split(s, idStoppedSeparator)
	return nil
}
