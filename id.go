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

import "fmt"

// ID node in a tree where the node is root, branch or leaf.
type ID struct {
	Value    string // The string value of this ID.
	Children []*ID  // An array of child ids, if any.
	parent   *ID    // The ID of the parent if set.
}

// GetParent returns the parent ID which is read only.
func (i *ID) GetParent() *ID { return i.parent }

// Find the first ID that matches the condition.
func (i *ID) Find(condition func(n *ID) bool) *ID {
	if condition(i) {
		return i
	}
	for _, c := range i.Children {
		f := c.Find(condition)
		if f != nil {
			return f
		}
	}
	return nil
}

// Fix sets the parent pointer for all descendents.
func (i *ID) Fix() {
	for _, c := range i.Children {
		c.parent = i
		c.Fix()
	}
}

// AddLeaf adds a childless Id to the Id.
func (i *ID) AddLeaf(value string) (uint32, *ID) {
	var l ID
	l.Value = value
	i.Children = append(i.Children, &l)
	return uint32(len(i.Children) - 1), &l
}

// AddChildren includes the other Ids provided in the list of children for this
// Id.
func (i *ID) AddChildren(others []*ID) {
	for _, o := range others {
		i.Children = append(i.Children, o)
	}
}

// GetNodeByIndexes returns the node at the integer indexes provided where each
// index of the array o is the level of the tree. To find the third, fourth and
// then second child of a tree the array could contain { 2, 3, 1 }.
func (i *ID) GetNodeByIndexes(o []uint32) (*ID, error) {
	l := 0
	c := i
	for l < len(o) {
		if len(c.Children) == 0 {
			return nil, fmt.Errorf("Node not found")
		}
		if o[l] >= uint32(len(c.Children)) {
			return nil, fmt.Errorf("Node not found")
		}
		c = c.Children[o[l]]
		l = l + 1
	}
	return c, nil
}
