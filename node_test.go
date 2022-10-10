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
		verifyOWID(t, s, c.Empty, false)

		// Verify when the seed is added to the response.
		c.Empty.Seed = d
		verifyOWID(t, s, c.Empty, true)
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
		verifyOWID(t, s, c.Failed, false)

		// Verify when the seed is added to the response.
		c.Failed.Seed = d
		verifyOWID(t, s, c.Failed, true)
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
		verifyOWID(t, s, c.Bid, false)

		// Verify when the seed is added to the response.
		c.Bid.Seed = d
		verifyOWID(t, s, c.Bid, true)
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
		verifyOWID(t, s, c.Empty, true)

		// Verify the children with the seed.
		for _, i := range c.Children {
			if i.Bid != nil {
				i.Bid.Seed = d
				verifyOWID(t, s, i.Bid, true)
			} else if i.Failed != nil {
				i.Failed.Seed = d
				verifyOWID(t, s, i.Failed, true)
			} else if i.Empty != nil {
				i.Empty.Seed = d
				verifyOWID(t, s, i.Empty, true)
			} else {
				t.Fatal("child has no value")
			}
		}

		// Find the first bid.
		f, err := c.FindFirst(func(n *Node) (bool, error) {
			return n.Bid != nil, nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if f.Bid == nil {
			t.Fatal("find first failed")
		}

		// Find all the empty nodes.
		m := make([]*Node, 0, 4)
		err = c.AddMatching(&m, func(n *Node) (bool, error) {
			return n.Empty != nil, nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(m) != 2 {
			t.Fatal("add 2 matching failed")
		}

		// Find no nodes.
		m = make([]*Node, 0, 4)
		err = c.AddMatching(&m, func(n *Node) (bool, error) {
			return false, nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(m) != 0 {
			t.Fatal("add 0 matching failed")
		}
	})
}
