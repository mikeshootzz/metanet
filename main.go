package main

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddy-dns"
)

func init() {
	caddy.RegisterModule(DNSProvider{})
}

// Interface guards
var _ caddy_dns.DNSProvider = (*DNSProvider)(nil)
