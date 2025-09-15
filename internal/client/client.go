package client

import (
	"bytes"
	"net/http"
	"qatarina-cli/internal/auth"
)

type Client struct {
	BaseURL string
	Token   string
}

func Default() *Client {
	token, _ := auth.LoadToken()
	return &Client{
		BaseURL: "http://localhost:4597/", // make this configurable later
		Token:   token,
	}
}

func (c *Client) Post(path string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.BaseURL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	return http.DefaultClient.Do(req)
}
