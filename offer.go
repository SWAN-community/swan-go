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
	"fmt"
	"owid"
	"strings"

	"github.com/google/uuid"
)

// Used to separate stopped advert IDs in a single string.
const impressionStoppedSeparator = "\r"

// Impression contains the information about the
// opportunity to advertise with a publisher. It is created by the SWAN host
// as an OWID and as such is signed by the SWAN host and not the publisher.
type Impression struct {
	base
	Placement   string     // A value assigned by the publisher for the advertisement slot on the web page
	PubDomain   string     // The domain that the advertisement slot will appear on
	UUID        []byte     // A unique identifier for this impression
	SWID        *owid.OWID // The Commmon Browser ID
	SID         *owid.OWID // The Signed In ID
	Preferences *owid.OWID // The privacy preferences string
	Stopped     []string   // List of domains of advert IDs that should not be shown
}

// Returns a new Impression with the correct version and impression information set.
func NewImpression() Impression {
	return Impression{base: base{typeVersion, typeImpression}}
}

// SWIDAsString as a base 64 string.
func (o *Impression) SWIDAsString() string {
	u, err := uuid.FromBytes(o.SWID.Payload)
	if err != nil {
		return o.SWID.PayloadAsPrintable()
	}
	return u.String()
}

// SIDAsString as a base 64 string.
func (o *Impression) SIDAsString() string {
	return o.SID.PayloadAsPrintable()
}

// PreferencesAsString as a base 64 string.
func (o *Impression) PreferencesAsString() string {
	return o.Preferences.PayloadAsString()
}

// StoppedAsArray returns an array of domains that should not be included in
// bids.
func (o *Impression) StoppedAsArray() []string {
	return o.Stopped
}

// IsStopped returns true if the URL provided is stopped.
func (o *Impression) IsStopped(u string) bool {
	for _, i := range o.StoppedAsArray() {
		if strings.EqualFold(u, i) {
			return true
		}
	}
	return false
}

// ImpressionFromOWID returns an Impression created from the OWID payload.
func ImpressionFromOWID(i *owid.OWID) (*Impression, error) {
	var o Impression
	buf := bytes.NewBuffer(i.Payload)
	err := o.setFromBuffer(buf)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// ImpressionFromNode returns an Impression created from the Node payload.
func ImpressionFromNode(n *owid.Node) (*Impression, error) {
	var o Impression
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

// AsByteArray returns the Impression as a byte array.
func (o *Impression) AsByteArray() ([]byte, error) {
	var buf bytes.Buffer
	o.writeToBuffer(&buf)
	return buf.Bytes(), nil
}

func (o *Impression) writeToBuffer(f *bytes.Buffer) error {
	o.base.version = typeVersion
	o.base.structType = typeImpression
	err := o.base.writeToBuffer(f)
	if err != nil {
		return err
	}
	err = writeString(f, o.Placement)
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
	err = writeString(f, strings.Join(o.Stopped, impressionStoppedSeparator))
	if err != nil {
		return err
	}
	return nil
}

func (o *Impression) setFromBuffer(f *bytes.Buffer) error {
	var err error
	err = o.base.setFromBuffer(f)
	if err != nil {
		return err
	}
	if o.structType != typeImpression {
		return fmt.Errorf(
			"type %s not valid for %s",
			typeAsString(o.structType),
			typeAsString(typeImpression))
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

func (o *Impression) setFromBufferVersion1(f *bytes.Buffer) error {
	var err error
	o.Placement, err = readString(f)
	if err != nil {
		return err
	}
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
	o.Stopped = strings.Split(s, impressionStoppedSeparator)
	return nil
}
