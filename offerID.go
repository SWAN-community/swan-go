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

// OfferID (aka Transaction ID or Bid ID) contains the information about the
// opportunity to advertise with a publisher. It is created by the SWAN host
// as an OWID and as such is signed by the SWAN host and not the publisher.
type OfferID struct {
	Placement   string // A value assigned by the publisher for the advertisement slot on the web page
	PubDomain   string // The domain that the advertisement slot will appear on
	CBID        string // The Commmon Browser ID (not the OWID version)
	SID         string // The Signed In ID (not the OWID version)
	Preferences string // The privacy preferences string (not the OWID version)
}

// NewOfferID creates a new OfferID instance from the string provided. The string
func NewOfferID(s string) (*OfferID, error) {

}

// AsByteArray returns the OfferID as a byte array.
func (o *OfferID) AsByteArray() string {

}

// AsString returns the OfferID as a base 64 encoded string.
func (o *OfferID) AsString() string {

}
