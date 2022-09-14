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
		if v != u.String() {
			t.Fatal()
		}
	})
	t.Run("pass", func(t *testing.T) {

		// Verify the identifier and check that they pass.
		verifyBase(t, s, &i.Base, true)
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
		verifyBase(t, s, &n.Base, true)
	})
	t.Run("json", func(t *testing.T) {

		// Get a JSON representation of the identifier.
		j, err := json.Marshal(i)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(j))

		// Create a new instance of the identifier from the JSON.
		n, err := IdentifierFromJson(j)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyBase(t, s, &n.Base, true)
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
		verifyBase(t, s, &n.Base, true)
	})
	t.Run("cookie", func(t *testing.T) {

		// Create a cookie pair and verify the correct result is returned.
		p, err := NewPairFromField("id", i)
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
		v1, err := IdentifierUnmarshalBase64([]byte(n.Value))
		if err != nil {
			t.Fatal(err)
		}
		verifyBase(t, s, &v1.Base, true)

		// Verify that the data is correct when passed out to a field.
		var v2 Identifier
		err = p.UnmarshalBase64(&v2)
		if err != nil {
			t.Fatal(err)
		}
		verifyBase(t, s, &v2.Base, true)
	})
	t.Run("fail", func(t *testing.T) {

		// Change the identifier and then verify them to confirm that they
		// do not pass verification now the target data has changed.
		i.IdType += " "
		verifyBase(t, s, &i.Base, false)
	})
}
