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

func TestIdentifier(t *testing.T) {
	s := owid.NewTestDefaultSigner(t)

	// Create the new identifier.
	u, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	i, err := NewIdentifier(s, "type", u)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("printable", func(t *testing.T) {
		v := i.AsPrintable()
		if v == "" {
			t.Fatal("identifier not printable")
		}
	})
	t.Run("pass", func(t *testing.T) {

		// Verify the identifier and check that they pass.
		verifyOWID(t, s, i, true)
	})
	t.Run("base64", func(t *testing.T) {

		// Get a base64 string representation of the identifier.
		b, err := i.MarshalBase64()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the identifier from the base64 string.
		n, err := IdentifierUnmarshalBase64(b)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, n, true)
		testCompareIdentifier(t, i, n)
	})
	t.Run("json", func(t *testing.T) {

		// Get a JSON representation of the identifier.
		j, err := json.Marshal(i)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(j))

		// Create a new instance of the identifier from the JSON.
		var n Identifier
		err = json.Unmarshal(j, &n)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, &n, true)
		testCompareIdentifier(t, i, &n)
	})
	t.Run("binary", func(t *testing.T) {

		// Get a binary representation of the identifier.
		b, err := i.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the identifier from the binary.
		var n Identifier
		err = n.UnmarshalBinary(b)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, &n, true)
		testCompareIdentifier(t, i, &n)
	})
	t.Run("cookie", func(t *testing.T) {

		// Create a cookie.
		c, err := i.AsHttpCookie(s.Domain, false)
		if err != nil {
			t.Fatal(err)
		}

		// Verify that the type is correct.
		v, err := IdentifierUnmarshalBase64([]byte(c.Value))
		if err != nil {
			t.Fatal(err)
		}
		verifyOWID(t, s, v, true)
		testCompareIdentifier(t, i, v)
	})
	t.Run("fail", func(t *testing.T) {

		// Change the identifier and then verify them to confirm that they
		// do not pass verification now the target data has changed.
		i.IdType += " "
		verifyOWID(t, s, i, false)
	})
}
