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
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Unit is a configuration for a command or app.
type Unit struct {
	// The order of the unit in Config.
	Seq int `json:"seq,omitempty"`
	// The alias for the Unit.
	Name string `json:"name,omitempty"`
	// The description of the Unit.
	Description string `json:"description,omitempty"`
	// The type of the Unit: app, command.
	Kind string `json:"kind,omitempty"`
	// The command to execute.
	Command string `json:"cmd,omitempty"`
	// The executed command arguments.
	Arguments []string `json:"args,omitempty"`
	// The directory the command starts in.
	WorkDirectory string `json:"workdir,omitempty"`
	// The higher the Priority the sooner this unit activates. The Manager
	// activates units with the same Priority in alphabetical order.
	Priority uint64 `json:"priority,omitempty"`
	// If set to true, the unit will not be started or stopped.
	Noop bool `json:"noop,omitempty"`
	// Not implemented.
	Wants []string `json:"wants,omitempty"`
	// Not implemented.
	Requires []string `json:"requires,omitempty"`
	// The names of the Units that need to start after this one.
	Before []string `json:"before,omitempty"`
	// The names of the Units that need to start prior to this one.
	After []string `json:"after,omitempty"`
	// The path to standard output file.
	StdOutFilePath string `json:"std_out_file_path,omitempty"`
	// The path to err output file.
	StdErrFilePath string `json:"std_err_file_path,omitempty"`
}

// NewUnit returns an instance of Unit.
func NewUnit(kind, name string) (*Unit, error) {
	name = strings.TrimSpace(name)
	kind = strings.TrimSpace(kind)
	if name == "" {
		return nil, fmt.Errorf("empty unit alias")
	}
	re := regexp.MustCompile("^[a-zA-Z0-9_-]{3,100}$")
	if !re.MatchString(name) {
		return nil, fmt.Errorf("invalid unit alias: %q", name)
	}
	switch kind {
	case "app", "command":
	case "":
		return nil, fmt.Errorf("unit %q: empty type", name)
	default:
		return nil, fmt.Errorf("unit %q: invalid %q type", name, kind)
	}
	return &Unit{Name: name, Kind: kind}, nil
}

func validateFilePath(fp string) error {
	fsfi, err := os.Stat(fp)
	if err == nil {
		if fsfi.IsDir() {
			return fmt.Errorf("file path is directory")
		}
		return nil
	}
	if os.IsNotExist(err) {
		// Check if parent exists
		fd := filepath.Dir(fp)
		_, parentErr := os.Stat(fd)
		if parentErr == nil {
			return nil
		}
		if os.IsNotExist(parentErr) {
			return fmt.Errorf("parent directory does not exist: %s", fd)
		}
		return fmt.Errorf("parent directory erred: %s", fd)
	}
	return err
}
