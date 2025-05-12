package plesk

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// UnmarshalCaddyfile parses the provider block in a Caddyfile.
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
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
				return d.Errf("unknown directive '%s'", d.Val())
			}
		}
	}
	return nil
}

var _ caddyfile.Unmarshaler = (*Provider)(nil)
