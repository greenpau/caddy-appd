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
	"go.uber.org/zap"
)

// Service represents an application instance.
type Service struct {
	Unit   *Unit   `json:"unit,omitempty"`
	Status *Status `json:"status,omitempty"`
}

// Manager manages services.
type Manager struct {
	Services []*Service `json:"services,omitempty"`
	logger   *zap.Logger
}

// NewManager parses config and creates Manager instance.
func NewManager(cfg *Config, logger *zap.Logger) (*Manager, error) {
	m := &Manager{}
	m.logger = logger
	return m, nil
}

// Start starts services.
func (m *Manager) Start() []*Status {
	return nil
}

// Stop stops services.
func (m *Manager) Stop() []*Status {
	return nil
}
