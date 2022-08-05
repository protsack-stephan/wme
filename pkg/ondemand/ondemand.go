// SDK for on-demand APIs. Includes article lookup API and projects API.
package ondemand

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/protsack-stephan/wme/schema/v1"
)

// ArticleRequest parameters required for article request.
type ArticleRequest struct {
	Project string `json:"project"`
	Name    string `json:"name"`
}

// Client to simplify work with WME on-demand API.
type Client struct {
	BaseURL     string
	HTTPClient  *http.Client
	accessToken string
}

// NewClient creates new on-demand client.
func NewClient() *Client {
	return &Client{
		BaseURL:    "https://api.enterprise.wikimedia.com/v1",
		HTTPClient: &http.Client{},
	}
}

func (c *Client) get(ctx context.Context, url string, v interface{}) (*http.Response, error) {
	body := bytes.NewBuffer([]byte{})

	if err := json.NewEncoder(body).Encode(v); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.BaseURL, url), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusIMUsed {
		data, err := ioutil.ReadAll(res.Body)

		if err != nil {
			return nil, err
		}

		defer res.Body.Close()
		return nil, fmt.Errorf("%s: %s", res.Status, string(data))
	}

	return res, nil
}

// SetAccessToken sets access token for authentication.
func (c *Client) SetAccessToken(accessToken string) {
	c.accessToken = accessToken
}

// GetAccessToken returns value of the access token.
func (c *Client) GetAccessToken() string {
	return c.accessToken
}

// Article triggers /pages/meta/{project}/{name} endpoint and returns current revision of an article.
func (c *Client) Article(ctx context.Context, req *ArticleRequest) (*schema.Page, error) {
	res, err := c.get(ctx, fmt.Sprintf("/pages/meta/%s/%s", req.Project, req.Name), nil)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	art := new(schema.Page)

	if err := json.NewDecoder(res.Body).Decode(art); err != nil {
		return nil, err
	}

	return art, nil
}

// Projects triggers /projects endpoint and returns list of available projects.
func (c *Client) Projects(ctx context.Context) ([]*schema.Project, error) {
	res, err := c.get(ctx, "/projects", nil)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var projects []*schema.Project

	dec := json.NewDecoder(res.Body)

	// read open bracket
	if _, err = dec.Token(); err != nil {
		return nil, err
	}

	// while the array contains project objects
	for dec.More() {
		project := new(schema.Project)

		if err := dec.Decode(project); err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	// read closing bracket
	if _, err = dec.Token(); err != nil {
		return nil, err
	}

	return projects, nil
}
