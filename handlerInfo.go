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
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"owid"
)

// infoModel data needed for the advert information interface.
type infoModel struct {
	OWIDs     map[*owid.OWID]interface{}
	Bid       *Bid
	Offer     *Offer
	Root      *owid.OWID
	ReturnURL template.HTML
}

func (m *infoModel) findOffer() (*owid.OWID, *Offer) {
	for k, v := range m.OWIDs {
		if o, ok := v.(*Offer); ok {
			return k, o
		}
	}
	return nil, nil
}

func (m *infoModel) findBid() *Bid {
	for _, v := range m.OWIDs {
		if b, ok := v.(*Bid); ok {
			return b
		}
	}
	return nil
}

func infoRole(s interface{}) string {
	_, fok := s.(*Failed)
	_, bok := s.(*Bid)
	_, eok := s.(*Empty)
	_, ook := s.(*Offer)
	if fok {
		return "Failed"
	}
	if bok {
		return "Bid"
	}
	if eok {
		return "Empty"
	}
	if ook {
		return "Offer"
	}
	return ""
}

func handlerInfo(s *services, h string) (http.HandlerFunc, error) {
	t := template.Must(template.New("info").Funcs(template.FuncMap{
		"role": infoRole,
	}).Parse(removeHTMLWhiteSpace(h)))
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the SWAN OWIDs from the form parameters.
		err := r.ParseForm()
		if err != nil {
			returnServerError(&s.config, w, err)
		}
		var m infoModel
		m.OWIDs = make(map[*owid.OWID]interface{})
		for _, vs := range r.Form {
			for _, v := range vs {
				o, err := owid.FromBase64(v)
				if err != nil {
					returnRequestError(&s.config, w, err, http.StatusBadRequest)
				}
				m.OWIDs[o], err = FromOWID(o)
				if err != nil {
					returnRequestError(&s.config, w, err, http.StatusBadRequest)
				}
			}
		}

		// Set the common fields.
		m.Bid = m.findBid()
		m.Root, m.Offer = m.findOffer()
		f, err := getReferer(r)
		if err != nil {
			returnRequestError(&s.config, w, err, http.StatusBadRequest)
		}
		m.ReturnURL = template.HTML(f)

		// Display the template form.
		err = t.Execute(w, m)
		if err != nil {
			fmt.Println(err.Error())
			returnServerError(&s.config, w, err)
			return
		}
	}, nil
}

func getReferer(r *http.Request) (string, error) {
	u, err := url.Parse(r.Header.Get("Referer"))
	if err != nil {
		return "", err
	}
	u.RawQuery = ""
	return u.String(), nil
}
