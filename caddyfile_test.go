// Copyright 2021 Paul Greenberg greenpau@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package systemd

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/google/go-cmp/cmp"
)

const tf string = "Testfile"

func TestParseCaddyfile(t *testing.T) {
	testcases := []struct {
		name      string
		d         *caddyfile.Dispenser
		want      string
		shouldErr bool
		err       error
	}{
		{
			name: "test parse config",
			d: caddyfile.NewTestDispenser(`
            systemd {
              command hostname {
                cmd hostname
              }
              app test-py-http-server {
                cmd python3 -m http.server 4080
              }
              app test-py-http-server-4081 {
                cmd python3
                args -m http.server 4081
              }
            }`),
			want: `{
			  "config": {
                "units": [
                  {
                    "name":"hostname",
					"cmd":"hostname",
					"kind":"command",
					"seq": 1
                  },
                  {
                    "name":"test-py-http-server",
					"cmd":"python3",
					"args":["-m","http.server","4080"],
					"kind":"app",
					"seq": 2
                  },
                  {
                    "name":"test-py-http-server-4081",
                    "cmd":"python3",
                    "args":["-m","http.server","4081"],
					"kind":"app",
					"seq": 3
                  }
                ]
              }
			}`,
		},
		{
			name: "test parse config with unsupported unit key",
			d: caddyfile.NewTestDispenser(`
            systemd {
              command foo {
                bar baz
              }
            }`),
			shouldErr: true,
			err:       fmt.Errorf("%s:%d - Error during parsing: unsupported %q key, import chain: ['']", tf, 4, "bar"),
		},
		{
			name: "test parse config with too few arg for unit arg",
			d: caddyfile.NewTestDispenser(`
            systemd {
              command foo {
                cmd
              }
            }`),
			shouldErr: true,
			err:       fmt.Errorf("%s:%d - Error during parsing: too few args for %q directive, import chain: ['']", tf, 4, "cmd"),
		},
		{
			name: "test parse config with too many arg for unit arg",
			d: caddyfile.NewTestDispenser(`
            systemd {
              command foo {
                noop foo
              }
            }`),
			shouldErr: true,
			err:       fmt.Errorf("%s:%d - Error during parsing: too many args for %q directive, import chain: ['']", tf, 4, "noop"),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			app, err := parseCaddyfile(tc.d, nil)
			if err != nil {
				if !tc.shouldErr {
					t.Fatalf("expected success, got: %v", err)
				}
				if diff := cmp.Diff(err.Error(), tc.err.Error()); diff != "" {
					t.Fatalf("unexpected error: %v, want: %v", err, tc.err)
				}
				return
			}
			if tc.shouldErr {
				t.Fatalf("unexpected success, want: %v", tc.err)
			}
			got := unpack(t, string(app.(httpcaddyfile.App).Value))
			want := unpack(t, tc.want)

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("parseCaddyfile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func unpack(t *testing.T, s string) (m map[string]interface{}) {
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("failed to parse %q: %v", s, err)
	}
	return m
}
