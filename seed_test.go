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
	"encoding/json"
	"testing"

	"github.com/SWAN-community/owid-go"
	"github.com/google/uuid"
)

func TestSeed(t *testing.T) {
	s := owid.NewTestDefaultSigner(t)

	// Create the new seed.
	d, err := NewSeed()
	if err != nil {
		t.Fatal(err)
	}

	// Add the simple fields.
	d.PubDomain = "test.com"
	d.Stopped = []string{"a.com", "b.com"}

	// Create the new preferences.
	d.Preferences, err = NewPreferences(s, true)
	if err != nil {
		t.Fatal(err)
	}

	// Create the new identifier.
	u, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	d.SWID, err = NewIdentifier(s, "type", u)
	if err != nil {
		t.Fatal(err)
	}

	// Create the new byte array.
	d.SID, err = NewByteArray(s, []byte{1, 2, 3, 4})
	if err != nil {
		t.Fatal(err)
	}

	err = d.Sign(s)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("pass", func(t *testing.T) {

		// Verify the seed and check that they pass.
		verifyBase(t, s, &d.Base, true)
	})
	t.Run("json", func(t *testing.T) {

		// Get a JSON representation of the seed.
		j, err := json.Marshal(d)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the seed from the JSON.
		n, err := SeedFromJson(j)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyBase(t, s, &n.Base, true)
	})
	t.Run("binary", func(t *testing.T) {

		// Get a binary representation of the seed.
		b, err := d.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the seed from the binary.
		var n Seed
		err = n.UnmarshalBinary(b)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyBase(t, s, &n.Base, true)
	})
	t.Run("fail", func(t *testing.T) {

		// Change the seed and then verify them to confirm that they
		// do not pass verification now the target data has changed.
		d.PubDomain = "different.com"
		verifyBase(t, s, &d.Base, false)
	})
}
