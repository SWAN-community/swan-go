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

func TestNode(t *testing.T) {
	s := owid.NewTestDefaultSigner(t)
	s.Domain = "receiver"
	d := createSeedTest(t, s)
	t.Run("empty", func(t *testing.T) {

		// Create an empty response.
		e, err := NewEmpty(s, d)
		if err != nil {
			t.Fatal(err)
		}

		// Create the JSON for the response.
		n := Node{Empty: e}
		b, err := json.Marshal(n)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))

		// Create a new node from the JSON.
		var c Node
		err = json.Unmarshal(b, &c)
		if err != nil {
			t.Fatal(err)
		}

		// Verify that it fails when there is no seed.
		verifyOWID(t, s, c.Empty.GetOWID(), false)

		// Verify when the seed is added to the response.
		c.Empty.Seed = d
		verifyOWID(t, s, c.Empty.GetOWID(), true)
	})
	t.Run("failed", func(t *testing.T) {

		// Create a failed response.
		f, err := NewFailed(s, d, "bad.host", "no response")
		if err != nil {
			t.Fatal(err)
		}

		// Create the JSON for the response.
		n := Node{Failed: f}
		b, err := json.Marshal(n)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))

		// Create a new node from the JSON.
		var c Node
		err = json.Unmarshal(b, &c)
		if err != nil {
			t.Fatal(err)
		}

		// Verify that it fails when there is no seed.
		verifyOWID(t, s, c.Failed.GetOWID(), false)

		// Verify when the seed is added to the response.
		c.Failed.Seed = d
		verifyOWID(t, s, c.Failed.GetOWID(), true)
	})
	t.Run("bid", func(t *testing.T) {

		// Create an bid response.
		i, err := NewBid(s, d, "media.url", "advertiser.url")
		if err != nil {
			t.Fatal(err)
		}

		// Create the JSON for the response.
		n := Node{Bid: i}
		b, err := json.Marshal(n)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))

		// Create a new node from the JSON.
		var c Node
		err = json.Unmarshal(b, &c)
		if err != nil {
			t.Fatal(err)
		}

		// Verify that it fails when there is no seed.
		verifyOWID(t, s, c.Bid.GetOWID(), false)

		// Verify when the seed is added to the response.
		c.Bid.Seed = d
		verifyOWID(t, s, c.Bid.GetOWID(), true)
	})
	t.Run("tree", func(t *testing.T) {

		// Create an empty response.
		e, err := NewEmpty(s, d)
		if err != nil {
			t.Fatal(err)
		}

		// Set the value and the children.
		n := Node{Empty: e, Children: make([]Node, 3, 3)}
		err = n.Children[0].SetEmpty(s, d)
		if err != nil {
			t.Fatal(err)
		}
		err = n.Children[1].SetFailed(s, d, "bad.host", "no response")
		if err != nil {
			t.Fatal(err)
		}
		err = n.Children[2].SetBid(s, d, "media.url", "advertiser.url")
		if err != nil {
			t.Fatal(err)
		}

		// Create the JSON for the response.
		b, err := json.Marshal(n)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))

		// Create a new node from the JSON.
		var c Node
		err = json.Unmarshal(b, &c)
		if err != nil {
			t.Fatal(err)
		}

		// Verify the values are as we would expect.
		if c.Failed != nil {
			t.Fatal("root failed should be nil")
		}
		if c.Bid != nil {
			t.Fatal("root bid should be nil")
		}
		c.Empty.Seed = d
		verifyOWID(t, s, c.Empty.GetOWID(), true)

		// Verify the children with the seed.
		for _, i := range c.Children {
			if i.Bid != nil {
				i.Bid.Seed = d
				verifyOWID(t, s, i.Bid.GetOWID(), true)
			} else if i.Failed != nil {
				i.Failed.Seed = d
				verifyOWID(t, s, i.Failed.GetOWID(), true)
			} else if i.Empty != nil {
				i.Empty.Seed = d
				verifyOWID(t, s, i.Empty.GetOWID(), true)
			} else {
				t.Fatal("child has no value")
			}
		}
	})
}
