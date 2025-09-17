package client

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/wakisa/qatarina-cli/internal/auth"
)

type Client struct {
	BaseURL string
	Token   string
}

func NewClient(url string) *Client {
	token, err := auth.LoadToken()
	if err != nil {
		fmt.Println(err)
	}
	return &Client{
		BaseURL: url,
		Token:   token,
	}

}

func Default() *Client {
	return NewClient("http://localhost:4597/")
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
