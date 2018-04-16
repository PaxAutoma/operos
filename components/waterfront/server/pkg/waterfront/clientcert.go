/*
Copyright 2018 Pax Automa Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package waterfront

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func MakeGenClientCertHandler(teamsterAddr string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getURL := url.URL{
			Scheme: "http",
			Host:   teamsterAddr,
			Path:   "/clientcert",
		}
		q := getURL.Query()
		q.Set("user", "admin")
		q.Add("group", "admin")
		q.Set("host", strings.SplitN(r.Host, ":", 2)[0])
		getURL.RawQuery = q.Encode()

		resp, err := http.Get(getURL.String())
		if err != nil {
			log.Printf("error while accessing Teamster: %v", err)
			return
		}
		defer resp.Body.Close()

		for name := range resp.Header {
			for _, value := range resp.Header[name] {
				w.Header().Set(name, value)
			}
		}
		w.WriteHeader(resp.StatusCode)

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("error while forwarding Teamster response to client: %v", err)
			return
		}
	})
}
