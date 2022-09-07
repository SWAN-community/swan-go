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

func TestPreferences(t *testing.T) {
	s := owid.NewTestDefaultSigner(t)
	t.Run("pass", func(t *testing.T) {

		// Create the new preferences with the flag set to true.
		p, err := NewPreferences(s, true)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the preferences and check that they pass.
		r, err := s.Verify(p.OWID)
		if err != nil {
			t.Fatal(err)
		}
		if !r {
			t.Fatal("Expected verification to pass")
		}
	})
	t.Run("fail", func(t *testing.T) {

		// Create the new preferences with the flag set to true.
		p, err := NewPreferences(s, true)
		if err != nil {
			t.Fatal(err)
		}

		// Change the preferences and then verify them to confirm that they
		// do not pass verification now the target data has changed.
		p.Data.UseBrowsingForPersonalization = false
		r, err := s.Verify(p.OWID)
		if err != nil {
			t.Fatal(err)
		}
		if r {
			t.Fatal("Expected verification to fail")
		}
	})
	t.Run("json", func(t *testing.T) {

		// Create the new preferences with the flag set to true.
		p, err := NewPreferences(s, true)
		if err != nil {
			t.Fatal(err)
		}

		// Get a JSON representation of the preferences.
		j, err := json.Marshal(p)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new instance of the preferences from the JSON.
		n, err := PreferencesFromJson(j)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the new instance with the signer.
		r, err := s.Verify(n.OWID)
		if err != nil {
			t.Fatal(err)
		}
		if !r {
			t.Fatal("Expected verification to pass")
		}
	})
	t.Run("binary", func(t *testing.T) {

		// Create the new preferences with the flag set to true.
		p, err := NewPreferences(s, true)
		if err != nil {
			t.Fatal(err)
		}

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
		r, err := s.Verify(n.OWID)
		if err != nil {
			t.Fatal(err)
		}
		if !r {
			t.Fatal("Expected verification to pass")
		}
	})
}
