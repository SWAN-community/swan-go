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

// Bid contains the information about the advert to be displayed.
type Bid struct {
	Response
	MediaURL      string `json:"mediaUrl"`      // The URL of the content of the advert provided in response
	AdvertiserURL string `json:"advertiserURL"` // The URL to direct the browser to if the advert is selected
}

func NewBid(s *owid.Signer, mediaUrl string, advertiserUrl string) (*Bid, error) {
	var err error
	a := &Bid{MediaURL: mediaUrl, AdvertiserURL: advertiserUrl}
	a.Version = swanVersion
	a.StructType = responseBid
	a.OWID, err = s.CreateOWIDandSign(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func BidFromJson(j []byte) (*Bid, error) {
	var a Bid
	err := json.Unmarshal(j, &a)
	if err != nil {
		return nil, err
	}
	a.OWID.Target = &a
	return &a, nil
}

func (a *Bid) ToBase64() (string, error) {
	b, err := a.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func BidFromBase64(value string) (*Bid, error) {
	var a Bid
	err := unmarshalString(&a, value)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *Bid) MarshalOwid() ([]byte, error) {
	return a.marshalOwid(func(b *bytes.Buffer) error { return a.marshal(b) })
}

func (a *Bid) MarshalBinary() ([]byte, error) {
	return a.marshalBinary(func(b *bytes.Buffer) error { return a.marshal(b) })
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

func (a *Bid) UnmarshalBinary(data []byte) error {
	return a.unmarshalBinary(a, data, func(b *bytes.Buffer) error {
		var err error
		if a.StructType != responseBid {
			return fmt.Errorf("struct type not bid '%d'", responseBid)
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
