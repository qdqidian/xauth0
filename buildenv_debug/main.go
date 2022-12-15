package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// plug in Caddy modules here
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "qidian.xyz/auth0/xauth0"
)

func main() {
	caddycmd.Main()
}
