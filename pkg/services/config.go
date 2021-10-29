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
)

// Config is a configuration of Manager.
type Config struct {
	Units   []*Unit `json:"units,omitempty"`
	chain   []int
	unitMap map[string]*Unit
}

// NewConfig returns an instance of Config.
func NewConfig() *Config {
	return &Config{
		unitMap: make(map[string]*Unit),
	}
}

// AddUnit adds a unit entry to Config.
func (cfg *Config) AddUnit(u *Unit) error {
	if _, exists := cfg.unitMap[u.Name]; exists {
		return fmt.Errorf("unit %q already exists", u.Name)
	}
	cfg.Units = append(cfg.Units, u)
	cfg.unitMap[u.Name] = u
	return nil
}

func (cfg *Config) validate() error {
	for _, u := range cfg.Units {
		for _, dep := range u.Before {
			if _, exists := cfg.unitMap[dep]; !exists {
				return fmt.Errorf("the %q in %q directive for unit %q is not found", dep, "before", u.Name)
			}
		}
		for _, dep := range u.After {
			if _, exists := cfg.unitMap[dep]; !exists {
				return fmt.Errorf("the %q in %q directive for unit %q is not found", dep, "after", u.Name)
			}
		}
	}
	return nil
}

func (cfg *Config) order() error {
	return nil
}

// Services validates the config and creates a list of services.
func (cfg *Config) Services() ([]*Service, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	if err := cfg.order(); err != nil {
		return nil, err
	}

	svcs := []*Service{}
	for i, u := range cfg.Units {
		svc := &Service{
			Seq:  i,
			Unit: u,
		}
		svcs = append(svcs, svc)
	}
	return svcs, nil
}
