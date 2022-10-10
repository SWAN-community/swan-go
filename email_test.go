/* ****************************************************************************
 * Copyright 2022 51 Degrees Mobile Experts Limited (51degrees.com)
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
	"github.com/SWAN-community/swift-go"
)

const testEmail = "email@example.com"

func TestEmail(t *testing.T) {
	s := owid.NewTestDefaultSigner(t)

	// Create the new email.
	e, err := NewEmail(s, testEmail)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("printable", func(t *testing.T) {
		v := e.AsPrintable()
		if v != testEmail {
			t.Fatal()
		}
	})
	t.Run("pass", func(t *testing.T) {

		// Verify the email and check that they pass.
		verifyOWID(t, s, e, true)
	})
	t.Run("base64", func(t *testing.T) {

		// Get a base64 string representation of the email.
		b, err := e.MarshalBase64()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the email from the base64 string.
		n, err := EmailUnmarshalBase64(b)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, n, true)
	})
	t.Run("json", func(t *testing.T) {

		// Get a JSON representation of the email.
		j, err := json.Marshal(e)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(j))

		// Create a new instance of the email from the JSON.
		var n Email
		err = json.Unmarshal(j, &n)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, &n, true)
	})
	t.Run("binary", func(t *testing.T) {

		// Get a binary representation of the email.
		b, err := e.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the email from the binary.
		var n Email
		err = n.UnmarshalBinary(b)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, &n, true)
	})
	t.Run("swift", func(t *testing.T) {

		// Get a binary representation of the email.
		b, err := e.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		// Create the SWIFT pair.
		p := &swift.Pair{}
		p.SetCreated(e.GetOWID().TimeStamp)
		p.SetExpires(e.GetOWID().GetExpires(1))
		p.SetValues([][]byte{b})

		// Create a new instance of the email from the SWIFT pair.
		var n Email
		err = n.UnmarshalSwift(p)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, &n, true)
		verifyCookie(t, n.GetOWID(), n.Cookie, 1)
	})
	t.Run("cookie", func(t *testing.T) {

		// Create a cookie.
		c, err := e.AsHttpCookie(s.Domain, false)
		if err != nil {
			t.Fatal(err)
		}

		// Verify that the type is correct.
		v, err := EmailUnmarshalBase64([]byte(c.Value))
		if err != nil {
			t.Fatal(err)
		}
		verifyOWID(t, s, v, true)
	})
	t.Run("fail", func(t *testing.T) {

		// Change the email and then verify them to confirm that they
		// do not pass verification now the target data has changed.
		e.Email = "different@test.com"
		verifyOWID(t, s, e, false)
	})
	t.Run("unsigned", func(t *testing.T) {
		e.OWID = nil
		verifyOWID(t, s, e, false)
	})
}
