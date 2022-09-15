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

	"github.com/SWAN-community/common-go"
	"github.com/SWAN-community/owid-go"
)

// Bid contains the information about the advert to be displayed.
type Bid struct {
	Response
	MediaURL      string `json:"mediaUrl"`      // The URL of the content of the advert provided in response
	AdvertiserURL string `json:"advertiserURL"` // The URL to direct the browser to if the advert is selected
}

func (b *Bid) GetOWID() *owid.OWID {
	if b.OWID.Target == nil {
		b.OWID.Target = b
	}
	return b.OWID
}

func NewBid(
	signer *owid.Signer,
	seed *Seed,
	mediaUrl string,
	advertiserUrl string) (*Bid, error) {
	var err error
	a := &Bid{MediaURL: mediaUrl, AdvertiserURL: advertiserUrl}
	a.Version = swanVersion
	a.StructType = responseBid
	a.Seed = seed
	a.OWID, err = signer.CreateOWIDandSign(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func BidUnmarshalBase64(value []byte) (*Bid, error) {
	var a Bid
	err := unmarshalBase64(&a, value)
	if err != nil {
		return nil, err
	}
	a.OWID.Target = &a
	return &a, nil
}

func (a *Bid) MarshalBase64() ([]byte, error) {
	return a.Response.marshalBase64(a.marshal)
}

func (a *Bid) MarshalOwid() ([]byte, error) {
	return a.Response.marshalOwid(a.marshal)
}

func (a *Bid) MarshalBinary() ([]byte, error) {
	return a.Response.marshalBinary(a.marshal)
}

func (a *Bid) UnmarshalBinary(data []byte) error {
	return a.Response.unmarshalBinary(a, data, func(b *bytes.Buffer) error {
		var err error
		if a.StructType != responseBid {
			return fmt.Errorf(
				"struct type '%d' not bid '%d'",
				a.StructType,
				responseBid)
		}
		a.MediaURL, err = common.ReadString(b)
		if err != nil {
			return err
		}
		a.AdvertiserURL, err = common.ReadString(b)
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *Bid) marshal(b *bytes.Buffer) error {
	err := common.WriteString(b, a.MediaURL)
	if err != nil {
		return err
	}
	err = common.WriteString(b, a.AdvertiserURL)
	if err != nil {
		return err
	}
	return nil
}
