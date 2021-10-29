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

package services

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestNewUnit(t *testing.T) {
	testcases := []struct {
		name      string
		unit      *Unit
		want      *Unit
		shouldErr bool
		err       error
	}{
		{
			name: "test new unit",
			unit: &Unit{Name: "hostname", Kind: "command"},
			want: &Unit{Name: "hostname", Kind: "command"},
		},
		{
			name:      "test empty unit alias",
			unit:      &Unit{},
			shouldErr: true,
			err:       fmt.Errorf("empty unit alias"),
		},
		{
			name:      "test invalid unit alias",
			unit:      &Unit{Name: "fo"},
			shouldErr: true,
			err:       fmt.Errorf("invalid unit alias: %q", "fo"),
		},
		{
			name:      "test empty unit type",
			unit:      &Unit{Name: "hostname"},
			shouldErr: true,
			err:       fmt.Errorf("unit %q: empty type", "hostname"),
		},
		{
			name:      "test invalid unit type",
			unit:      &Unit{Name: "hostname", Kind: "foo"},
			shouldErr: true,
			err:       fmt.Errorf("unit %q: invalid %q type", "hostname", "foo"),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewUnit(tc.unit.Kind, tc.unit.Name)
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
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("unit mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
