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
)

func TestSalt(t *testing.T) {
	const testSalt = "123456"
	s := owid.NewTestDefaultSigner(t)

	// Create the new salt with the from string method.
	a, err := NewSaltFromString(s, testSalt)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("printable", func(t *testing.T) {
		v := a.AsPrintable()
		if v != testSalt {
			t.Fatal(v)
		}
	})
	t.Run("pass", func(t *testing.T) {

		// Verify the salt and check that they pass.
		verifyOWID(t, s, a, true)
	})
	t.Run("base64", func(t *testing.T) {

		// Get a base64 string representation of the salt.
		b, err := a.MarshalBase64()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the salt from the base64 string.
		n, err := SaltUnmarshalBase64(b)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, n, true)
	})
	t.Run("json", func(t *testing.T) {

		// Get a JSON representation of the salt.
		j, err := json.Marshal(a)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(j))

		// Create a new instance of the salt from the JSON.
		var n Salt
		err = json.Unmarshal(j, &n)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, &n, true)
	})
	t.Run("binary", func(t *testing.T) {

		// Get a binary representation of the salt.
		b, err := a.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the salt from the binary.
		var n Salt
		err = n.UnmarshalBinary(b)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, &n, true)
	})
	t.Run("cookie", func(t *testing.T) {

		// Create a cookie pair and verify the correct result is returned.
		p, err := NewPairFromField("salt", a)
		if err != nil {
			t.Fatal(err)
		}
		c, err := p.AsCookie(s.Domain, false)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new pair from the cookie.
		n, err := NewPairFromCookie(c)
		if err != nil {
			t.Fatal(err)
		}

		// Verify that the type is correct.
		v, err := SaltUnmarshalBase64([]byte(n.Value))
		if err != nil {
			t.Fatal(err)
		}
		verifyOWID(t, s, v, true)
	})
	t.Run("fail", func(t *testing.T) {

		// Change the salt and then verify them to confirm that they
		// do not pass verification now the target data has changed.
		a.Salt = []byte{2}
		verifyOWID(t, s, a, false)
	})
}
