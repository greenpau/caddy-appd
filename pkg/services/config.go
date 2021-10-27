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

// Unit is a configuration for a command or app.
type Unit struct {
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	Wants         []string `json:"wants,omitempty"`
	Requires      []string `json:"requires,omitempty"`
	Before        []string `json:"before,omitempty"`
	After         []string `json:"after,omitempty"`
	WorkDirectory string   `json:"workdir,omitempty"`
	Command       string   `json:"cmd,omitempty"`
	Arguments     []string `json:"args,omitempty"`
	Priority      uint64   `json:"priority,omitempty"`
}

// Config is a configuration of Manager.
type Config struct {
	Units []*Unit `json:"units,omitempty"`
}

// NewConfig returns an instance of Config.
func NewConfig() *Config {
	return &Config{}
}

// NewUnit returns an instance of Unit.
func NewUnit(s string) *Unit {
	return &Unit{
		Name: s,
	}
}

// AddUnit adds a unit entry to Config.
func (cfg *Config) AddUnit(u *Unit) {
	cfg.Units = append(cfg.Units, u)
}
