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

// Writeable is the base structure for data that can be changed and persisted
// including identifiers, email, salt, and preferences.
type Writeable struct {
	Base
	Cookie    *Cookie `json:"-"`         // Cookie data
	Persisted bool    `json:"persisted"` // True if the value has been stored.
}

// getCookie returns the cookie instance. Used by the Entry interface.
func (m *Writeable) getCookie() *Cookie { return m.Cookie }
