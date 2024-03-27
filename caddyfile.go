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
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/greenpau/caddy-appd/pkg/services"
)

func init() {
	httpcaddyfile.RegisterGlobalOption("appd", parseCaddyfile)
}

// parseCaddyfile sets up a service manager.
//
// Syntax:
//
// appd {
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

var argRules = map[string]argRule{
	"cmd":  argRule{Min: 1, Max: 255},
	"args": argRule{Min: 1, Max: 255},
	"noop": argRule{},
}

type argRule struct {
	Min int
	Max int
}

func parseCaddyfile(d *caddyfile.Dispenser, _ interface{}) (interface{}, error) {
	app := new(App)
	app.Config = services.NewConfig()

	if !d.Next() {
		return nil, d.ArgErr()
	}

	for d.NextBlock(0) {
		switch d.Val() {
		case "command", "app":
			dk := d.Val()
			args := d.RemainingArgs()
			if len(args) != 1 {
				return nil, d.ArgErr()
			}
			unit, err := services.NewUnit(dk, args[0])
			if err != nil {
				return nil, d.Errf("%s", err)
			}
			for nesting := d.Nesting(); d.NextBlock(nesting); {
				k := d.Val()
				v := d.RemainingArgs()
				if err := validateArg(k, v); err != nil {
					return nil, d.Errf("%s", err)
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
				case "noop":
					unit.Noop = true
				default:
					// return nil, d.Errf("k: %v, v: %v", k, v)
					return nil, d.Errf("unsupported %q key", k)
				}
			}
			if err := app.Config.AddUnit(unit); err != nil {
				return nil, d.Err(err.Error())
			}
		default:
			return nil, d.ArgErr()
		}
	}

	return httpcaddyfile.App{
		Name:  appName,
		Value: caddyconfig.JSON(app, nil),
	}, nil
}

func validateArg(k string, v []string) error {
	r, exists := argRules[k]
	if !exists {
		return nil
	}
	if r.Min > len(v) {
		return fmt.Errorf("too few args for %q directive", k)
	}
	if r.Max < len(v) {
		return fmt.Errorf("too many args for %q directive", k)
	}
	return nil
}
