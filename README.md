# caddy-appd

<a href="https://github.com/greenpau/caddy-appd/actions/" target="_blank"><img src="https://github.com/greenpau/caddy-appd/workflows/build/badge.svg?branch=main"></a>
<a href="https://pkg.go.dev/github.com/greenpau/caddy-appd" target="_blank"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"></a>
<a href="https://caddy.community" target="_blank"><img src="https://img.shields.io/badge/community-forum-ff69b4.svg"></a>
<a href="https://caddyserver.com/docs/modules/appd" target="_blank"><img src="https://img.shields.io/badge/caddydocs-appd-green.svg"></a>

Service Management App for [Caddy v2](https://github.com/caddyserver/caddy).

Please ask questions either here or via LinkedIn. I am happy to help you! @greenpau

Please see other plugins:
* [caddy-security](https://github.com/greenpau/caddy-security)
* [caddy-trace](https://github.com/greenpau/caddy-trace)

<!-- begin-markdown-toc -->
## Table of Contents

* [Overview](#overview)
* [Getting Started](#getting-started)

<!-- end-markdown-toc -->

## Overview

The `caddy-appd` is a Caddy app that starts and stops non-Caddy
applications.

The primary use case is enabling starting the applications proxied by `caddy`
at startup. This way, there is no need to orchestrate the starting of services
prior to the starting of `caddy` itself.

## Getting Started

For example, the following configuration starts up `webapp1`
in `/usr/local/www/webapp` directory. The app is listening on port 8080.
When requests arrive for `webapp1.myfiosgateway.com`, they are being proxied
to `webapp1`.

```
{
  appd {
    app webapp1 {
      workdir /usr/local/www/webapp
      cmd /usr/local/bin/webapp
      args foo bar --foo=bar --port=8080
    }
  }
}

webapp1.myfiosgateway.com {
  route {
    reverse_proxy localhost:8080
  }
}
```

There is a sample config in `assets/config/Caddyfile`.

Test its configuration:

```bash
curl https://localhost:8443/version
curl https://localhost:8443/pytest/foo
```
