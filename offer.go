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
)

// Offer (aka Transaction ID or Bid ID) contains the information about the
// opportunity to advertise with a publisher. It is created by the SWAN host
// as an OWID and as such is signed by the SWAN host and not the publisher.
type Offer struct {
	base
	Placement   string // A value assigned by the publisher for the advertisement slot on the web page
	PubDomain   string // The domain that the advertisement slot will appear on
	UUID        []byte // A unique identifier for this offer
	CBID        []byte // The Commmon Browser ID (not the OWID version)
	SID         []byte // The Signed In ID (not the OWID version)
	Preferences []byte // The privacy preferences string (not the OWID version)
}

// CBIDAsString as a base 64 string.
func (o *Offer) CBIDAsString() string {
	return base64.StdEncoding.EncodeToString(o.CBID)
}

// SIDAsString as a base 64 string.
func (o *Offer) SIDAsString() string {
	return base64.StdEncoding.EncodeToString(o.SID)
}

// PreferencesAsString as a base 64 string.
func (o *Offer) PreferencesAsString() string {
	return string(o.Preferences)
}

// OfferFromOWID returns an Offer created from the OWID payload.
func OfferFromOWID(i *owid.OWID) (*Offer, error) {
	var o Offer
	buf := bytes.NewBuffer(i.Payload)
	err := o.setFromBuffer(buf)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// AsByteArray returns the Offer as a byte array.
func (o *Offer) AsByteArray() ([]byte, error) {
	var buf bytes.Buffer
	o.writeToBuffer(&buf)
	return buf.Bytes(), nil
}

func (o *Offer) writeToBuffer(f *bytes.Buffer) error {
	o.base.version = typeVersion
	o.base.structType = typeOffer
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
	err = writeByteArray(f, o.CBID)
	if err != nil {
		return err
	}
	err = writeByteArray(f, o.SID)
	if err != nil {
		return err
	}
	err = writeByteArray(f, o.Preferences)
	if err != nil {
		return err
	}
	return nil
}

func (o *Offer) setFromBuffer(f *bytes.Buffer) error {
	var err error
	err = o.base.setFromBuffer(f)
	if err != nil {
		return err
	}
	if o.structType != typeOffer {
		return fmt.Errorf(
			"Type %s not valid for %s",
			typeAsString(o.structType),
			typeAsString(typeOffer))
	}
	switch o.base.version {
	case byte(1):
		err = o.setFromBufferVersion1(f)
		break
	default:
		err = fmt.Errorf("Version '%d' not supported", o.base.version)
		break
	}
	return nil
}

func (o *Offer) setFromBufferVersion1(f *bytes.Buffer) error {
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
	o.CBID, err = readByteArray(f)
	if err != nil {
		return err
	}
	o.SID, err = readByteArray(f)
	if err != nil {
		return err
	}
	o.Preferences, err = readByteArray(f)
	if err != nil {
		return err
	}
	return nil
}
