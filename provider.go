package plesk

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/libdns/libdns"
)

// Provider configures DNS updates against Plesk.
type Provider struct {
	BaseURL  string `json:"base_url,omitempty"`
	APIKey   string `json:"api_key,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the module information.
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

// GetRecords lists all DNS records in the zone.
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
	var out []libdns.Record
	for _, r := range recs {
		rr := libdns.RR{
			Name: r.Host,
			Type: r.Type,
			TTL:  time.Duration(r.TTL) * time.Second,
			Data: r.Value,
		}
		out = append(out, rr)
	}
	return out, nil
}

// AppendRecords creates new records in the zone.
func (p *Provider) AppendRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	client := NewClient(p.BaseURL, p.Username, p.Password, p.APIKey)
	vals := url.Values{"domain": {zone}}
	var created []libdns.Record
	for _, rec := range recs {
		rr := rec.RR()
		body := map[string]interface{}{
			"type":  rr.Type,
			"host":  rr.Name,
			"value": rr.Data,
			"ttl":   int(rr.TTL.Seconds()),
		}
		var resp struct {
			ID int `json:"id"`
		}
		if err := client.doRequest("POST", "/dns/records", vals, body, &resp); err != nil {
			return nil, err
		}
		created = append(created, rr)
	}
	return created, nil
}

// DeleteRecords removes records by matching fields and deletes via ID.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	client := NewClient(p.BaseURL, p.Username, p.Password, p.APIKey)
	vals := url.Values{"domain": {zone}}
	var recList []struct {
		ID    int    `json:"id"`
		Type  string `json:"type"`
		Host  string `json:"host"`
		Value string `json:"value"`
		TTL   int    `json:"ttl"`
	}
	if err := client.doRequest("GET", "/dns/records", vals, nil, &recList); err != nil {
		return nil, err
	}
	var deleted []libdns.Record
	for _, rec := range recs {
		rr := rec.RR()
		for _, r := range recList {
			if r.Type == rr.Type && r.Host == rr.Name && r.Value == rr.Data && r.TTL == int(rr.TTL.Seconds()) {
				path := fmt.Sprintf("/dns/records/%d", r.ID)
				if err := client.doRequest("DELETE", path, nil, nil, nil); err != nil {
					return nil, err
				}
				deleted = append(deleted, rr)
				break
			}
		}
	}
	return deleted, nil
}

var _ libdns.RecordGetter = (*Provider)(nil)
var _ libdns.RecordAppender = (*Provider)(nil)
var _ libdns.RecordDeleter = (*Provider)(nil)
