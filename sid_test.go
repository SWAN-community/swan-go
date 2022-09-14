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
	"testing"

	"github.com/SWAN-community/owid-go"
)

func TestSID(t *testing.T) {
	s := owid.NewTestDefaultSigner(t)

	// Create the base SID.
	d := testCreateSID(t, s)

	t.Run("pass", func(t *testing.T) {

		// Verify the SID and check that they pass.
		verifyOWID(t, s, d.GetOWID(), true)
	})
	t.Run("same", func(t *testing.T) {

		// Verify that another SID with the same input values results in the
		// same byte array.
		n := testCreateSID(t, s)
		testCompareByteArray(t, d, n)
	})
}

func testCompareByteArray(t *testing.T, a *ByteArray, b *ByteArray) {
	if len(a.Data) != len(b.Data) {
		t.Fatal("byte array length mismatch")
	}
	for i := 0; i < len(a.Data); i++ {
		if a.Data[i] != b.Data[i] {
			t.Fatalf("byte array difference at '%d'", i)
		}
	}
}

func testCreateSID(t *testing.T, s *owid.Signer) *ByteArray {
	// Create the new email.
	e, err := NewEmail(s, "email@example.com")
	if err != nil {
		t.Fatal(err)
	}

	// Create the new salt with the from string method.
	a, err := NewSaltFromString(s, "1234")
	if err != nil {
		t.Fatal(err)
	}

	// Create the SID.
	d, err := NewSID(s, e, a)
	if err != nil {
		t.Fatal(err)
	}

	return d
}
