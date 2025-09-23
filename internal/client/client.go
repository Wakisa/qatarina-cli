package client

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

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

// Default creates a new client that connects to default URL
// or host specified in the environment variable `QATARINA_HOST`
func Default() *Client {
	url := os.Getenv("QATARINA_HOST")
	if url == "" {
		url = "http://localhost:4597"
	}
	return NewClient(url)
}

func (c *Client) Post(path string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.BaseURL, path), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	return http.DefaultClient.Do(req)
}
