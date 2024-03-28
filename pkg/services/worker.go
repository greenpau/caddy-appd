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
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

var workerStopTimeout = time.Second * 4

type WorkerKind int

const (
	UnknownWorker StateKind = iota
	CommandWorker
	ApplicationWorker
)

func (k WorkerKind) String() string {
	return [...]string{"unknown", "command", "application"}[k]
}

func (k WorkerKind) EnumIndex() int {
	return int(k)
}

func (k WorkerKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToLower(k.String()))
}

type worker struct {
	mu     sync.RWMutex
	ID     uint
	Cmd    *exec.Cmd
	Pid    int
	logger *zap.Logger
}

func newWorker(id uint, binPath string, args []string, stdOutFilePath, stdErrFilePath string, logger *zap.Logger) (*worker, error) {
	w := &worker{
		ID:     id,
		logger: logger,
	}

	cmd := exec.Command(binPath, args...)

	var outFile, errFile *os.File
	var err error

	if stdOutFilePath == "" {
		cmd.Stdout = os.Stdout
	} else {
		outFile, err = os.OpenFile(stdOutFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed opening output file: %w", err)
		}
		defer outFile.Close()
		cmd.Stdout = outFile
	}

	if stdErrFilePath == "" {
		if stdOutFilePath != "" {
			cmd.Stderr = outFile
		} else {
			cmd.Stderr = os.Stderr
		}
	} else {
		errFile, err = os.OpenFile(stdErrFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed opening error file: %w", err)
		}
		defer errFile.Close()
		cmd.Stderr = errFile
	}

	w.Cmd = cmd
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	w.Pid = cmd.Process.Pid
	return w, nil
}

func (w *worker) stop() (*State, *Status) {
	w.mu.Lock()
	defer w.mu.Unlock()

	state := &State{
		Current: UnknownState,
	}
	status := &Status{
		Current: UnknownStatus,
	}
	if w.Pid < 1 {
		state.Current = CompletedState
		status.Current = FailureStatus
		status.Error = fmt.Errorf("pid is 0")
		return state, status
	}

	if w.Cmd == nil {
		state.Current = CompletedState
		status.Current = FailureStatus
		status.Error = fmt.Errorf("cmd exec is nil")
		return state, status
	}

	if w.Cmd.Process == nil {
		state.Current = CompletedState
		status.Current = FailureStatus
		status.Error = fmt.Errorf("process is nil")
		return state, status
	}

	cmdResult := make(chan error, 1)
	go func() {
		_, err := w.Cmd.Process.Wait()
		cmdResult <- err
	}()

	if err := w.Cmd.Process.Signal(os.Interrupt); err != nil {
		state.Current = CompletedState
		status.Current = FailureStatus
		status.Error = err
		return state, status
	}

	select {
	case err := <-cmdResult:
		state.Current = CompletedState
		status.Current = FailureStatus
		status.Error = err
		break
	case <-time.After(workerStopTimeout):
		if err := w.Cmd.Process.Kill(); err != nil {
			state.Current = CompletedState
			status.Current = FailureStatus
			status.Error = fmt.Errorf("force terminated failed: %w", err)
			return state, status
		}
		state.Current = CompletedState
		status.Current = FailureStatus
		status.Error = fmt.Errorf("force terminated process")
		break
	}

	return state, status
}

func newAdhocWorker(binPath string, args []string, stdOutFilePath, stdErrFilePath string) error {
	cmd := exec.Command(binPath, args...)
	var outFile, errFile *os.File
	var err error

	if stdOutFilePath == "" {
		cmd.Stdout = os.Stdout
	} else {
		outFile, err = os.OpenFile(stdOutFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("failed opening output file: %w", err)
		}
		defer outFile.Close()
		cmd.Stdout = outFile
	}

	if stdErrFilePath == "" {
		if stdOutFilePath != "" {
			cmd.Stderr = outFile
		} else {
			cmd.Stderr = os.Stderr
		}
	} else {
		errFile, err = os.OpenFile(stdErrFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("failed opening error file: %w", err)
		}
		defer errFile.Close()
		cmd.Stderr = errFile
	}

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
