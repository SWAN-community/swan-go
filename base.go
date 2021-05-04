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
	"owid"
)

// Version used for persisting SWAN data structures.
const typeVersion byte = 1

// Type structures that include base.
const (
	typeBid    byte = iota
	typeID     byte = iota
	typeFailed byte = iota
	typeEmpty  byte = iota
)

// base used as the first fields of any SWAN data structure.
type base struct {
	version    byte // Used to indicate the version encoding of the type.
	structType byte // Used to indicate the type of struct that follows.
}

// FromOWID returns a point to a structure of the SWAN type contained in the
// OWID.
func FromOWID(o *owid.OWID) (interface{}, error) {
	var b base
	f := bytes.NewBuffer(o.Payload)
	err := b.setFromBuffer(f)
	if err != nil {
		return nil, err
	}
	var i interface{}
	switch b.structType {
	case typeBid:
		var v Bid
		v.base = b
		_ = v.setFromBufferVersion1(f)
		i = &v
	case typeID:
		var v ID
		v.base = b
		_ = v.setFromBufferVersion1(f)
		i = &v
	case typeFailed:
		var v Failed
		v.base = b
		_ = v.setFromBufferVersion1(f)
		i = &v
	case typeEmpty:
		var v Empty
		v.base = b
		_ = v.setFromBuffer(f)
		i = &v
	}
	return i, nil
}

// FromNode returns a point to a structure of the SWAN type contained in the
// Node's OWID.
func FromNode(n *owid.Node) (interface{}, error) {
	o, err := n.GetOWID()
	if err != nil {
		return nil, err
	}
	return FromOWID(o)
}

func (b *base) writeToBuffer(f *bytes.Buffer) error {
	err := writeByte(f, b.version)
	if err != nil {
		return err
	}
	err = writeByte(f, b.structType)
	if err != nil {
		return err
	}
	return nil
}

func (b *base) setFromBuffer(f *bytes.Buffer) error {
	var err error
	b.version, err = readByte(f)
	if err != nil {
		return err
	}
	b.structType, err = readByte(f)
	if err != nil {
		return err
	}
	return nil
}

func typeAsString(b byte) string {
	switch b {
	case typeBid:
		return "Bid"
	case typeID:
		return "ID"
	case typeFailed:
		return "Failed"
	case typeEmpty:
		return "Empty"
	default:
		return "Unknown"
	}
}
