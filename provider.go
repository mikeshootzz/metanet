package plesk

import (
	"context"
	"fmt"
	"net/url"

	"github.com/caddyserver/caddy/v2"
	"github.com/libdns/libdns"
)

// Provider configures the Plesk DNS provider.
type Provider struct {
	BaseURL  string `json:"base_url,omitempty"`
	APIKey   string `json:"api_key,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns module info.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.plesk",
		New: func() caddy.Module { return new(Provider) },
	}
}

// Provision sets defaults.
func (p *Provider) Provision(ctx caddy.Context) error {
	if p.BaseURL == "" {
		p.BaseURL = "https://localhost/api/v2"
	}
	return nil
}

// GetRecords lists all records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	client := NewClient(p.BaseURL, p.Username, p.Password, p.APIKey)
	vals := url.Values{"domain": {zone}}
	var recs []struct {
		ID    int    `json:"id"`
		Type  string `json:"type"`
		Host  string `json:"host"`
		Value string `json:"value"`
		TTL   int    `json:"ttl"`
	}
	if err := client.doRequest("GET", "/dns/records", vals, nil, &recs); err != nil {
		return nil, err
	}
	out := make([]libdns.Record, len(recs))
	for i, r := range recs {
		out[i] = libdns.Record{
			ID:    fmt.Sprint(r.ID),
			Type:  r.Type,
			Name:  r.Host,
			Value: r.Value,
			TTL:   r.TTL,
		}
	}
	return out, nil
}

// AppendRecord creates a new record in the zone.
func (p *Provider) AppendRecord(ctx context.Context, zone string, rec libdns.Record) (libdns.Record, error) {
	client := NewClient(p.BaseURL, p.Username, p.Password, p.APIKey)
	vals := url.Values{"domain": {zone}}
	body := map[string]interface{}{
		"type":  rec.Type,
		"host":  rec.Name,
		"value": rec.Value,
		"ttl":   rec.TTL,
	}
	var resp struct {
		ID int `json:"id"`
	}
	if err := client.doRequest("POST", "/dns/records", vals, body, &resp); err != nil {
		return rec, err
	}
	rec.ID = fmt.Sprint(resp.ID)
	return rec, nil
}

// DeleteRecord removes a record by ID.
func (p *Provider) DeleteRecord(ctx context.Context, zone, id string) error {
	client := NewClient(p.BaseURL, p.Username, p.Password, p.APIKey)
	return client.doRequest("DELETE", "/dns/records/"+id, nil, nil, nil)
}

var _ libdns.RecordGetter = (*Provider)(nil)
var _ libdns.RecordAppender = (*Provider)(nil)
var _ libdns.RecordRemover = (*Provider)(nil)
