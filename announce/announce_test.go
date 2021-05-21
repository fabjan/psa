//	Copyright 2021 Fabian Bergström
//	
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//	
//			http://www.apache.org/licenses/LICENSE-2.0
//	
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.

package announce

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type mockReq struct {
	method string
	url    string
	body   string
	form   map[string]string
	header map[string]string
}

func TestDiscorcian_makeRequest(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		msg     string
		wantErr bool
		want    mockReq
	}{
		{
			name: "a normal message",
			url:  "http://hook.example.com",
			msg:  "test message, please ignore",
			want: mockReq{
				method: "POST",
				url:    "http://hook.example.com",
				header: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
				form:   map[string]string{"content": "test message, please ignore"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			if err != nil {
				t.Errorf("bad test case, inalid URL: <%s>", tt.url)
			}
			s := DiscordHook(u)
			got, err := s.makeRequest(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Discorcian.makeRequest() [%s] error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			errs := checkReq(got, tt.want)
			for _, err := range errs {
				t.Errorf("Discorcian.makeRequest() [%s] request check failed: %v", tt.name, err)
			}
		})
	}
}
func TestSlacker_makeRequest(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		msg     string
		wantErr bool
		want    mockReq
	}{
		{
			name: "a normal message",
			url:  "http://hook.example.com",
			msg:  "test message, please ignore",
			want: mockReq{
				method: "POST",
				url:    "http://hook.example.com",
				header: map[string]string{"Content-Type": "application/json"},
				body:   `{"text":"test message, please ignore"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			if err != nil {
				t.Errorf("bad test case, inalid URL: <%s>", tt.url)
			}
			s := SlackHook(u)
			got, err := s.makeRequest(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Slacker.makeRequest() [%s] error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			errs := checkReq(got, tt.want)
			for _, err := range errs {
				t.Errorf("Slacker.makeRequest() [%s] request check failed: %v", tt.name, err)
			}
		})
	}
}

func checkReq(got *http.Request, want mockReq) []error {
	errs := []error{}
	if want.method != "" && got.Method != want.method {
		err := fmt.Errorf("method: got %q, want %q", got.Method, want.method)
		errs = append(errs, err)
	}
	if want.url != "" && got.URL.String() != want.url {
		err := fmt.Errorf("URL: got %q, want %q", got.URL.String(), want.url)
		errs = append(errs, err)
	}
	for key, wantV := range want.header {
		gotV := got.Header.Get(key)
		if gotV != wantV {
			err := fmt.Errorf("headers[%s]: got %q, want %q", key, gotV, wantV)
			errs = append(errs, err)
		}
	}
	if want.body != "" {
		bytes, err := io.ReadAll(got.Body)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not read body: %w", err))
		}
		gotBody := strings.TrimSpace(string(bytes))
		if gotBody != want.body {
			err = fmt.Errorf("body: got %q, want %q", gotBody, want.body)
			errs = append(errs, err)
		}
	}
	for key, wantV := range want.form {
		gotV := got.FormValue(key)
		if gotV != wantV {
			err := fmt.Errorf("form[%s]: got %q, want %q", key, gotV, wantV)
			errs = append(errs, err)
		}
	}
	return errs
}
