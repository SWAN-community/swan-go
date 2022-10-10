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
	"bytes"
	"encoding/json"

	"github.com/SWAN-community/common-go"
)

// PreferencesData
// https://github.com/OneKey-Network/addressability-framework/blob/main/mvp-spec/model/preferences-data.md
type PreferencesData struct {
	UseBrowsingForPersonalization bool `json:"use_browsing_for_personalization"`
}

func (p *PreferencesData) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["use_browsing_for_personalization"] = p.UseBrowsingForPersonalization
	return json.Marshal(m)
}

func (p *PreferencesData) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	err := common.WriteBool(&b, p.UseBrowsingForPersonalization)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (p *PreferencesData) UnmarshalBinary(data []byte) error {
	var err error
	b := bytes.NewBuffer(data)
	p.UseBrowsingForPersonalization, err = common.ReadBool(b)
	if err != nil {
		return err
	}
	return nil
}
