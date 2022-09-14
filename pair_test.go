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

func TestPair(t *testing.T) {
	pairByteArray := []byte{1, 2, 3, 4}

	s := owid.NewTestDefaultSigner(t)

	// Create the new byte array.
	a, err := NewByteArray(s, pairByteArray)
	if err != nil {
		t.Fatal(err)
	}

	// Create the pair.
	p, err := NewPairFromField("test", a)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("base64 good", func(t *testing.T) {

		// Unmarshal the value to the byte array.
		var c ByteArray
		err = p.UnmarshalBase64(&c)
		if err != nil {
			t.Fatal(err)
		}

		// Check the byte array is valid.
		v, err := s.Verify(c.OWID)
		if err != nil {
			t.Fatal(err)
		}
		if !v {
			t.Fatal("verification should pass")
		}
	})
	t.Run("base64 bad", func(t *testing.T) {

		// Corrupt the base 64 array.
		p.Value = " " + p.Value

		// Unmarshal the value to the byte array.
		var c ByteArray
		err = p.UnmarshalBase64(&c)
		if err == nil {
			t.Fatal("corrupt base 64 data should result in error")
		}
	})
}
