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
	"fmt"
	"owid"
)

// WinningOWID gets the winning OWID.
func WinningOWID(o *owid.Node) (*owid.OWID, error) {
	w, err := WinningNode(o)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, fmt.Errorf("No winning bid")
	}
	return w.GetOWID()
}

// WinningBid gets the winning bid from the winner's Processor OWID.
func WinningBid(o *owid.Node) (*Bid, error) {
	w, err := WinningOWID(o)
	if err != nil {
		return nil, err
	}
	return BidFromOWID(w)
}

// WinningNode gets the winning Processor OWID node for the transaction.
func WinningNode(o *owid.Node) (*owid.Node, error) {
	w := o.Find(func(n *owid.Node) bool {
		v, ok := n.Value.(float64)
		return ok && v >= 0
	})
	if w != nil {
		for w != nil {
			i, ok := w.Value.(float64)
			if ok && int(i) < len(w.Children) && int(i) >= 0 {
				w = w.Children[int(i)]
			} else {
				break
			}
		}
	}
	return w, nil
}
