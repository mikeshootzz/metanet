package main

import (
	"context"
	"fmt"
	"net/url"

	caddy_dns "github.com/caddyserver/caddy/v2/modules/caddy-dns"
)

// Ensure provider implements interface
var _ caddy_dns.DNSProvider = (*DNSProvider)(nil)

// Records retrieves all DNS records for the zone.
func (p *DNSProvider) Records(zone string) ([]caddy_dns.Record, error) {
	client := NewClient(p.BaseURL, p.Username, p.Password, p.APIKey)
	// Plesk API: GET /dns/records?domain=example.com
	vals := url.Values{}
	vals.Set("domain", zone)
	var recs []struct {
		ID    int    `json:"id"`
		Type  string `json:"type"`
		Host  string `json:"host"`
		Value string `json:"value"`
		Opt   string `json:"opt,omitempty"`
		TTL   int    `json:"ttl"`
	}
	if err := client.doRequest("GET", "/dns/records", vals, nil, &recs); err != nil {
		return nil, err
	}
	records := make([]caddy_dns.Record, len(recs))
	for i, r := range recs {
		records[i] = caddy_dns.Record{
			ID:       fmt.Sprint(r.ID),
			Type:     r.Type,
			Name:     r.Host,
			Value:    r.Value,
			TTL:      r.TTL,
			Priority: 0,
			Port:     0,
			Weight:   0,
			Flags:    0,
		}
	}
	return records, nil
}

func (p *DNSProvider) GetZoneRecords(ctx context.Context, zone string) ([]caddy_dns.Record, error) {
	return p.Records(zone)
}

func (p *DNSProvider) CreateZoneRecord(ctx context.Context, zone string, rec caddy_dns.Record) (caddy_dns.Record, error) {
	client := NewClient(p.BaseURL, p.Username, p.Password, p.APIKey)
	vals := url.Values{}
	vals.Set("domain", zone)
	reqBody := map[string]interface{}{
		"type":  rec.Type,
		"host":  rec.Name,
		"value": rec.Value,
		"ttl":   rec.TTL,
	}
	var created struct {
		ID int `json:"id"`
	}
	if err := client.doRequest("POST", "/dns/records", vals, reqBody, &created); err != nil {
		return rec, err
	}
	rec.ID = fmt.Sprint(created.ID)
	return rec, nil
}

func (p *DNSProvider) UpdateZoneRecord(ctx context.Context, zone string, existingID string, rec caddy_dns.Record) (caddy_dns.Record, error) {
	client := NewClient(p.BaseURL, p.Username, p.Password, p.APIKey)
	path := "/dns/records/" + existingID
	reqBody := map[string]interface{}{
		"type":  rec.Type,
		"host":  rec.Name,
		"value": rec.Value,
		"ttl":   rec.TTL,
	}
	if err := client.doRequest("PUT", path, nil, reqBody, nil); err != nil {
		return rec, err
	}
	rec.ID = existingID
	return rec, nil
}

func (p *DNSProvider) DeleteZoneRecord(ctx context.Context, zone string, id string) error {
	client := NewClient(p.BaseURL, p.Username, p.Password, p.APIKey)
	path := "/dns/records/" + id
	return client.doRequest("DELETE", path, nil, nil, nil)
}
