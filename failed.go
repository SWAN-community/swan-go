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
	"fmt"
	"owid"
)

// Failed contains details about the request that was not signed by the
// recipient.
type Failed struct {
	base
	Host  string // The domain that did not respond.
	Error string // The error message to add to the tree.
}

// FailedFromOWID returns a Failed created from the OWID payload.
func FailedFromOWID(i *owid.OWID) (*Failed, error) {
	var n Failed
	f := bytes.NewBuffer(i.Payload)
	err := n.setFromBuffer(f)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// AsByteArray returns the Failed as a byte array.
func (n *Failed) AsByteArray() ([]byte, error) {
	var f bytes.Buffer
	n.writeToBuffer(&f)
	return f.Bytes(), nil
}

func (n *Failed) writeToBuffer(f *bytes.Buffer) error {
	n.version = typeVersion
	n.structType = typeFailed
	err := n.base.writeToBuffer(f)
	if err != nil {
		return err
	}
	err = writeString(f, n.Host)
	if err != nil {
		return err
	}
	err = writeString(f, n.Error)
	if err != nil {
		return err
	}
	return nil
}

func (n *Failed) setFromBuffer(f *bytes.Buffer) error {
	err := n.base.setFromBuffer(f)
	if err != nil {
		return err
	}
	if n.structType != typeFailed {
		return fmt.Errorf(
			"Type %s not valid for %s",
			typeAsString(n.structType),
			typeAsString(typeFailed))
	}
	switch n.base.version {
	case byte(1):
		err = n.setFromBufferVersion1(f)
		break
	default:
		err = fmt.Errorf("Version '%d' not supported", n.base.version)
		break
	}
	return err
}

func (n *Failed) setFromBufferVersion1(f *bytes.Buffer) error {
	var err error
	n.Host, err = readString(f)
	if err != nil {
		return err
	}
	n.Error, err = readString(f)
	if err != nil {
		return err
	}
	return nil
}
