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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/SWAN-community/swift-go"
)

// SWAN is the base structure for all actions. It includes the scheme for the
// SWAN Operator URLs, the Operator domain and the access key needed by the
// SWAN Operator.
type SWAN struct {
	Scheme    string // The HTTP or HTTPS scheme to use for SWAN requests
	Operator  string // Domain name of the SWAN Operator access node
	AccessKey string // SWAN access key provided by the SWAN Operator
}

// Decrypt contains the string to be decrypted via the call to SWAN.
type Decrypt struct {
	SWAN
	Encrypted string // The encrypted string to be decrypted by SWAN
}

// Client is used for actions where a request from a web browser is available.
// It is mainly used to set the home node from the public IP address of the web
// browser.
type Client struct {
	SWAN
	Request *http.Request // The HTTP request from the web browser
}

// Operation has members for all the parameters for a storage operation
// involving a URL that is requested by the web browser.
type Operation struct {
	Client
	// The URL to return to with the encrypted data appended to it.
	ReturnUrl string
	// The access node that will be used to decrypt the result of the storage
	// operation. Defaults to the access node that started the storage
	// operation.
	AccessNode      string
	Title           string // The title of the progress UI page.
	Message         string // The text of the message in the progress UI.
	ProgressColor   string // The HTML color for the progress indicator.
	BackgroundColor string // The HTML color for the progress UI background.
	MessageColor    string // The HTML color for the message text.
	NodeCount       int    // Number of storage nodes to use for operations.
	// DisplayUserInterface true if a progress UI should be displayed during the
	// storage operation, otherwise false.
	DisplayUserInterface bool
	// PostMessageOnComplete true if at the end of the operation the resulting
	// data should be returned to the parent using JavaScript postMessage,
	// otherwise false. Default false.
	PostMessageOnComplete bool
	// UseHomeNode true if the home node can be used if it contains current
	// data. False if the SWAN network should be consulted irrespective of the
	// state of data held on the home node. Default true.
	UseHomeNode bool
	// JavaScript true if the response for storage operations should be
	// JavaScript include that will continue the operation. This feature
	// requires cookies to be sent for DOM inserted JavaScript elements. Default
	// false.
	JavaScript bool
	// Optional array of strings that can be used to pass state information to
	// the party that retrieves the results of the storage operation. For
	// example; passing information between a Publisher and User Interface
	// Provider such as a CMP in the storage operation.
	State []string
}

// Update operation from a User Interface Provider where the preferences, email
// and salt have been captured. The RID is returned from a previous call to
// swan.CreateRID.
type Update struct {
	Operation
	RID   *Identifier  // Random [browser] Id - see CreateRID
	Pref  *Preferences // Preference for marketing - see CreatePreferences
	Email *Email       // Email address - see CreateEmail
	Salt  *ByteArray   // Salt for SID - see CreateByteArray
}

// Fetch operation to retrieve the SWAN data for use with a call to Decrypt or
// DecryptRaw.
type Fetch struct {
	Operation
	Existing []*Pair // Existing SWAN data pairs
}

// Stop operation to block an advert domain or identifier.
type Stop struct {
	Operation
	Host string // Advert host to block
}

// Connection stores the static details that are used when creating a new swan
// request.
type Connection struct {
	operation Operation
}

// NewConnection creates a new SWAN connection based on the operation provided.
func NewConnection(operation Operation) *Connection {
	return &Connection{operation: operation}
}

// NewFetch creates a new fetch operation using the default in the connection.
//
// request http request from a web browser
//
// returnUrl return URL after the operation completes
//
// existing if any values already exist then use these if none are available in
// SWAN
func (c *Connection) NewFetch(
	request *http.Request,
	returnUrl string,
	existing []*Pair) *Fetch {
	f := Fetch{}
	f.Operation = c.operation
	f.Request = request
	f.ReturnUrl = returnUrl
	f.Existing = existing
	return &f
}

// NewUpdate creates a new fetch operation using the default in the connection.
//
// request http request from a web browser
//
// returnUrl return URL after the operation completes
func (c *Connection) NewUpdate(
	request *http.Request,
	returnUrl string) *Update {
	p := Update{}
	p.Operation = c.operation
	p.Request = request
	p.ReturnUrl = returnUrl
	return &p
}

// NewStop creates a new stop operation using the default in the connection.
//
// request http request from a web browser
//
// returnUrl return URL after the operation completes
//
// host associated with the advert to stop
func (c *Connection) NewStop(
	request *http.Request,
	returnUrl string,
	host string) *Stop {
	s := Stop{}
	s.Operation = c.operation
	s.Request = request
	s.ReturnUrl = returnUrl
	s.Host = host
	return &s
}

// NewClient creates a new request.
//
// request http request from a web browser
func (c *Connection) NewClient(request *http.Request) *Client {
	l := Client{}
	l.SWAN = c.operation.SWAN
	l.Request = request
	return &l
}

// NewDecrypt creates a new decrypt request using the default in the
// connection.
//
// encrypted the base 64 encoded SWAN data to be decrypted
func (c *Connection) NewDecrypt(encrypted string) *Decrypt {
	e := Decrypt{}
	e.SWAN = c.operation.SWAN
	e.Encrypted = encrypted
	return &e
}

// NewSWAN creates a new request using the default in the connection.
func (c *Connection) NewSWAN() *SWAN {
	s := c.operation.SWAN
	return &s
}

// GetURL contacts the SWAN operator domain with the access key and returns a
// URL string that the web browser should be immediately directed to.
func (f *Fetch) GetURL() (string, *Error) {
	q := url.Values{}
	err := f.setData(&q)
	if err != nil {
		return "", &Error{Err: err}
	}
	return requestAsString(&f.SWAN, "fetch", q)
}

// GetURL contacts the SWAN operator domain with the access key and returns a
// URL string that the web browser should be directed to.
func (u *Update) GetURL() (string, *Error) {
	q := url.Values{}
	err := u.setData(&q)
	if err != nil {
		return "", &Error{Err: err}
	}
	return requestAsString(&u.SWAN, "update", q)
}

// GetValues returns the values that can be used to configure a web browser with
// the information contained in the Update operation. Ensure the access key and
// other values that are specific to an operation are not included in the
// resulting values.
func (u *Update) GetValues() (url.Values, error) {
	q := url.Values{}
	err := u.setData(&q)
	if err != nil {
		return nil, err
	}
	q.Del("accessKey") // Known only to this party and must never be shared
	q.Del("rid")       // Not to be shared with other browsers
	// Used for home node operations that depend on the specific browser
	q.Del("remoteAddr")
	q.Del("X-Forwarded-For")
	return q, nil
}

// GetURL contacts the SWAN operator domain with the access key and returns a
// URL string that the web browser should be directed to.
func (s *Stop) GetURL() (string, *Error) {
	q := url.Values{}
	err := s.setData(&q)
	if err != nil {
		return "", &Error{Err: err}
	}
	return requestAsString(&s.SWAN, "stop", q)
}

// Decrypt returns SWAN key value pairs for the data contained in the encrypted
// string.
func (c *Connection) Decrypt(encrypted string) ([]*Pair, *Error) {
	return c.NewDecrypt(encrypted).decrypt()
}

// DecryptRaw returns key value pairs for the raw SWAN data contained in the
// encrypted string. Must only be used by User Interface Providers.
func (c *Connection) DecryptRaw(
	encrypted string) (map[string]interface{}, *Error) {
	return c.NewDecrypt(encrypted).decryptRaw()
}

// CreateRID returns a new RID in OWID format from the SWAN Operator. Only
// SWAN Operators can create legitimate RIDs.
func (c *Connection) CreateRID() (*Identifier, *Error) {
	return c.NewSWAN().createRID()
}

// HomeNode returns the SWAN home node associated with the web browser.
func (c *Connection) HomeNode(r *http.Request) (string, *Error) {
	return c.NewClient(r).homeNode()
}

func (c *Client) homeNode() (string, *Error) {
	q := url.Values{}
	err := c.setData(&q)
	if err != nil {
		return "", &Error{Err: err}
	}
	return requestAsString(&c.SWAN, "home-node", q)
}

func (e *Decrypt) decrypt() ([]*Pair, *Error) {
	var p []*Pair
	q := url.Values{}
	err := e.setData(&q)
	if err != nil {
		return nil, &Error{Err: err}
	}
	b, se := requestAsByteArray(&e.SWAN, "decrypt", q)
	if se != nil {
		return nil, se
	}
	err = json.Unmarshal(b, &p)
	if err != nil {
		return nil, &Error{Err: err}
	}
	return p, nil
}

func (e *Decrypt) decryptRaw() (map[string]interface{}, *Error) {
	r := make(map[string]interface{})
	q := url.Values{}
	err := e.setData(&q)
	if err != nil {
		return nil, &Error{Err: err}
	}
	b, se := requestAsByteArray(&e.SWAN, "decrypt-raw", q)
	if se != nil {
		return nil, se
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, &Error{Err: err}
	}
	return r, nil
}

func (s *SWAN) createRID() (*Identifier, *Error) {
	b, se := requestAsByteArray(s, "create-rid", url.Values{})
	if se != nil {
		return nil, se
	}
	var i Identifier
	err := i.UnmarshalBinary(b)
	if err != nil {
		return nil, &Error{Err: err}
	}
	return &i, nil
}

func requestAsByteArray(
	s *SWAN,
	a string,
	q url.Values) ([]byte, *Error) {

	// Verify the provided parameters.
	if s.Scheme == "" {
		return nil, &Error{Err: fmt.Errorf("scheme must be provided")}
	}
	if s.Operator == "" {
		return nil, &Error{Err: fmt.Errorf("operator must be provided")}
	}
	if s.AccessKey == "" {
		return nil, &Error{Err: fmt.Errorf("accessKey must be provided")}
	}

	// Construct the SWAN URL.
	var u url.URL
	u.Scheme = s.Scheme
	u.Host = s.Operator
	u.Path = "/swan/api/v1/" + a

	// Add the access key to the data.
	q.Set("accessKey", s.AccessKey)

	// Post the parameters to the SWAN url.
	res, err := http.PostForm(u.String(), q)
	if err != nil {
		return nil, &Error{Err: err}
	}

	// Read the response.
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &Error{Err: err, Response: res}
	}

	// If the status code is not OK then return the response and status code
	// as an error message.
	if res.StatusCode != http.StatusOK {
		return nil, &Error{Err: fmt.Errorf(string(b)), Response: res}
	}

	return b, nil
}

func requestAsString(
	s *SWAN,
	a string,
	q url.Values) (string, *Error) {
	b, err := requestAsByteArray(s, a, q)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *Client) setData(q *url.Values) error {
	if c.Request == nil {
		return errorMissing("Request")
	}
	swift.SetHomeNodeHeaders(c.Request, q)
	return nil
}

func (e *Decrypt) setData(q *url.Values) error {
	if e.Encrypted == "" {
		return errorMissing("Encrypted")
	}
	q.Set("encrypted", e.Encrypted)
	return nil
}

func (s *Stop) setData(q *url.Values) error {
	err := s.Operation.setData(q)
	if err != nil {
		return err
	}
	if s.Host == "" {
		return errorMissing("host")
	}
	q.Set("host", s.Host)
	return nil
}

func (o *Operation) setData(q *url.Values) error {
	err := o.Client.setData(q)
	if err != nil {
		return err
	}
	if o.ReturnUrl == "" {
		return errorMissing("ReturnURL")
	}
	_, err = url.Parse(o.ReturnUrl)
	if err != nil {
		return err
	}
	q.Set("returnUrl", o.ReturnUrl)
	if o.AccessNode != "" {
		q.Set("accessNode", o.AccessNode)
	}
	if o.Title != "" {
		q.Set("title", o.Title)
	}
	if o.Message != "" {
		q.Set("message", o.Message)
	}
	if o.ProgressColor != "" {
		q.Set("progressColor", o.ProgressColor)
	}
	if o.BackgroundColor != "" {
		q.Set("backgroundColor", o.BackgroundColor)
	}
	if o.MessageColor != "" {
		q.Set("messageColor", o.MessageColor)
	}
	if o.NodeCount != 0 {
		q.Set("nodeCount", fmt.Sprintf("%d", o.NodeCount))
	}
	q.Set("displayUserInterface", fmt.Sprintf("%t",
		o.DisplayUserInterface))
	q.Set("postMessageOnComplete", fmt.Sprintf("%t",
		o.PostMessageOnComplete))
	q.Set("useHomeNode", fmt.Sprintf("%t", o.UseHomeNode))
	q.Set("javaScript", fmt.Sprintf("%t", o.JavaScript))
	for _, s := range o.State {
		q.Add("state", s)
	}
	return nil
}

func (f *Fetch) setData(q *url.Values) error {
	err := f.Operation.setData(q)
	if err != nil {
		return err
	}
	if f.Existing != nil {
		for _, v := range f.Existing {
			if v.Key == "rid" || v.Key == "pref" {
				q.Set(v.Key, v.Value)
			}
		}
	}
	return nil
}

func (u *Update) setData(q *url.Values) error {
	var s []byte
	err := u.Operation.setData(q)
	if err != nil {
		return err
	}
	if u.RID != nil {
		s, err = u.RID.MarshalBase64()
		if err != nil {
			return err
		}
		q.Set("rid", string(s))
	}
	if u.Pref != nil {
		s, err = u.Pref.MarshalBase64()
		if err != nil {
			return err
		}
		q.Set("pref", string(s))
	}
	if u.Email != nil {
		s, err = u.Email.MarshalBase64()
		if err != nil {
			return err
		}
		q.Set("email", string(s))
	}
	if u.Salt != nil {
		s, err = u.Salt.MarshalBase64()
		if err != nil {
			return err
		}
		q.Set("salt", string(s))
	}
	return nil
}
