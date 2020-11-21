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
	"html/template"
	"strings"
)

var captureTemplate = newHTMLTemplate("capture", `
<!DOCTYPE html>
<html>
<head>
    <link rel="icon" href="data:;base64,=">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <title>{{ .Title }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body {
            margin: 0;
            padding: 0;
            font-family: nunito, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            background-color: {{ .BackgroundColor }};
        }
        input {
            border-radius: 0;
        }
        form {
            background-color: white;
            border-color: black;
            border-width: 2px;
            border-style: solid;
            color: black;
            padding: 1em;
        }
        p {
            display: inline;
        }
        ul {
            list-style: none;
            padding: 0;
            margin: 0;
        }
        li {
            vertical-align: middle;
        }
        li span {
            margin-left: 0.5em;
            font-size: 0.8em;
            color: slategray;
        }
        h1 {
            text-align: center;
            font-size: 1.5em;
        }
        label {
            font-size: 1.2em;
        }
        @media only screen and (max-width: 480px) {
        label {
            font-size: 0.9em;
        }
        }
        .textbox {
            margin-bottom: 0.5em;
            line-height: 1em;
            padding: 0.2em;
            box-sizing: border-box; 
        }
        .textbox-size {
            width: 480px;
            font-size: 1.5em;
        }
        @media only screen and (max-width: 480px) {
        .textbox-size {
            width: 100%;
            font-size: 1.2em;
        }
        }
        .disabled {
            background-color: lightgray;
            border: none;
        }
        .checkbox input {
            margin: 0.5em auto;
            border-radius: 5px;
            height: 1.2em;
            width: 1.2em;
        }
        .checkbox label {
            margin-left: 0.5em;
            display: inline;
        }
        @media only screen and (max-width: 480px) {
            .checkbox span {
                display: block;
            }
        }
        .button {
            float: right;
            margin-top: 1em;
        }
        .button-link {
            background-color: transparent;
            text-decoration: underline;
            border: none;
            cursor: pointer;
        }        
        .button input, .button a {
            -webkit-appearance: none;
            display: block;
            padding: 0.5em;
            background-color: black;
            text-decoration: none;
            color: white;
            border: none;
            border-radius: 0;
        }
        .button-size input, .button-size a {
            font-size: 1.2em;
        }
        @media only screen and (max-width: 480px) {
        .button-size input, .button-size a {
            font-size: 1em;
        }
        }
        .error {
            margin-bottom: 1em;
            color: darkred;
        }
    </style>
</head>
<body>
    <form method="POST">
        <ul>
            <li>
                <h1>{{ .Title }}</h1>
            </li>
            <li>
                <label for="cbid">Common Browser Id</label>
                <input style="display: none;" type="submit" value="Update"/>
                <input style="float:right" class="button-link" type="submit" value="Reset" name="reset-cbid"/>    
            </li>
            <li>
                <input class="textbox disabled textbox-size" type="text" id="cbid" name="cbid" value="{{ .CBID }}" readonly/>
            </li>
            <li>
                <label for="email">Email</label>
            </li>
            <li>
                <input class="textbox textbox-size" type="text" id="email" name="email" value="{{ .Email }}"/>
            </li>
            <li>
                <div class="checkbox">
                    <input type="checkbox" id="allow" name="allow" {{ if eq .Allow "on" }} checked {{ end }}/>
                    <label for="allow">Personalize marketing</label>
                </div>
            </li>
            <li>
                <div class="button button-size">
                    <input type="submit" value="Update"/>
                </div>
                <div class="button button-size" style="float: left;">
                    <input style="background-color: grey;" type="submit" value="Reset" name="reset-all"/>
                </div>            
            </li>
        </ul>
    </form>
</body>
</html>`)

func newHTMLTemplate(n string, h string) *template.Template {
	c := removeHTMLWhiteSpace(h)
	return template.Must(template.New(n).Parse(c))
}

// Removes white space from the HTML string provided whilst retaining valid
// HTML.
func removeHTMLWhiteSpace(h string) string {
	var sb strings.Builder
	for i, r := range h {

		// Only write out runes that are not control characters.
		if r != '\r' && r != '\n' && r != '\t' {

			// Only write this rune if the rune is not a space, or if it is a
			// space the preceding rune is not a space.
			if i == 0 || r != ' ' || h[i-1] != ' ' {
				sb.WriteRune(r)
			}
		}
	}
	return sb.String()
}
