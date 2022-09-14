// Package realtime is a SDK for working with realtime API v2 BETA.
package realtime

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/protsack-stephan/wme/schema/v2"
)

// Filter payload for filters in realtime API.
type Filter struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

// ArticlesRequest request for filtering and fields in realtime API.
type ArticlesRequest struct {
	Since   time.Time `json:"since,omitempty"`
	Fields  []string  `json:"fields,omitempty"`
	Filters []Filter  `json:"filters,omitempty"`
}

// NewClient create new realtime client.
func NewClient() *Client {
	return &Client{
		BaseURL:    "https://realtime-beta.enterprise.wikimedia.com/v2",
		HTTPClient: &http.Client{},
	}
}

// Client realtime streaming client to simplify work with WME realtime API.
type Client struct {
	BaseURL     string
	HTTPClient  *http.Client
	accessToken string
}

// SetAccessToken sets access token for authentication.
func (c *Client) SetAccessToken(accessToken string) {
	c.accessToken = accessToken
}

// GetAccessToken returns value of the access token.
func (c *Client) GetAccessToken() string {
	return c.accessToken
}

// Articles opens and listens articles stream.
func (cl *Client) Articles(ctx context.Context, req *ArticlesRequest, cb func(art *schema.Article) error) error {
	return cl.subscribe(ctx, "/articles", req, func(data []byte) error {
		art := new(schema.Article)

		if err := json.Unmarshal(data, art); err != nil {
			return err
		}

		return cb(art)
	})
}

func (c *Client) subscribe(ctx context.Context, url string, body interface{}, cb func(data []byte) error) error {
	bod := bytes.NewBuffer([]byte{})

	if body != nil {
		data, err := json.Marshal(body)

		if err != nil {
			return err
		}

		if _, err := bod.Write(data); err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.BaseURL, url), bod)

	if err != nil {
		return err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "application/x-ndjson")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusIMUsed {
		data, err := ioutil.ReadAll(res.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("%s: %s", res.Status, string(data))
	}

	scn := bufio.NewScanner(res.Body)
	scn.Buffer([]byte{}, 20971520) // this is important as we are encountering large messages (approx 20MB)

	for scn.Scan() {
		if err := cb([]byte(scn.Text())); err != nil {
			return err
		}
	}

	return nil
}
