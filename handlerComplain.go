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
	"fmt"
	"net/http"
	"net/url"
	"owid"
	"strings"
	"text/template"
)

var complaintSubjectTemplate = newComplaintTemplate(
	"subject",
	"Complaint: {{ .Organization }}")
var complaintBodyTemplate = newComplaintTemplate("body", `
To whom it may concern,

I believe that "{{ .Organization }}" used my personal information without a 
legal basis on '{{ .Date }}'. 

	Common Browser Identifier:	{{ .CBID }}
	Signed in Identifier:		{{ .SID }}

I provided you the following permissions for use of this data.

	Personalize Marketing: {{ .Preferences }}

You cryptographically signed this information. We therefore agree that you were
in posession of the information.

As an organization operating in '{{ .Country }}' you are bound by the following 
rules.

	{{ .DataProtectionURL }}

I would be grateful if you can respond by email to this address within 7 
working days.

Regards,

[INSERT YOU NAME]

--- DO NOT CHANGE THE TEXT BELOW THIS LINE ---
{{ .OfferID }} {{ .SWANOWID }}
--- DO NOT CHANGE THE TEXT ABOVE THIS LINE ---`)

// Complaint used to format an email template.
type Complaint struct {
	Offer             *Offer // The offer that the complaint relates to
	DataProtectionURL string
	Organization      string
	Country           string
	offerID           *owid.OWID
	swanOWID          *owid.OWID
}

// Date to use in the email template.
func (c *Complaint) Date() string {
	return c.swanOWID.Date.Format("2006-01-02")
}

// CBID to use in the email template.
func (c *Complaint) CBID() string {
	return c.Offer.CBIDAsString()
}

// SID to use in the email template.
func (c *Complaint) SID() string {
	return c.Offer.SIDAsString()
}

// Preferences string to use in the email template.
func (c *Complaint) Preferences() string {
	return c.Offer.PreferencesAsString()
}

// OfferID as a string
func (c *Complaint) OfferID() string {
	return c.offerID.AsString()
}

// SWANOWID as a string
func (c *Complaint) SWANOWID() string {
	return c.swanOWID.AsString()
}

func newComplaintTemplate(n string, b string) *template.Template {
	t, err := template.New(n).Parse(strings.TrimSpace(b))
	if err != nil {
		panic(err)
	}
	return t
}

func newComplaint(
	offerID *owid.OWID,
	swanID *owid.OWID) (*Complaint, error) {
	var err error

	// Set the static information associated with the complaint.
	var c Complaint
	c.DataProtectionURL = "https://ico.org.uk/make-a-complaint/your-personal-information-concerns/"
	c.Country = "Europe"

	// Work out the offer ID from the OWID provided.
	c.Offer, err = OfferFromOWID(offerID)
	if err != nil {
		return nil, err
	}

	// Set the OWIDs as strings.
	c.offerID = offerID
	c.swanOWID = swanID

	// Set the organisation as the domain for the moment.
	c.Organization = swanID.Domain

	// Return the complain data structure ready for the template email.
	return &c, nil
}

func handlerComplaintEmail(s *services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the form values from the input request.
		err := r.ParseForm()
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}

		// Check that the offer ID and the SWAN ID are present.
		if r.Form.Get("offerid") == "" {
			returnAPIError(
				&s.config,
				w,
				fmt.Errorf("'offerid' missing"),
				http.StatusBadRequest)
			return
		}
		if r.Form.Get("swanowid") == "" {
			returnAPIError(
				&s.config,
				w,
				fmt.Errorf("'swanowid' missing"),
				http.StatusBadRequest)
			return
		}

		// Get the SWAN OWIDs from the parameters.
		offerID, err := owid.FromBase64(r.Form.Get("offerid"))
		if err != nil {
			returnAPIError(
				&s.config,
				w,
				fmt.Errorf("'offerid' not a valid OWID"),
				http.StatusBadRequest)
			return
		}
		swanOWID, err := owid.FromBase64(r.Form.Get("swanowid"))
		if err != nil {
			returnAPIError(
				&s.config,
				w,
				fmt.Errorf("'swanowid' not a valid OWID"),
				http.StatusBadRequest)
			return
		}

		// Create the complaint object.
		c, err := newComplaint(offerID, swanOWID)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusBadRequest)
			return
		}

		// Get the strings for the subject and the body.
		var subject bytes.Buffer
		err = complaintSubjectTemplate.Execute(&subject, c)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}
		var body bytes.Buffer
		err = complaintBodyTemplate.Execute(&body, c)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}

		// Create the URL for the email.
		u := fmt.Sprintf("mailto:info@%s?subject=%s&body=%s",
			c.swanOWID.Domain,
			url.PathEscape(subject.String()),
			url.PathEscape(body.String()))

		// Return the URL as a text string.
		b := []byte(u)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
		_, err = w.Write(b)
		if err != nil {
			returnAPIError(&s.config, w, err, http.StatusInternalServerError)
			return
		}
	}
}
