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

func TestByteArray(t *testing.T) {
	s := owid.NewTestDefaultSigner(t)

	// Create the new byte array.
	a, err := NewByteArray(s, []byte{1, 2, 3, 4})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("pass", func(t *testing.T) {

		// Verify the byte array and check that they pass.
		verifyBase(t, s, &a.Base, true)
	})
	t.Run("base64", func(t *testing.T) {

		// Get a base64 string representation of the byte array.
		b, err := a.ToBase64()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the byte array from the base64 string.
		n, err := ByteArrayFromBase64(b)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyBase(t, s, &n.Base, true)
	})
	t.Run("json", func(t *testing.T) {

		// Get a JSON representation of the byte array.
		j, err := json.Marshal(a)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(j))

		// Create a new instance of the byte array from the JSON.
		n, err := ByteArrayFromJson(j)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyBase(t, s, &n.Base, true)
	})
	t.Run("binary", func(t *testing.T) {

		// Get a binary representation of the byte array.
		b, err := a.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the byte array from the binary.
		var n ByteArray
		err = n.UnmarshalBinary(b)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyBase(t, s, &n.Base, true)
	})
	t.Run("cookie", func(t *testing.T) {

		// Create a cookie pair and verify the correct result is returned.
		p := Pair{Key: "byteArray", Value: a}
		c, err := p.AsCookie(s.Domain, false)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new pair from the cookie.
		n, err := NewPairFromCookie(c, &ByteArray{})
		if err != nil {
			t.Fatal(err)
		}

		// Verify that the type is correct.
		if v, ok := n.Value.(*ByteArray); ok {
			verifyBase(t, s, &v.Base, true)
		} else {
			t.Fatal("wrong type")
		}
	})
	t.Run("fail", func(t *testing.T) {

		// Change the byte array and then verify them to confirm that they
		// do not pass verification now the target data has changed.
		a.Data = []byte{4, 3, 2, 1}
		verifyBase(t, s, &a.Base, false)
	})
}
