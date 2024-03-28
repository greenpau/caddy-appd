// Copyright 2024 Paul Greenberg greenpau@outlook.com
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
	"encoding/json"
	"strings"
)

type StateKind int

const (
	UnknownState StateKind = iota
	PendingState
	RunningState
	StoppedState
	PausedState
	CompletedState
)

// State represent the last recorded state of a service.
type State struct {
	Current     StateKind `json:"current,omitempty"`
	ServiceName string    `json:"service_name,omitempty"`
	Error       error     `json:"error,omitempty"`
}

// NewState creates State instance.
func NewState(name string, kind StateKind) *State {
	st := &State{
		ServiceName: name,
		Current:     kind,
	}
	return st
}

func (k StateKind) String() string {
	return [...]string{"Unknown", "Pending", "Running", "Stopped", "Paused", "Completed"}[k]
}

func (k StateKind) EnumIndex() int {
	return int(k)
}

func (k StateKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToLower(k.String()))
}
