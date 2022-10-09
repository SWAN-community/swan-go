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

func TestFailed(t *testing.T) {
	s := owid.NewTestDefaultSigner(t)
	d := createSeedTest(t, s)

	// Create the new failed.
	f, err := NewFailed(s, d, "test.com", "message")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("pass", func(t *testing.T) {

		// Verify the failed and check that they pass.
		verifyOWID(t, s, f, true)
	})
	t.Run("base64", func(t *testing.T) {

		// Get a base64 string representation of the failed.
		b, err := f.MarshalBase64()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the failed from the base64 string.
		n, err := FailedUnmarshalBase64(b)
		if err != nil {
			t.Fatal(err)
		}
		n.Seed = d

		// Verify the new instance with the signer.
		verifyOWID(t, s, n, true)
	})
	t.Run("json", func(t *testing.T) {

		// Get a JSON representation of the failed.
		j, err := json.Marshal(f)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(j))

		// Create a new instance of the failed from the JSON.
		var n Failed
		err = json.Unmarshal(j, &n)
		if err != nil {
			t.Fatal(err)
		}
		n.Seed = d

		// Verify the new instance with the signer.
		verifyOWID(t, s, &n, true)
	})
	t.Run("binary", func(t *testing.T) {

		// Get a binary representation of the failed.
		b, err := f.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the failed from the binary.
		var n Failed
		err = n.UnmarshalBinary(b)
		if err != nil {
			t.Fatal(err)
		}
		n.Seed = d

		// Verify the new instance with the signer.
		verifyOWID(t, s, &n, true)
	})
	t.Run("response", func(t *testing.T) {

		// Get a base64 string representation of the bid.
		b, err := f.MarshalBase64()
		if err != nil {
			t.Fatal(err)
		}

		// Use the method of response to work out the structure type.
		a, err := ResponseFromBase64(b)
		if err != nil {
			t.Fatal(err)
		}

		if n, ok := a.(*Failed); ok {
			// Verify the new instance with the signer.
			n.Seed = d
			verifyOWID(t, s, n, true)
		} else {
			t.Fatal("bid invalid")
		}
	})
	t.Run("fail", func(t *testing.T) {

		// Change the failed and then verify them to confirm that they
		// do not pass verification now the target data has changed.
		f.Host = "different.com"
		verifyOWID(t, s, f, false)
	})
	t.Run("unsigned", func(t *testing.T) {
		f.OWID = nil
		verifyOWID(t, s, f, false)
	})
}
