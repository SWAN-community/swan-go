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

// Empty contains nothing. Used for most OWIDs that just sign the root and
// themselves.
type Empty struct {
	base
}

// EmptyFromOWID returns an Empty created from the OWID payload.
func EmptyFromOWID(o *owid.OWID) (*Bid, error) {
	var b Bid
	f := bytes.NewBuffer(o.Payload)
	err := b.setFromBuffer(f)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// AsByteArray returns the Empty as a byte array.
func (e *Empty) AsByteArray() ([]byte, error) {
	var f bytes.Buffer
	e.writeToBuffer(&f)
	return f.Bytes(), nil
}

func (e *Empty) writeToBuffer(f *bytes.Buffer) error {
	e.version = typeVersion
	e.structType = typeEmpty
	return e.base.writeToBuffer(f)
}

func (e *Empty) setFromBuffer(f *bytes.Buffer) error {
	return e.base.setFromBuffer(f)
}
