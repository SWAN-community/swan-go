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
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"owid"
	"swift"
	"time"
)

// Seperator used for an array of string values.
const listSeparator = "\r\n"

// handlerRawAsJSON returns the original data held in the the operation.
// Used by user interfaces to get the operations details for dispaly, or to
// continue a storage operation after time has passed waiting for the user.
// This method should never be used for passing for purposes other than for
// users editing their data.
func handlerRawAsJSON(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check caller can access.
		if s.getAccessAllowed(w, r) == false {
			returnAPIError(&s.config, w,
				errors.New("Not authorized"),
				http.StatusUnauthorized)
			return
		}

		// Decode the query string to form the byte array.
		d, err := base64.RawURLEncoding.DecodeString(r.Form.Get("data"))
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Decrypt the string with the access node.
		v, err := decryptAndDecode(s.swift, r.Host, d)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Create a map of key value pairs.
		p := make(map[string]string)

		// Unpack or copy the SWIFT key value pairs to the map.
		for _, v := range v.Pairs() {
			switch v.Key() {
			case "swid":
				// SWID does not get the OWID removed. It's is copied.
				o := getOWIDFromSWIFTPair(s, v)
				if o != nil {
					p[v.Key()] = o.AsString()
				}
				break
			case "email":
				// Email is unpacked so that the original value can be
				// displayed.
				b := unpackOWID(s, v)
				if b != nil {
					p[v.Key()] = string(b)
				}
				break
			case "pref":
				// Allow preferences are unpacked so that the original value can
				// be displayed.
				b := unpackOWID(s, v)
				if b != nil {
					p[v.Key()] = string(b)
				}
				break
			}
		}

		// If there is no valid SWID create a new one.
		if p["swid"] == "" {
			o, err := createSWID(s, r)
			if err != nil {
				returnAPIError(
					&s.config,
					w,
					err,
					http.StatusInternalServerError)
				return
			}
			p["swid"] = o.AsString()
		}

		// Set the values needed by the UIP to continue the operation.
		p["title"] = v.HTML.Title
		p["backgroundColor"] = v.HTML.BackgroundColor
		p["messageColor"] = v.HTML.MessageColor
		p["progressColor"] = v.HTML.ProgressColor
		p["message"] = v.HTML.Message
		p["returnUrl"] = v.State()[0]
		p["accessNode"] = v.State()[1]
		p["displayUserInterface"] = v.State()[2]
		p["postMessageOnComplete"] = v.State()[3]

		// Turn the map of Raw SWAN data into a JSON string.
		j, err := json.Marshal(p)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}

		// Send the JSON string.
		sendGzipJSON(s, w, j)
	}
}

// handlerDataAsJSON turns the the "data" parameter into an array of key value
// pairs where the value is encoded as an OWID using the credentials of the SWAN
// Operator.
// If the timestamp of the data provided has expired then an error is returned.
// The Email value is converted to a hashed version before being returned.
func handlerDataAsJSON(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check caller can access.
		if s.getAccessAllowed(w, r) == false {
			returnAPIError(&s.config, w,
				errors.New("Not authorized"),
				http.StatusUnauthorized)
			return
		}

		// Decode the query string to form the byte array.
		d, err := base64.RawURLEncoding.DecodeString(r.Form.Get("data"))
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Decrypt the string with the access node.
		o, err := decryptAndDecode(s.swift, r.Host, d)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Validate that the timestamp has not expired.
		if o.IsTimeStampValid() == false {
			returnAPIError(
				&s.config,
				w,
				fmt.Errorf("data expired"),
				http.StatusBadRequest)
			return
		}

		// Copy the key value pairs from SWIFT to SWAN. This is needed to
		// turn the email into a SID, and to convert the stopped domains from
		// byte arrays to a single string.
		v, err := convertPairs(s, r, o.Pairs())
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Turn the SWAN Pairs into a JSON string.
		j, err := json.Marshal(v)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}

		// Send the JSON string.
		sendGzipJSON(s, w, j)
	}
}

func sendGzipJSON(s *services, w http.ResponseWriter, j []byte) {
	g := gzip.NewWriter(w)
	defer g.Close()
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	_, err := g.Write(j)
	if err != nil {
		returnAPIError(&s.config, w, err, http.StatusInternalServerError)
	}
}

// unpackOWID return the payload from the OWID value, or nil if the OWID is not
// valid or present.
func unpackOWID(s *services, p *swift.Pair) []byte {
	o := getOWIDFromSWIFTPair(s, p)
	if o != nil {
		return o.Payload
	}
	return nil
}

// getOWIDFromSWIFTPair return the OWID, or nil if the OWID is not valid or
// present.
func getOWIDFromSWIFTPair(s *services, p *swift.Pair) *owid.OWID {
	if len(p.Values()) == 1 && len(p.Values()[0]) > 0 {
		o, err := owid.FromByteArray(p.Values()[0])
		if err != nil {
			if s.config.Debug == true {
				log.Println(err.Error())
			}
		} else {
			return o
		}
	}
	return nil
}

// verifyOWID if debug is enabled then the OWID is verified before being
// returned to the
func verifyOWID(s *services, v []byte) error {
	if s.config.Debug {
		o, err := owid.FromByteArray(v)
		if err != nil {
			return err
		}
		b, err := o.Verify(s.config.Scheme)
		if err != nil {
			return err
		}
		if b == false {
			return fmt.Errorf("OWID failed verification")
		}
	}
	return nil
}

// Copy the SWIFT results to the SWAN pairs. If the key is the email then it
// will be converted to a SID. An error is returned if the SWIFT results are
// not usable.
func convertPairs(
	s *services,
	r *http.Request,
	p []*swift.Pair) ([]*Pair, error) {
	var err error
	t := time.Now().UTC().Add(time.Second * s.config.ValueTimeout)
	w := make([]*Pair, len(p))
	for i, v := range p {
		switch v.Key() {
		case "email":
			err = verifyOWID(s, v.Values()[0])
			if err != nil {
				return nil, err
			}
			w[i], err = getSID(s, r, v)
			if err != nil {
				return nil, err
			}
			break
		case "pref":
			err = verifyOWID(s, v.Values()[0])
			if err != nil {
				return nil, err
			}
			w[i] = copyValue(v)
			break
		case "swid":
			err = verifyOWID(s, v.Values()[0])
			if err != nil {
				return nil, err
			}
			w[i] = copyValue(v)
			break
		case "stop":
			w[i], err = getStopped(v)
			if err != nil {
				return nil, err
			}
			break
		default:
			w[i] = copyValue(v)
			break
		}
		w[i].Expires = t
	}
	return w, nil
}

// Converts the array of stopped values into a single string seperated by the
// listSeparator.
func getStopped(p *swift.Pair) (*Pair, error) {
	var f bytes.Buffer
	for i, v := range p.Values() {
		_, err := f.Write(v)
		if err != nil {
			return nil, err
		}
		if i < len(p.Value())-1 {
			f.Write([]byte(listSeparator))
		}
	}
	return NewPairFromSWIFT(p, f.Bytes()), nil
}

func copyValue(p *swift.Pair) *Pair {
	return NewPairFromSWIFT(p, p.Values()[0])
}

// getSID turns the email address that is contained in the Value OWID into
// a hashed version in a new OWID with this SWAN Operator as the creator.
func getSID(s *services, r *http.Request, p *swift.Pair) (*Pair, error) {
	var v Pair
	v.Key = "sid"
	v.Created = p.Created()
	v.Expires = p.Expires()
	if len(p.Values()[0]) > 0 {
		sid, err := createSID(p.Values()[0])
		if err != nil {
			return nil, err
		}
		o, err := createOWID(s, r, sid)
		if err != nil {
			return nil, err
		}
		v.Value, err = o.AsByteArray()
		if err != nil {
			return nil, err
		}
	}
	return &v, nil
}

func createOWID(s *services, r *http.Request, v []byte) (*owid.OWID, error) {

	// Get the creator associated with this SWAN domain.
	c, err := s.owid.GetCreator(r.Host)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf(
			"No creator for '%s'. Use http[s]://%s/owid/register to setup "+
				"domain.",
			r.Host,
			r.Host)
	}

	// Create and sign the OWID.
	o, err := c.CreateOWIDandSign(v)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// TODO : What hashing algorithm do we want to use to turn email address into
// hashes?
func createSID(email []byte) ([]byte, error) {
	o, err := owid.FromByteArray(email)
	if err != nil {
		return nil, err
	}
	hasher := sha1.New()
	hasher.Write(o.Payload)
	return hasher.Sum(nil), nil
}
