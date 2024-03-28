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

package appd

import (
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/greenpau/caddy-appd/pkg/services"
	"go.uber.org/zap"
)

var (
	appName = "appd"

	// Interface guards
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.Module      = (*App)(nil)
	_ caddy.App         = (*App)(nil)
)

func init() {
	caddy.RegisterModule(App{})
}

// App implements systemd-like service manager.
type App struct {
	Name    string           `json:"-"`
	Config  *services.Config `json:"config,omitempty"`
	manager *services.Manager
	logger  *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  caddy.ModuleID(appName),
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the service manager.
func (app *App) Provision(ctx caddy.Context) error {
	app.Name = appName
	app.logger = ctx.Logger(app)

	app.logger.Info(
		"provisioning app instance",
		zap.String("app", app.Name),
	)

	manager, err := services.NewManager(app.Config, app.logger)
	if err != nil {
		app.logger.Error(
			"failed configuring app instance",
			zap.String("app", app.Name),
			zap.Error(err),
		)
		return err
	}
	app.manager = manager

	app.logger.Info(
		"provisioned app instance",
		zap.String("app", app.Name),
	)
	return nil
}

// Start starts the service manager and associated services.
func (app App) Start() error {
	app.logger.Debug(
		"starting service manager",
		zap.String("app", app.Name),
	)

	if msgs := app.manager.Start(); msgs != nil {
		for _, msg := range msgs {
			app.logger.Error(
				"failed to start service",
				zap.String("app", app.Name),
				zap.String("service_name", msg.ServiceName),
				zap.Error(msg.Error),
			)
		}
		if len(msgs) > 0 {
			return fmt.Errorf("service manager failed to start services")
		}
	}

	app.logger.Debug(
		"started service manager",
		zap.String("app", app.Name),
	)

	return nil
}

// Stop stops the service manager and associated services.
func (app App) Stop() error {
	app.logger.Debug(
		"stopping service manager",
		zap.String("app", app.Name),
	)

	if msgs := app.manager.Stop(); msgs != nil {
		for _, msg := range msgs {
			app.logger.Error(
				"failed to stop service",
				zap.String("app", app.Name),
				zap.String("service_name", msg.ServiceName),
				zap.Error(msg.Error),
			)
		}
		if len(msgs) > 0 {
			return fmt.Errorf("service manager failed to stop services")
		}
	}

	app.logger.Debug(
		"stopped service manager",
		zap.String("app", app.Name),
	)
	return nil
}
