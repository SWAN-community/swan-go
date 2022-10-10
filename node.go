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

import "github.com/SWAN-community/owid-go"

// Node that contains a response and an optional array of child response.
type Node struct {
	Bid      *Bid    `json:"bid,omitempty"`
	Failed   *Failed `json:"failed,omitempty"`
	Empty    *Empty  `json:"empty,omitempty"`
	Children []Node  `json:"children,omitempty"`
}

func (n *Node) SetEmpty(signer *owid.Signer, seed *Seed) error {
	var err error
	n.Empty, err = NewEmpty(signer, seed)
	return err
}

func (n *Node) SetBid(
	signer *owid.Signer,
	seed *Seed,
	mediaUrl string,
	advertiserUrl string) error {
	var err error
	n.Bid, err = NewBid(signer, seed, mediaUrl, advertiserUrl)
	return err
}

func (n *Node) SetFailed(
	signer *owid.Signer,
	seed *Seed,
	host string,
	message string) error {
	var err error
	n.Failed, err = NewFailed(signer, seed, host, message)
	return err
}

// FindFirst returns the first node in the tree under this one that when passed
// the function provided returns true. Nil is returned if no nodes fulfill the
// criteria.
func (n *Node) FindFirst(f func(*Node) (bool, error)) (*Node, error) {
	r, err := f(n)
	if err != nil {
		return nil, err
	}
	if r {
		return n, nil
	}
	for _, i := range n.Children {
		c, err := i.FindFirst(f)
		if err != nil {
			return nil, err
		}
		if c != nil {
			return c, nil
		}
	}
	return nil, nil
}

// AddMatching adds any matching nodes to the array passed into the method.
func (n *Node) AddMatching(m *[]*Node, f func(*Node) (bool, error)) error {
	r, err := f(n)
	if err != nil {
		return err
	}
	if r {
		*m = append(*m, n)
	}
	for _, i := range n.Children {
		err := i.AddMatching(m, f)
		if err != nil {
			return err
		}
	}
	return nil
}
