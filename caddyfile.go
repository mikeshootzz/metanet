package plesk

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// UnmarshalCaddyfile parses the provider block in a Caddyfile.
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextBlock(0) {
			switch d.Val() {
			case "base_url":
				if !d.NextArg() {
					return d.Err("missing base_url argument")
				}
				p.BaseURL = d.Val()
			case "api_key":
				if !d.NextArg() {
					return d.Err("missing api_key argument")
				}
				p.APIKey = d.Val()
			case "username":
				if !d.NextArg() {
					return d.Err("missing username argument")
				}
				p.Username = d.Val()
			case "password":
				if !d.NextArg() {
					return d.Err("missing password argument")
				}
				p.Password = d.Val()
			default:
				return d.Errf("unknown directive '%s'", d.Val())
			}
		}
	}
	return nil
}

var _ caddyfile.Unmarshaler = (*Provider)(nil)
