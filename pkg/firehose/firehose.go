package firehose

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/protsack-stephan/wme/schema/v1"
)

// EventID shows metadata for the event.
type EventID struct {
	Topic     string    `json:"topic"`
	Partition int       `json:"partition"`
	Dt        time.Time `json:"dt"`
	Timestamp int64     `json:"timestamp"`
	Offset    int       `json:"offset"`
}

// Event server side event structure for firehose.
type Event struct {
	ID   []*EventID   `json:"id"`
	Data *schema.Page `json:"data"`
}

// Client firehose streaming client to simplify work with WME realtime API.
type Client struct {
	BaseURL     string
	HTTPClient  *http.Client
	accessToken string
}

// NewClient create new realtime client.
func NewClient() *Client {
	return &Client{
		BaseURL:    "https://firehose.enterprise.wikimedia.com/v1",
		HTTPClient: &http.Client{},
	}
}

func (c *Client) subscribe(ctx context.Context, since time.Time, url string, cb func(evt *Event)) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s?since=%s", c.BaseURL, url, since.UTC().Format(time.RFC3339)), nil)

	if err != nil {
		return err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusIMUsed {
		data, err := io.ReadAll(res.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("%s: %s", res.Status, string(data))
	}

	scn := bufio.NewScanner(res.Body)
	buf := []byte{}
	scn.Buffer(buf, 20971520) // this is important as we are encountering large messages (approx 20MB)

	evt := new(Event)

	for scn.Scan() {
		if strings.HasPrefix(scn.Text(), "id:") {
			if err := json.Unmarshal([]byte(scn.Text()[len("id:"):]), &evt.ID); err != nil {
				return err
			}
		}

		if strings.HasPrefix(scn.Text(), "data:") {
			evt.Data = new(schema.Page)

			if err := json.Unmarshal([]byte(scn.Text()[len("data:"):]), evt.Data); err != nil {
				return err
			}
		}

		if len(evt.ID) > 0 && evt.Data != nil {
			cb(evt)
			evt = new(Event)
		}
	}

	return nil
}

// SetAccessToken sets access token for authentication.
func (c *Client) SetAccessToken(accessToken string) {
	c.accessToken = accessToken
}

// GetAccessToken returns value of the access token.
func (c *Client) GetAccessToken() string {
	return c.accessToken
}

// PageUpdate opens connection to page update stream.
func (c *Client) PageUpdate(ctx context.Context, since time.Time, cb func(evt *Event)) error {
	return c.subscribe(ctx, since, "/page-update", cb)
}

// PageDelete opens connection to page delete stream.
func (c *Client) PageDelete(ctx context.Context, since time.Time, cb func(evt *Event)) error {
	return c.subscribe(ctx, since, "/page-delete", cb)
}

// PageVisibility opens connection to page visibility stream.
func (c *Client) PageVisibility(ctx context.Context, since time.Time, cb func(evt *Event)) error {
	return c.subscribe(ctx, since, "/page-visibility", cb)
}
