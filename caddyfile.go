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

package systemd

import (
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/greenpau/caddy-systemd/pkg/services"
)

func init() {
	httpcaddyfile.RegisterGlobalOption("systemd", parseCaddyfile)
}

// parseCaddyfile sets up a service manager.
//
// Syntax:
//
// systemd {
//   <command|app> <alias> {
//     workdir <path/to/dir>
//     cmd <path/to/command> [args]
//     args [arg1] [arg2] ... [argN]
//   }
//
//   command hostname {
//     cmd hostname
//   }
//
//   app test-py-http-server {
//     cmd python3 -m http.server 4080
//   }
//
//   app test-py-http-server-4081
//     cmd python3
//     args -m http.server 4081
//   }
// }

func parseCaddyfile(d *caddyfile.Dispenser, _ interface{}) (interface{}, error) {
	app := new(App)
	app.Config = services.NewConfig()

	if !d.Next() {
		return nil, d.ArgErr()
	}

	for d.NextBlock(0) {
		switch d.Val() {
		case "command", "app":
			args := d.RemainingArgs()
			if len(args) != 1 {
				return nil, d.ArgErr()
			}
			unit := services.NewUnit(args[0])
			for nesting := d.Nesting(); d.NextBlock(nesting); {
				k := d.Val()
				v := d.RemainingArgs()
				if len(v) < 1 {
					return nil, d.Errf("too few args for %q", k)
				}
				switch k {
				case "cmd":
					if len(v) > 1 {
						unit.Command = v[0]
						unit.Arguments = append(unit.Arguments, v[1:]...)
						break
					}
					unit.Command = v[0]
				case "args":
					unit.Arguments = append(unit.Arguments, v...)
				default:
					// return nil, d.Errf("k: %v, v: %v", k, v)
					return nil, d.Errf("unsupported %q key", k)
				}
			}
			app.Config.AddUnit(unit)
		default:
			return nil, d.ArgErr()
		}
	}

	return httpcaddyfile.App{
		Name:  appName,
		Value: caddyconfig.JSON(app, nil),
	}, nil
}
