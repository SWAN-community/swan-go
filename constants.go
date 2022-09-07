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

// The version to use when creating SWAN instances.
const swanVersion byte = 1

// The character used to separate fields when building byte arrays for signing.
// See OneKey signature specification.
// https://github.com/OneKey-Network/addressability-framework/blob/main/mvp-spec/security-signatures.md
const oneKeyFieldSeparator = "\u2063"
