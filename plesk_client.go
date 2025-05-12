package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client wraps HTTP calls to Plesk.
type Client struct {
	BaseURL  string
	Username string
	Password string
	APIKey   string
	http     *http.Client
}

// NewClient constructs a new Plesk API client.
func NewClient(baseURL, username, password, apiKey string) *Client {
	return &Client{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
		APIKey:   apiKey,
		http:     http.DefaultClient,
	}
}

func (c *Client) doRequest(method, path string, query url.Values, body interface{}, result interface{}) error {
	u := fmt.Sprintf("%s%s?%s", c.BaseURL, path, query.Encode())
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
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		data, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("plesk error: %s", string(data))
	}
	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}
