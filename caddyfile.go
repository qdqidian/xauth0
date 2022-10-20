package xauth0

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	httpcaddyfile.RegisterHandlerDirective("auth0", parseCaddyfile)
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	a := &HTTPAuth0Auth{}

	for h.Next() {
		args := h.RemainingArgs()

		switch len(args) {
		case 1:
			a.URL = args[0]
			a.HeaderPrefix = "auth0.user."
		case 2:
			a.URL = args[0]
			a.HeaderPrefix = args[1]
		default:
			return nil, h.ArgErr()
		}
	}

	return a, nil
}
