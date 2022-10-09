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
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/SWAN-community/owid-go"
	"github.com/SWAN-community/swift-go"
	"github.com/google/uuid"
)

func TestResponse(t *testing.T) {

	// Setup the default test signer.
	g := owid.NewTestDefaultSigner(t)

	// Create the data to use with the test.
	rid, err := NewIdentifier(g, "paf_browser_id", uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	ridA, err := rid.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	email, err := NewEmail(g, testEmail)
	if err != nil {
		t.Fatal(err)
	}
	emailA, err := email.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	salt, err := NewSaltFromString(g, testSalt)
	if err != nil {
		t.Fatal(err)
	}
	saltA, err := salt.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	pref, err := NewPreferences(g, true)
	if err != nil {
		t.Fatal(err)
	}
	prefA, err := pref.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("pass", func(t *testing.T) {
		m := &ModelResponse{}
		r := &swift.Results{}
		responseAddPair(r, "rid", [][]byte{ridA})
		responseAddPair(r, "pref", [][]byte{prefA})
		responseAddPair(r, "salt", [][]byte{saltA})
		responseAddPair(r, "email", [][]byte{emailA})
		err := m.UnmarshalSwift(r)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(m.RID.OWID.Signature, rid.OWID.Signature) {
			t.Fatal("rid")
		}
		if !bytes.Equal(m.Pref.OWID.Signature, pref.OWID.Signature) {
			t.Fatal("pref")
		}
		if !bytes.Equal(m.Email.OWID.Signature, email.OWID.Signature) {
			t.Fatal("email")
		}
		if !bytes.Equal(m.Salt.OWID.Signature, salt.OWID.Signature) {
			t.Fatal("salt")
		}
		err = m.SetValidity(60)
		if err != nil {
			t.Fatal(err)
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))
	})
}

func responseAddPair(r *swift.Results, key string, value [][]byte) {
	t := time.Now().UTC()
	p := &swift.Pair{}
	p.SetKey(key)
	p.SetValues(value)
	p.SetCreated(t)
	p.SetExpires(t.Add(time.Hour))
	r.Pairs = append(r.Pairs, p)
}
