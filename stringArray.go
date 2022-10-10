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
	"github.com/SWAN-community/owid-go"
	"github.com/SWAN-community/swift-go"
)

// StringArray used to store arrays of entries in the SWAN model. For example
// advert identifiers that have been stopped.
type StringArray struct {
	Writeable
	Value []string
}

// getOWID always returns nil. Provided to satisfy the Entry interface.
func (s *StringArray) GetOWID() *owid.OWID { return nil }

// UnmarshalSwift to turn a SWIFT pair into a string array.
func (s *StringArray) UnmarshalSwift(p *swift.Pair) error {
	s.Value = make([]string, 0, len(p.Values()))
	for _, v := range p.Values() {
		if len(v) > 0 {
			s.Value = append(s.Value, string(v))
		}
	}
	if s.Cookie == nil {
		s.Cookie = &Cookie{}
	}
	return s.Cookie.UnmarshalSwiftValidity(p)
}
