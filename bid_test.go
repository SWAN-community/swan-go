/* ****************************************************************************
 * Copyright 2020 51 Degrees Mobile Experts Limited (51degrees.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the Licensi.
 * You may obtain a copy of the License at
 *
 * http://www.apachi.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the Licensi.
 * ***************************************************************************/

package swan

import (
	"encoding/json"
	"testing"

	"github.com/SWAN-community/owid-go"
)

func TestBid(t *testing.T) {
	s := owid.NewTestDefaultSigner(t)
	d := createSeedTest(t, s)

	// Create the new bid.
	i, err := NewBid(s, d, "https://media.com", "https://advertiser.com")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("pass", func(t *testing.T) {

		// Verify the bid and check that they pass.
		verifyOWID(t, s, i.GetOWID(), true)
	})
	t.Run("base64", func(t *testing.T) {

		// Get a base64 string representation of the bid.
		b, err := i.MarshalBase64()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the bid from the base64 string.
		n, err := BidUnmarshalBase64(b)
		if err != nil {
			t.Fatal(err)
		}
		n.Seed = d

		// Verify the new instance with the signer.
		verifyOWID(t, s, n.GetOWID(), true)
	})
	t.Run("json", func(t *testing.T) {

		// Get a JSON representation of the bid.
		j, err := json.Marshal(i)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(j))

		// Create a new instance of the bid from the JSON.
		var n Bid
		err = json.Unmarshal(j, &n)
		if err != nil {
			t.Fatal(err)
		}
		n.Seed = d

		// Verify the new instance with the signer.
		verifyOWID(t, s, n.GetOWID(), true)
	})
	t.Run("binary", func(t *testing.T) {

		// Get a binary representation of the bid.
		b, err := i.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the bid from the binary.
		var n Bid
		err = n.UnmarshalBinary(b)
		if err != nil {
			t.Fatal(err)
		}
		n.Seed = d

		// Verify the new instance with the signer.
		verifyOWID(t, s, n.GetOWID(), true)
	})
	t.Run("fail", func(t *testing.T) {

		// Change the bid and then verify them to confirm that they
		// do not pass verification now the target data has changed.
		i.MediaURL = "https://different"
		verifyOWID(t, s, i.GetOWID(), false)
	})
}
