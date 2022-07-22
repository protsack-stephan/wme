package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// LoginRequest parameters required for login request.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse payload for the login response.
type LoginResponse struct {
	IDToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// RevokeTokenRequest parameters required for revoke token request.
type RevokeTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenRequest parameters required for refresh token request.
type RefreshTokenRequest struct {
	Username     string `json:"username"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse payload for the refresh token response.
type RefreshTokenResponse struct {
	IDToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// NewClient create new http client for WME authentication.
func NewClient() *Client {
	return &Client{
		BaseURL:    "https://auth.enterprise.wikimedia.com/v1",
		HTTPClient: &http.Client{},
	}
}

// Client http client to simplify work with WME authentication.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func (c *Client) post(url string, v interface{}) (*http.Response, error) {
	body := bytes.NewBuffer([]byte{})

	if err := json.NewEncoder(body).Encode(v); err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Post(fmt.Sprintf("%s%s", c.BaseURL, url), "application/json", body)

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

// Login triggers login endpoint and returns fresh set of tokens.
func (c *Client) Login(req *LoginRequest) (*LoginResponse, error) {
	res, err := c.post("/login", req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	rsp := new(LoginResponse)

	if err := json.NewDecoder(res.Body).Decode(rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

// RevokeToken invalidates refresh token and all related access tokens.
func (c *Client) RevokeToken(req *RevokeTokenRequest) error {
	_, err := c.post("/token-revoke", req)
	return err
}

// RefreshToken gets new set of tokens using refresh token.
func (c *Client) RefreshToken(req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	res, err := c.post("/token-refresh", req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	rsp := new(RefreshTokenResponse)

	if err := json.NewDecoder(res.Body).Decode(rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}
