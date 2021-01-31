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
    <link href="/bootstrap.min.css" rel="stylesheet">
    <style>
    .modal {
        display: block;
    }
    .reset {
        float:right;
        border: none;
        background: none;
        font-size: 0.8em;
        text-decoration: underline;
    }
    body, .blur {
        height: 100vh;
        width: 100vw;
    }
    body {
        background-repeat: no-repeat;
        background-position: center;
        background-color: {{ .BackgroundColor }};
        backdrop-filter: blur(4px);
    }
    .blur {
        background-color: rgba(0,0,0, 0.4);
    }
    @media only screen and (max-width: 600px) {
        body {
            background-image: url(//{{ .PublisherHost }}/background-600.png);
        }
    }
    @media only screen and (min-width: 600px) {
        body {
            background-image: url(//{{ .PublisherHost }}/background-767.png);
        }
    }
    @media only screen and (min-width: 768px) {
        {
            background-image: url(//{{ .PublisherHost }}/background-991.png);
        }
    }
    @media only screen and (min-width: 992px) {
        body {
            background-image: url(//{{ .PublisherHost }}/background-1199.png);
        }
    }
    @media only screen and (min-width: 1200px) {
        body {
            background-image: url(//{{ .PublisherHost }}/background.png);
        }
    }
    </style>
</head>
<body>
    <div class="blur"></div>
    <form method="POST">
        <div class="modal" style="display: block" tabindex="-1" role="dialog">
            <div class="modal-dialog modal-dialog-centered" role="document">
                <div class="modal-content">
                    <div class="modal-header">
                    <h5 class="modal-title">{{ .Title }}</h5>
                    <button type="submit" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true">×</span>
                    </button>
                </div>
                <div class="modal-body">
                    <div class="pt-3 pb-3">
                        <div class="form-group mb-6">
                            <label for="cbid">Common Browser ID (CBID)</label>
                            <input class="button-link reset" type="submit" value="Reset" name="reset-cbid"/>
                            <input type="text" class="form-control" id="cbid" name="cbid" value="{{ .CBID }}" readonly>
                            <small id="cbidHelp" class="form-text text-muted">
                                This ID enables your ad-funded access to the Open Web. You can reset it at any time. Recipients promise not to use for personalized marketing or associate this with your identity without consent.
                            </small>
                        </div>
                        <div class="form-group form-check mb-6 pl-2 py-4 bg-light">
                            <div>
                            <input type="checkbox" id="allow" name="allow" {{ if eq .Allow "on" }} checked {{ end }}>
                            <label class="form-check-label small" for="allow">Personalize Marketing</label>
                            <small id="allowHelp" class="form-text text-muted">
                                Tick to receive fewer more personalised adverts on this device.
                            </small>
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="email">Email address (optional)</label>
                            <input type="email" class="form-control" id="email" name="email" aria-describedby="emailHelp" placeholder="Optional email" value="{{ .Email }}">
                            <small id="emailHelp" class="form-text text-muted">
                                Advertisers, publishers and their partners will receive only an encrypted version of your email.
                                If you provide an email, and consent to personalised marketing, advertisers will try and show you fewer more personalized adverts across all your devices. 
                            </small>
                        </div>
                    </div>        
                </div>
                <div class="modal-footer">
                    <button type="submit" class="w-75 mx-auto btn btn-primary text-center">Update</button>
                </div>
            </div>
        </div>
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
