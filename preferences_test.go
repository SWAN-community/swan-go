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
	"github.com/SWAN-community/swift-go"
)

func TestPreferences(t *testing.T) {
	s := owid.NewTestDefaultSigner(t)

	// Create the new preferences.
	p, err := NewPreferences(s, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("printable", func(t *testing.T) {
		v := p.AsPrintable()
		if v != "{\"use_browsing_for_personalization\":true}" {
			t.Fatal(v)
		}
	})
	t.Run("pass", func(t *testing.T) {

		// Verify the preferences and check that they pass.
		verifyOWID(t, s, p, true)
	})
	t.Run("base64", func(t *testing.T) {

		// Get a base64 string representation of the preferences.
		b, err := p.MarshalBase64()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the preferences from the base64 string.
		n, err := PreferencesUnmarshalBase64(b)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, n, true)
	})
	t.Run("json", func(t *testing.T) {

		// Get a JSON representation of the preferences.
		j, err := json.Marshal(p)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(j))

		// Create a new instance of the preferences from the JSON.
		var n Preferences
		err = json.Unmarshal(j, &n)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, &n, true)
	})
	t.Run("binary", func(t *testing.T) {

		// Get a binary representation of the preferences.
		b, err := p.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the preferences from the binary.
		var n Preferences
		err = n.UnmarshalBinary(b)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, &n, true)
	})
	t.Run("swift", func(t *testing.T) {

		// Get a binary representation of the preferences.
		b, err := p.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		// Create the SWIFT pair.
		a := &swift.Pair{}
		a.SetCreated(p.GetOWID().TimeStamp)
		a.SetExpires(p.GetOWID().GetExpires(1))
		a.SetValues([][]byte{b})

		// Create a new instance of the preferences from the binary.
		var n Preferences
		err = n.UnmarshalSwift(a)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		verifyOWID(t, s, &n, true)
		verifyCookie(t, n.GetOWID(), n.Cookie, 1)
	})
	t.Run("cookie", func(t *testing.T) {

		// Create a cookie.
		c, err := p.AsHttpCookie(s.Domain, false)
		if err != nil {
			t.Fatal(err)
		}

		// Verify that the type is correct.
		v, err := PreferencesUnmarshalBase64([]byte(c.Value))
		if err != nil {
			t.Fatal(err)
		}
		verifyOWID(t, s, v, true)
	})
	t.Run("fail", func(t *testing.T) {

		// Change the preferences and then verify them to confirm that they
		// do not pass verification now the target data has changed.
		p.Data.UseBrowsingForPersonalization = false
		verifyOWID(t, s, p, false)
	})
	t.Run("unsigned", func(t *testing.T) {
		p.OWID = nil
		verifyOWID(t, s, p, false)
	})
}
