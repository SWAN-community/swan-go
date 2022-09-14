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
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/SWAN-community/owid-go"
)

func TestByteArray(t *testing.T) {
	testByteArray := []byte{1, 2, 3, 4}
	s := owid.NewTestDefaultSigner(t)

	// Create the new byte array.
	a, err := NewByteArray(s, testByteArray)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("printable", func(t *testing.T) {
		v := a.AsPrintable()
		if v != hex.EncodeToString(testByteArray) {
			t.Fatal()
		}
	})
	t.Run("pass", func(t *testing.T) {

		// Verify the byte array and check that they pass.
		verifyBase(t, s, &a.Base, true)
	})
	t.Run("base64", func(t *testing.T) {

		// Get a base64 string representation of the byte array.
		b, err := a.MarshalBase64()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the byte array from the base64 string.
		n, err := ByteArrayUnmarshalBase64(b)
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
		p, err := NewPairFromField("byteArray", a)
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
		v1, err := ByteArrayUnmarshalBase64([]byte(n.Value))
		if err != nil {
			t.Fatal(err)
		}
		verifyBase(t, s, &v1.Base, true)

		// Verify that the data is correct when passed out to a field.
		var v2 ByteArray
		err = p.UnmarshalBase64(&v2)
		if err != nil {
			t.Fatal(err)
		}
		verifyBase(t, s, &v2.Base, true)
	})
	t.Run("fail", func(t *testing.T) {

		// Change the byte array and then verify them to confirm that they
		// do not pass verification now the target data has changed.
		a.Data = []byte{4, 3, 2, 1}
		verifyBase(t, s, &a.Base, false)
	})
}
