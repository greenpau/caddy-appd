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
	"sort"
	"sync"

	"go.uber.org/zap"
)

// Manager manages services.
type Manager struct {
	mu          sync.Mutex
	Services    []*Service `json:"services,omitempty"`
	provisioned bool
	started     bool
	logger      *zap.Logger
}

// NewManager parses config and creates Manager instance.
func NewManager(cfg *Config, logger *zap.Logger) (*Manager, error) {
	m := &Manager{}
	m.logger = logger
	cfg.unitOrderAsc()
	logger.Debug("initializing manager", zap.Any("configuration", cfg))
	for i, unit := range cfg.Units {
		svc, err := NewService(i, unit, logger)
		if err != nil {
			return nil, err
		}
		m.Services = append(m.Services, svc)
	}
	logger.Debug("configured services", zap.Any("services", m.Services))
	m.provisioned = true
	return m, nil
}

// Start starts services.
func (m *Manager) Start() []*Status {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.provisioned {
		return []*Status{
			{
				Current:     FailureStatus,
				ServiceName: "all",
				Error:       fmt.Errorf("provisioning has failed"),
			}}
	}

	for _, svc := range m.Services {
		if svc.Unit.Noop {
			m.logger.Debug("skipped starting service",
				zap.String("service_name", svc.Unit.Name),
				zap.String("kind", svc.Unit.Kind),
				zap.String("reason", "noop"),
				zap.Int("seq_id", svc.Seq),
			)
			continue
		}

		if svc.Unit.StdOutFilePath != "" {
			if err := validateFilePath(svc.Unit.StdOutFilePath); err != nil {
				return []*Status{
					{
						Current:     FailureStatus,
						ServiceName: svc.Unit.Name,
						Error:       err,
					}}
			}
		}

		if svc.Unit.StdErrFilePath != "" {
			if err := validateFilePath(svc.Unit.StdErrFilePath); err != nil {
				return []*Status{
					{
						Current:     FailureStatus,
						ServiceName: svc.Unit.Name,
						Error:       err,
					}}
			}
		}

		if err := svc.Start(); err != nil {
			return []*Status{svc.Status}
		}
	}

	m.started = true
	return nil
}

// Stop stops services.
func (m *Manager) Stop() []*Status {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.provisioned {
		return []*Status{
			{
				Current:     FailureStatus,
				ServiceName: "all",
				Error:       fmt.Errorf("provisioning has failed"),
			}}
	}

	desc := func(a, b int) bool {
		return m.Services[a].Seq > m.Services[b].Seq
	}

	sort.Slice(m.Services, desc)

	svcErrors := []*Status{}

	for _, svc := range m.Services {
		if svc.Unit.Noop {
			m.logger.Debug("skipped stopping service",
				zap.String("service_name", svc.Unit.Name),
				zap.String("kind", svc.Unit.Kind),
				zap.String("reason", "noop"),
				zap.Int("seq_id", svc.Seq),
			)
			continue
		}
		if err := svc.Stop(); err != nil {
			svcErrors = append(svcErrors, svc.Status)
		}
	}

	m.started = false
	return svcErrors
}
