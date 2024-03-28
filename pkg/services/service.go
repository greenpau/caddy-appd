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

	"go.uber.org/zap"
)

// Service represents an application instance.
type Service struct {
	Seq    int        `json:"seq,omitempty"`
	Unit   *Unit      `json:"unit,omitempty"`
	Status *Status    `json:"status,omitempty"`
	State  *State     `json:"state,omitempty"`
	Kind   WorkerKind `json:"kind,omitempty"`
	logger *zap.Logger
	worker *worker
}

// NewService creates Service instance.
func NewService(seq int, unit *Unit, logger *zap.Logger) (*Service, error) {
	var k WorkerKind = WorkerKind(UnknownWorker)
	switch unit.Kind {
	case "command":
		k = WorkerKind(CommandWorker)
	case "app":
		k = WorkerKind(ApplicationWorker)
	default:
		return nil, fmt.Errorf("unsupported service kind: %s", unit.Kind)
	}

	svc := &Service{
		Seq:    seq + 1,
		Unit:   unit,
		Status: NewStatus(unit.Name, PendingStatus),
		State:  NewState(unit.Name, PendingState),
		Kind:   k,
		logger: logger,
	}
	return svc, nil
}

// Start starts Service instance.
func (svc *Service) Start() error {
	svc.logger.Debug("starting service",
		zap.String("service_name", svc.Unit.Name),
		zap.String("kind", svc.Unit.Kind),
		zap.Int("seq_id", svc.Seq),
	)

	switch svc.Kind {
	case WorkerKind(CommandWorker):
		err := newAdhocWorker(svc.Unit.Command, svc.Unit.Arguments, svc.Unit.StdOutFilePath, svc.Unit.StdErrFilePath)
		if err != nil {
			svc.State.Current = CompletedState
			svc.Status.Current = FailureStatus
			svc.Status.Error = err
			return err
		} else {
			svc.State.Current = CompletedState
			svc.State.Error = nil
			svc.Status.Current = SuccessStatus
			svc.Status.Error = nil
		}
	case WorkerKind(ApplicationWorker):
		w, err := newWorker(uint(svc.Unit.Seq), svc.Unit.Command, svc.Unit.Arguments, svc.Unit.StdOutFilePath, svc.Unit.StdErrFilePath, svc.logger)
		if err != nil {
			svc.State.Current = CompletedState
			svc.Status.Current = FailureStatus
			svc.Status.Error = err
			return err
		} else {
			svc.worker = w
			svc.State.Current = CompletedState
			svc.State.Error = nil
			svc.Status.Current = SuccessStatus
			svc.Status.Error = nil
		}
	default:
		svc.State.Current = CompletedState
		svc.Status.Current = FailureStatus
		svc.Status.Error = fmt.Errorf("unsupported worker type: %s", svc.Kind)
		return svc.Status.Error
	}

	return nil
}

// Stop stops Service instance.
func (svc *Service) Stop() error {
	switch svc.Kind {
	case WorkerKind(CommandWorker):
		svc.logger.Debug("skipped stopping service",
			zap.String("service_name", svc.Unit.Name),
			zap.String("kind", svc.Unit.Kind),
			zap.String("reason", "command"),
			zap.Int("seq_id", svc.Seq),
		)
		return nil
	case WorkerKind(ApplicationWorker):
		svc.logger.Debug("stopping service",
			zap.String("service_name", svc.Unit.Name),
			zap.String("kind", svc.Unit.Kind),
			zap.Int("seq_id", svc.Seq),
		)
		workerState, workerStatus := svc.worker.stop()
		svc.State.Current = workerState.Current
		svc.State.Error = workerState.Error
		svc.Status.Current = workerStatus.Current
		svc.Status.Error = workerStatus.Error
		if svc.Status.Error != nil {
			svc.logger.Debug("failed stopping service",
				zap.String("service_name", svc.Unit.Name),
				zap.String("kind", svc.Unit.Kind),
				zap.Int("seq_id", svc.Seq),
			)
			return svc.Status.Error
		}
		svc.logger.Debug("stopped service",
			zap.String("service_name", svc.Unit.Name),
			zap.String("kind", svc.Unit.Kind),
			zap.Int("seq_id", svc.Seq),
		)
	default:
		svc.logger.Debug("stopping service",
			zap.String("service_name", svc.Unit.Name),
			zap.String("kind", svc.Unit.Kind),
			zap.Int("seq_id", svc.Seq),
		)
	}

	svc.State.Current = CompletedState
	svc.Status.Current = SuccessStatus
	return nil
}
