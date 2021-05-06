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
	"io/ioutil"
	"net/http"
	"net/url"
	"owid"
	"swift"
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
	ReturnUrl             *url.URL
	AccessNode            string
	Title                 string
	Message               string
	ProgressColor         string
	BackgroundColor       string
	MessageColor          string
	NodeCount             int
	DisplayUserInterface  bool
	PostMessageOnComplete bool
	UseHomeNode           bool
	JavaScript            bool
	State                 []string
}

// Update operation from a User Interface Provider where the preferences, email
// and salt have been captured. The SWID is returned from a previous call to
// swan.CreateSWID.
type Update struct {
	Operation
	SWID  string
	Pref  bool
	Email string
	Salt  []byte
}

// Fetch operation to retrieve the SWAN data.
type Fetch struct {
	Operation
}

// Stop operation to block an advert domain or identifier.
type Stop struct {
	Operation
	Host string // Advert host to block
}

// Connection stores the static details that are used when creating a new
// swan request.
type Connection struct {
	operation Operation
}

// NewConnection creates a new SWAN connection based on the operation provided.
func NewConnection(operation Operation) *Connection {
	return &Connection{operation: operation}
}

// NewFetch creates a new fetch operation using the default in the connection.
// request http request from a web browser
// returnUrl return URL after the operation completes
func (c *Connection) NewFetch(
	request *http.Request,
	returnUrl *url.URL) *Fetch {
	f := Fetch{}
	f.Operation = c.operation
	f.Request = request
	f.ReturnUrl = returnUrl
	return &f
}

// NewUpdate creates a new fetch operation using the default in the connection.
// request http request from a web browser
// returnUrl return URL after the operation completes
func (c *Connection) NewUpdate(
	request *http.Request,
	returnUrl *url.URL) *Update {
	p := Update{}
	p.Operation = c.operation
	p.Request = request
	p.ReturnUrl = returnUrl
	return &p
}

// NewStop creates a new stop operation using the default in the connection.
// request http request from a web browser
// returnUrl return URL after the operation completes
// host associated with the advert to stop
func (c *Connection) NewStop(
	request *http.Request,
	returnUrl *url.URL,
	host string) *Stop {
	s := Stop{}
	s.Operation = c.operation
	s.Request = request
	s.ReturnUrl = returnUrl
	s.Host = host
	return &s
}

// NewClient creates a new request.
// request http request from a web browser
func (c *Connection) NewClient(request *http.Request) *Client {
	l := Client{}
	l.SWAN = c.operation.SWAN
	l.Request = request
	return &l
}

// NewDecrypt creates a new decrypt request using the default in the
// connection.
// encrypted the base 64 encoded data to be decrypted
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
// URL string that the web browser should be directed to.
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
// creator used to create the OWIDs for the data in the Update structure
func (u *Update) GetURL(creator *owid.Creator) (string, *Error) {
	q := url.Values{}
	err := u.setData(&q, creator)
	if err != nil {
		return "", &Error{Err: err}
	}
	return requestAsString(&u.SWAN, "update", q)
}

// GetValues returns the values that can be used to configure a web browser with
// the information contained in the Update operation. Ensure the access key is
// not included in the resulting values.
// creator used to create the OWIDs for the data in the Update structure
func (u *Update) GetValues(
	creator *owid.Creator) (url.Values, error) {
	q := url.Values{}
	err := u.setData(&q, creator)
	if err != nil {
		return nil, err
	}
	q.Del("accessKey")
	q.Del("swid")
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

// Decrypt returns key value pairs for the data contained in the encrypted
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

// CreateSWID returns a new SWID in OWID format.
func (c *Connection) CreateSWID() (*owid.OWID, *Error) {
	return c.NewSWAN().createSWID()
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

func (s *SWAN) createSWID() (*owid.OWID, *Error) {
	b, se := requestAsByteArray(s, "create-swid", url.Values{})
	if se != nil {
		return nil, se
	}
	o, err := owid.FromByteArray(b)
	if err != nil {
		return nil, &Error{Err: err}
	}
	return o, nil
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
	b, err := ioutil.ReadAll(res.Body)
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
		return fmt.Errorf("Request required")
	}
	swift.SetHomeNodeHeaders(c.Request, q)
	return nil
}

func (e *Decrypt) setData(q *url.Values) error {
	if e.Encrypted == "" {
		return fmt.Errorf("Encrypted required")
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
		return fmt.Errorf("host required")
	}
	q.Set("host", s.Host)
	return nil
}

func (o *Operation) setData(q *url.Values) error {
	err := o.Client.setData(q)
	if err != nil {
		return err
	}
	if o.ReturnUrl == nil {
		return fmt.Errorf("ReturnURL required")
	}
	q.Set("returnUrl", o.ReturnUrl.String())
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

func (u *Update) setData(q *url.Values, c *owid.Creator) error {
	err := u.Operation.setData(q)
	if err != nil {
		return err
	}
	if u.Pref {
		err = setSWANData(c, q, "pref", []byte("on"))
	} else {
		err = setSWANData(c, q, "pref", []byte("off"))
	}
	if err != nil {
		return err
	}
	err = setSWANData(c, q, "email", []byte(u.Email))
	if err != nil {
		return err
	}
	err = setSWANData(c, q, "salt", u.Salt)
	if err != nil {
		return err
	}
	q.Set("swid", u.SWID)
	return nil
}

// setSWANData uses the creator to turn the value v into an OWID before setting
// that OWID as a base 64 string in the query values q against the key k.
// c owid creator for the User Interface Provider
// q collection of key value pairs
// k the key for the SWAN value
// v the raw value to be used as the payload for the OWID
func setSWANData(c *owid.Creator, q *url.Values, k string, v []byte) error {
	o, err := c.CreateOWIDandSign(v)
	if err != nil {
		return err
	}
	q.Set(k, o.AsString())
	return nil
}
