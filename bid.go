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

	"github.com/SWAN-community/owid-go"
)

// Bid contains the information about the advert to be displayed.
type Bid struct {
	base
	MediaURL      string // The URL of the content of the advert provided in response
	AdvertiserURL string // The URL to direct the browser to if the advert is selected
}

// BidFromOWID returns a Bid created from the OWID payload.
func BidFromOWID(i *owid.OWID) (*Bid, error) {
	var b Bid
	f := bytes.NewBuffer(i.Payload)
	err := b.setFromBuffer(f)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// BidFromNode returns a Bid created from the Node payload.
func BidFromNode(n *owid.Node) (*Bid, error) {
	var b Bid
	o, err := n.GetOWID()
	if err != nil {
		return nil, err
	}
	f := bytes.NewBuffer(o.Payload)
	err = b.setFromBuffer(f)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// AsByteArray returns the Bid as a byte array.
func (b *Bid) AsByteArray() ([]byte, error) {
	var f bytes.Buffer
	b.writeToBuffer(&f)
	return f.Bytes(), nil
}

func (b *Bid) writeToBuffer(f *bytes.Buffer) error {
	b.version = typeVersion
	b.structType = typeBid
	err := b.base.writeToBuffer(f)
	if err != nil {
		return err
	}
	err = writeString(f, b.MediaURL)
	if err != nil {
		return err
	}
	err = writeString(f, b.AdvertiserURL)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bid) setFromBuffer(f *bytes.Buffer) error {
	err := b.base.setFromBuffer(f)
	if err != nil {
		return err
	}
	if b.structType != typeBid {
		return fmt.Errorf(
			"Type %s not valid for %s",
			typeAsString(b.structType),
			typeAsString(typeBid))
	}
	switch b.base.version {
	case byte(1):
		err = b.setFromBufferVersion1(f)
		break
	default:
		err = fmt.Errorf("Version '%d' not supported", b.base.version)
		break
	}
	return err
}

func (b *Bid) setFromBufferVersion1(f *bytes.Buffer) error {
	var err error
	b.MediaURL, err = readString(f)
	if err != nil {
		return err
	}
	b.AdvertiserURL, err = readString(f)
	if err != nil {
		return err
	}
	return nil
}
