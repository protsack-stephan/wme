package firehose

import "net/http"

// Client firehose streaming client to simplify work with WME realtime API.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient create new realtime client.
func NewClient() *Client {
	return &Client{
		BaseURL:    "https://firehose.enterprise.wikimedia.com/v1",
		HTTPClient: &http.Client{},
	}
}

func (c *Client) PageUpdate() {}

func (c *Client) PageDelete() {}

func (c *Client) PageVisibility() {}
