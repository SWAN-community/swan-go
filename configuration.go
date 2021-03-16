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
	"os"
	"time"
)

// Configuration maps to the appsettings.json settings file.
type Configuration struct {
	Scheme  string        `json:"scheme"` // The scheme to use for requests
	Debug   bool          `json:"debug"`
	Timeout time.Duration `json:"valueTimeout"` // Seconds until the value provided expires and must be revalidated with SWAN
}

// NewConfig creates a new instance of configuration from the file provided.
func newConfig(file string) Configuration {
	var c Configuration
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&c)
	return c
}
