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
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"owid"
	"strings"
	"swift"
	"time"
)

// Seperator used for an array of string values.
const listSeparator = " "

// The time format to use when adding the validation time to the response.
const ValidationTimeFormat = "2006-01-02T15:04:05Z07:00"

// handlerDecryptRawAsJSON returns the original data held in the the operation.
// Used by user interfaces to get the operations details for dispaly, or to
// continue a storage operation after time has passed waiting for the user.
// This method should never be used for passing for purposes other than for
// users editing their data.
func handlerDecryptRawAsJSON(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check caller can access.
		if s.getAccessAllowed(w, r) == false {
			return
		}

		// Get the SWIFT results from the request.
		o := getSWIFTResults(s, w, r)
		if o == nil {
			return
		}

		// Create a map of key value pairs.
		p := make(map[string]interface{})

		// Unpack or copy the SWIFT key value pairs to the map.
		for _, v := range o.Pairs() {
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
			case "salt":
				// Salt is unpacked so that the email hash can be preserved.
				b := unpackOWID(s, v)
				if b != nil {
					p[v.Key()] = string(b)
				}
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
		if p["swid"] == nil {
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
		p["title"] = o.HTML.Title
		p["backgroundColor"] = o.HTML.BackgroundColor
		p["messageColor"] = o.HTML.MessageColor
		p["progressColor"] = o.HTML.ProgressColor
		p["message"] = o.HTML.Message
		p["state"] = o.State()

		// Turn the map of Raw SWAN data into a JSON string.
		j, err := json.Marshal(p)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}

		// Send the JSON string.
		sendGzipJSON(s, w, r, j)
	}
}

// handlerDecryptAsJSON turns the the "encrypted" parameter into an array of key
// value pairs where the value is encoded as an OWID using the credentials of
// the SWAN Operator.
// If the timestamp of the data provided has expired then an error is returned.
// The Email value is converted to a hashed version before being returned.
func handlerDecryptAsJSON(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check caller can access.
		if s.getAccessAllowed(w, r) == false {
			return
		}

		// Get the SWIFT results from the request.
		o := getSWIFTResults(s, w, r)
		if o == nil {
			return
		}

		// Validate that the timestamp has not expired.
		if o.IsTimeStampValid() == false {
			returnAPIError(
				&s.config,
				w,
				fmt.Errorf("data expired and can no longer be used"),
				http.StatusBadRequest)
			return
		}

		// Copy the key value pairs from SWIFT to SWAN. This is needed to
		// turn the email into a SID, and to convert the stopped domains from
		// byte arrays to a single string.
		v, err := convertPairs(s, r, o.Map())
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
		sendGzipJSON(s, w, r, j)
	}
}

// Check that the encrypted parameter is present and if so decodes and decrypts
// it to return the SWIFT results. If there is an error then the method will be
// responsible for handling the response.
func getSWIFTResults(
	s *services,
	w http.ResponseWriter,
	r *http.Request) *swift.Results {

	// Validate that the encrypted parameter is present.
	v := r.Form.Get("encrypted")
	if v == "" {
		returnAPIError(
			&s.config,
			w,
			fmt.Errorf("Missing 'encrypted' parameter"),
			http.StatusBadRequest)
		return nil
	}

	// Decode the query string to form the byte array.
	d, err := base64.RawURLEncoding.DecodeString(v)
	if err != nil {
		returnAPIError(&s.config, w, err, http.StatusBadRequest)
		return nil
	}

	// Decrypt the string with the access node.
	o, err := decryptAndDecode(s.swift, r.Host, d)
	if err != nil {
		returnAPIError(&s.config, w, err, http.StatusBadRequest)
		return nil
	}

	return o
}

// sendGzipJSON responds with the JSON payload provided. If debug is enabled
// then the response is set to the logger.
func sendGzipJSON(
	s *services,
	w http.ResponseWriter,
	r *http.Request,
	j []byte) {
	if s.config.Debug {
		log.Println(string(j))
	}
	sendResponse(s, w, "application/json", j)
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

// verifyOWID confirms that the OWID provided has a valid signature.
func verifyOWID(s *services, o *owid.OWID) error {
	b, err := o.Verify(s.config.Scheme)
	if err != nil {
		return err
	}
	if b == false {
		return fmt.Errorf("OWID failed verification")
	}
	return nil
}

// verifyOWIDIfDebug confirms that the OWID byte array provided has a valid
// signature only if debug mode is enabled.
func verifyOWIDIfDebug(s *services, v []byte) error {
	if s.config.Debug {
		o, err := owid.FromByteArray(v)
		if err != nil {
			return err
		}
		return verifyOWID(s, o)
	}
	return nil
}

// Copy the SWIFT results to the SWAN pairs. If the key is the email then it
// will be converted to a SID. An additional pair is written to contain the
// validation time in UTC. An error is returned if the SWIFT results are
// not usable.
func convertPairs(
	s *services,
	r *http.Request,
	p map[string]*swift.Pair) ([]*Pair, error) {
	var err error
	var m time.Time
	w := make([]*Pair, 0, len(p)+1)

	for _, v := range p {
		// Turn the raw SWAN data into the SWAN data ready for readonly use.
		switch v.Key() {
		case "email":
			n := p["salt"]
			if len(v.Values()) > 0 && len(n.Values()) > 0 {
				// verify email
				err = verifyOWIDIfDebug(s, v.Values()[0])
				if err != nil {
					return nil, err
				}
				// verify salt
				err = verifyOWIDIfDebug(s, n.Values()[0])
				if err != nil {
					return nil, err
				}
				s, err := getSID(s, r, v, n)
				if err != nil {
					return nil, err
				}
				w = append(w, s)
			}
			break
		case "salt":
			// Don't do anything with salt as we have used it when
			// creating the SID.
			break
		case "pref":
			if len(v.Values()) > 0 {
				err = verifyOWIDIfDebug(s, v.Values()[0])
				if err != nil {
					return nil, err
				}
				w = append(w, copyValue(v))
			}
			break
		case "swid":
			if len(v.Values()) > 0 {
				err = verifyOWIDIfDebug(s, v.Values()[0])
				if err != nil {
					return nil, err
				}
				w = append(w, copyValue(v))
			}
			break
		case "stop":
			s, err := getStopped(v)
			if err != nil {
				return nil, err
			}
			w = append(w, s)
			break
		default:
			w = append(w, copyValue(v))
			break
		}
	}

	// Find the expiry date furthest in the future. This will be used to set the
	// val pair to indicate the caller when they should recheck the network.
	for _, v := range w {
		if m.Before(v.Expires) {
			m = v.Expires
		}
	}

	// Add a final pair to indicate when the caller should revalidate the
	// SWAN data with the network. This is recommended for the caller, but not
	// compulsory.
	t := time.Now().UTC()
	e := t.Add(s.config.RevalidateSecondsDuration()).Format(
		ValidationTimeFormat)
	w = append(w, &Pair{
		Key:     "val",
		Created: t,
		Expires: m,
		Value:   e,
	})
	return w, nil
}

// Converts the array of stopped values into a single string seperated by the
// listSeparator.
func getStopped(p *swift.Pair) (*Pair, error) {
	s := make([]string, 0, len(p.Values()))
	for _, v := range p.Values() {
		if len(v) > 0 {
			s = append(s, string(v))
		}
	}
	return NewPairFromSWIFT(p, strings.Join(s, listSeparator)), nil
}

// copyValue turns the SWIFT pair into a SWAN pair taking the first value and
// base 64 encoding it as a string.
func copyValue(p *swift.Pair) *Pair {
	return NewPairFromSWIFT(
		p,
		base64.RawStdEncoding.EncodeToString(p.Values()[0]))
}

// getSID turns the email address that is contained in the Value OWID into
// a hashed version in a new OWID with this SWAN Operator as the creator.
func getSID(s *services, r *http.Request, p *swift.Pair, n *swift.Pair) (*Pair, error) {
	v := &Pair{
		Key:     "sid",
		Created: p.Created(),
		Expires: p.Expires(),
	}
	if len(p.Values()[0]) > 0 && len(n.Values()[0]) > 0 {
		sid, err := createSID(p.Values()[0], n.Values()[0])
		if err != nil {
			return nil, err
		}
		o, err := createOWID(s, r, sid)
		if err != nil {
			return nil, err
		}
		v.Value, err = o.AsBase64()
		if err != nil {
			return nil, err
		}
	}
	return v, nil
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

// Create the SID by salting the email address and creating an sha256 hashes
func createSID(email []byte, salt []byte) ([]byte, error) {
	o1, err := owid.FromByteArray(email)
	if err != nil {
		return nil, err
	}
	o2, err := owid.FromByteArray(salt)
	if err != nil {
		return nil, err
	}
	hasher := sha256.New()
	hasher.Write(append(o1.Payload, o2.Payload...))
	return hasher.Sum(nil), nil
}
