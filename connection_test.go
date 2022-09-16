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
	"net/http/httptest"
	"testing"

	"github.com/SWAN-community/owid-go"
	"github.com/google/uuid"
)

const testReturnUrl = "http://not.a.host"
const testOperatorDomain = "operator.swan"

func TestConnectionFetch(t *testing.T) {
	c, r := newTestConnection()
	f := c.NewFetch(r, testReturnUrl)
	d, err := f.GetURL()
	if d != "" || err == nil {
		t.Fatal("ToDo: needs host")
	}
}

func TestConnectionUpdate(t *testing.T) {
	var err error
	g := owid.NewTestDefaultSigner(t)
	c, r := newTestConnection()
	f := c.NewUpdate(r, testReturnUrl)
	f.RID, err = NewIdentifier(g, "paf_browser_id", uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	f.Email, err = NewEmail(g, testEmail)
	if err != nil {
		t.Fatal(err)
	}
	f.Salt, err = NewSaltFromString(g, testSalt)
	if err != nil {
		t.Fatal(err)
	}
	f.Pref, err = NewPreferences(g, true)
	if err != nil {
		t.Fatal(err)
	}
	d, err := f.GetURL()
	if d != "" || err == nil {
		t.Fatal("ToDo: needs host")
	}
}

func TestConnectionDecrypt(t *testing.T) {
	c, _ := newTestConnection()
	f := c.NewDecrypt("1234")
	d, err := f.Decrypt()
	if d != nil || err == nil {
		t.Fatal("ToDo: needs host")
	}
}

func TestConnectionDecryptRaw(t *testing.T) {
	c, _ := newTestConnection()
	f := c.NewDecrypt("1234")
	d, err := f.DecryptRaw()
	if d != nil || err == nil {
		t.Fatal("ToDo: needs host")
	}
}

func newTestConnection() (*Connection, *http.Request) {
	o := Operation{}
	o.Scheme = "http"
	o.ReturnUrl = testReturnUrl
	o.Operator = testOperatorDomain
	o.AccessKey = "A"
	c := NewConnection(o)
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	return c, r
}
