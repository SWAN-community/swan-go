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
	"bytes"
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/SWAN-community/owid-go"
	"github.com/SWAN-community/swift-go"
	"github.com/google/uuid"
)

// testModelResponse used for these tests
type testModelResponse struct {
	ModelResponse
	rid   []byte // RID as byte array
	pref  []byte // Pref as byte array
	salt  []byte // Salt as byte array
	email []byte // Eamil as byte array
}

const testValiditySeconds = 60

func TestSetValidity(t *testing.T) {

	// Setup the default test signer.
	g := owid.NewTestDefaultSigner(t)

	// Build the model.
	m := responseBuildModel(t, g)

	// Check that when all the entities have no cookie value the duration is
	// the content.
	t.Run("zero", func(t *testing.T) {

		// Set the validity to 60 seconds.
		m.SetValidity(testValiditySeconds)

		// Check the duration between the created and expires times is the
		// constant.
		d := m.Val.Expires.Sub(m.Val.Created)
		if math.Ceil(d.Seconds()) != testValiditySeconds {
			t.Fatalf("expected %d found %v", testValiditySeconds, d.Seconds())
		}

		// Check the cookie.
		c := m.Val.AsHttpCookie("test.host", true)
		if c.Value == "" {
			t.Fatal("validity cookie must have time value")
		}
	})

	// Check that when an entity has a cookie expiry time before the current
	t.Run("entity earlier", func(t *testing.T) {

		// Set the validity to 60 seconds.
		m.SetValidity(testValiditySeconds)

		// Change the expiring date of the RID to 1/2 the validity period.
		m.RID.GetCookie().Expires = m.Val.Expires.Add(
			-time.Second * testValiditySeconds / 2)

		// Reset the validity now there is an entity with an earlier time. This
		// result in the expires time changing.
		m.SetValidity(testValiditySeconds)

		// Check the duration between the created and expires times is the
		// constant.
		d := m.Val.Expires.Sub(m.Val.Created)
		if math.Ceil(d.Seconds()) != testValiditySeconds/2 {
			t.Fatalf("expected %d found %v", testValiditySeconds/2, d.Seconds())
		}
	})
}

func TestResponse(t *testing.T) {

	// Setup the default test signer.
	g := owid.NewTestDefaultSigner(t)

	// Build the model.
	s := responseBuildModel(t, g)

	t.Run("pass", func(t *testing.T) {
		m := &ModelResponse{}
		r := &swift.Results{}
		responseAddPair(r, "rid", [][]byte{s.rid})
		responseAddPair(r, "pref", [][]byte{s.pref})
		responseAddPair(r, "salt", [][]byte{s.salt})
		responseAddPair(r, "email", [][]byte{s.email})
		err := m.UnmarshalSwift(r)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(m.RID.OWID.Signature, s.RID.OWID.Signature) {
			t.Fatal("rid")
		}
		if !bytes.Equal(m.Pref.OWID.Signature, s.Pref.OWID.Signature) {
			t.Fatal("pref")
		}
		if !bytes.Equal(m.Email.OWID.Signature, s.Email.OWID.Signature) {
			t.Fatal("email")
		}
		if !bytes.Equal(m.Salt.OWID.Signature, s.Salt.OWID.Signature) {
			t.Fatal("salt")
		}
		err = m.SetValidity(testValiditySeconds)
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

// Create the data to use with the test.
func responseBuildModel(t *testing.T, g *owid.Signer) *testModelResponse {
	var err error
	m := &testModelResponse{}
	m.RID, err = NewIdentifier(g, "rid", uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	m.rid, err = m.RID.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	m.Email, err = NewEmail(g, testEmail)
	if err != nil {
		t.Fatal(err)
	}
	m.email, err = m.Email.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	m.Salt, err = NewSaltFromString(g, testSalt)
	if err != nil {
		t.Fatal(err)
	}
	m.salt, err = m.Salt.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	m.Pref, err = NewPreferences(g, true)
	if err != nil {
		t.Fatal(err)
	}
	m.pref, err = m.Pref.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	return m
}
