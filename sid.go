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
	"crypto/sha256"

	"github.com/SWAN-community/owid-go"
)

// NewSID generates the SID by hashing the salt and the email to create a sha256
// hash. If the email address is empty an empty byte array is returned.
func NewSID(signer *owid.Signer, email *Email, salt *Salt) (*Identifier, error) {
	if len(email.Email) == 0 {
		return NewIdentifierFromByteArray(signer, "sid", []byte{})
	}
	hasher := sha256.New()
	hasher.Write(append([]byte(email.Email), salt.Salt...))
	b := hasher.Sum(nil)
	return NewIdentifierFromByteArray(signer, "sid", b)
}
