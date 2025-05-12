# Plesk Module for Caddy

This package implements a DNS provider module for [Caddy v2](https://github.com/caddyserver/caddy), allowing you to automate DNS-01 challenges (and dynamic record updates) against a Plesk server’s API.

## Module Name

```text
dns.providers.plesk
```

## Installation

You must build a custom Caddy binary that bundles this module. The easiest way is with [xcaddy](https://github.com/caddyserver/xcaddy):

```bash
xcaddy build \
  --with github.com/mikeshootzz/metanet
```

That produces a `caddy` executable in your working directory with `dns.providers.plesk` registered.

## Configuration

### Provider Options

| Option     | Required? | Default                    | Description                         |
|------------|-----------|----------------------------|-------------------------------------|
| `base_url` | no        | `https://localhost/api/v2` | Full URL to your Plesk API endpoint |
| `api_key`  | no        | ―                          | X-API-Key header value              |
| `username` | no        | ―                          | HTTP Basic Auth username            |
| `password` | no        | ―                          | HTTP Basic Auth password            |

You must supply **either** `api_key` **or** both `username` + `password`, depending on how your Plesk instance is secured.

### Environment Variables

You can also reference environment variables in your Caddyfile:

```caddyfile
tls {
  dns plesk {
    base_url   {env.PLESK_BASE_URL}
    api_key    {env.PLESK_API_KEY}
    username   {env.PLESK_USER}
    password   {env.PLESK_PASS}
  }
}
```

```bash
export PLESK_BASE_URL="https://plesk.example.com/api/v2"
export PLESK_API_KEY="abc123"
export PLESK_USER="admin"
export PLESK_PASS="secret"
```

## Caddyfile Examples

### Globally (default for all sites)

```caddyfile
{
  # Use Plesk as the default DNS-01 issuer for all sites
  acme_dns plesk {
    base_url   https://plesk.example.com/api/v2
    api_key    YOUR_API_KEY
  }
}

# all sites from here on will use DNS-01 via Plesk
example.com, www.example.com {
  reverse_proxy localhost:8080
}
```

### Per-site

```caddyfile
example.com {
  reverse_proxy localhost:8080

  tls {
    dns plesk {
      base_url https://plesk.example.com/api/v2
      username admin
      password yourPassword
    }
  }
}
```

## JSON (Caddy API) Configuration

```json
{
  "apps": {
    "tls": {
      "automation": {
        "policies": [
          {
            "issuers": [
              {
                "module": "acme",
                "challenges": {
                  "dns": {
                    "provider": {
                      "name": "plesk",
                      "base_url": "https://plesk.example.com/api/v2",
                      "api_key": "YOUR_API_KEY"
                    }
                  }
                }
              }
            ]
          }
        ]
      }
    }
  }
}
```

## How It Works

1. **GetRecords** — queries your Plesk server for existing DNS records in the zone.
2. **AppendRecords** — issues a POST to create the `_acme-challenge` TXT records.
3. **DeleteRecords** — cleans up TXT records after validation.

All HTTP calls go through Plesk’s REST API (`/dns/records` endpoints) with JSON request/response bodies.

## Troubleshooting

### API Access Denied

If you see:
```
plesk error: {"code":0,"message":"Access to API is disabled by admin access policy for <your IP>"}
```
log into Plesk → Tools & Settings → API → IP Access List and whitelist your Caddy server’s egress IP.
