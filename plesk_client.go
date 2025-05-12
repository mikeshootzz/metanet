package plesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client handles requests to the Plesk API.
type Client struct {
	BaseURL    string
	Username   string
	Password   string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new API client.
func NewClient(baseURL, username, password, apiKey string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
		APIKey:     apiKey,
		HTTPClient: http.DefaultClient,
	}
}

// doRequest makes an HTTP call to Plesk.
func (c *Client) doRequest(method, path string, params url.Values, body interface{}, out interface{}) error {
	u := c.BaseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return err
		}
	}
	req, err := http.NewRequest(method, u, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("X-API-Key", c.APIKey)
	}
	if c.Username != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("plesk error: %s", data)
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}
