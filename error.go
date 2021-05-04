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
	"net/http"
)

// Error is used to pass back errors from methods that call APIs. If the
// Response member is set then the called method can use this information in
// its response. If it is not set then an internal server error can be assumed.
type Error struct {
	Err      error          // The underlying error message.
	Response *http.Response // The HTTP response that caused the error.
}

// StatusCode returns the status code of the response.
func (e *Error) StatusCode() int {
	if e.Response != nil {
		return e.Response.StatusCode
	}
	return 0
}

// Error returns the error message as a string from an HTTPError reference.
func (e *Error) Error() string {
	if e != nil && e.Err != nil {
		return e.Err.Error()
	}
	return "empty error"
}
