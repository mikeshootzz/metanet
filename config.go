package main

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// DNSProvider holds the configuration for the Plesk DNS provider.
type DNSProvider struct {
	APIKey   string `json:"api_key,omitempty`  // X-API-Key header
	Username string `json:"username,omitempty` // basic auth user
	Password string `json:"password,omitempty` // basic auth pass
	BaseURL  string `json:"base_url,omitempty` // e.g. https://plesk.example.com/api/v2
}

// CaddyModule returns the module information.
func (DNSProvider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.plesk",
		New: func() caddy.Module { return new(DNSProvider) },
	}
}

// Provision sets defaults.
func (p *DNSProvider) Provision(ctx caddy.Context) error {
	if p.BaseURL == "" {
		p.BaseURL = "https://localhost/api/v2"
	}
	return nil
}

// UnmarshalCaddyfile parses the DNS provider configuration from Caddyfile.
func (p *DNSProvider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextBlock() {
			switch d.Val() {
			case "base_url":
				p.BaseURL = d.NextArg()
			case "api_key":
				p.APIKey = d.NextArg()
			case "username":
				p.Username = d.NextArg()
			case "password":
				p.Password = d.NextArg()
			default:
				return d.Errf("unknown option '%s'", d.Val())
			}
		}
	}
	return nil
}
