package phish

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	url    string
	client *http.Client
}

func NewClient(url string) *Client {
	return &Client{
		url: url,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (c *Client) TrackOpen(rid string) (*http.Response, error) {
	u, err := url.JoinPath(c.url, fmt.Sprintf("track?rid=%s", rid))
	if err != nil {
		return nil, err
	}

	return c.client.Get(u)
}

func (c *Client) TrackClick(rid string) (*http.Response, error) {
	u, err := url.JoinPath(c.url, fmt.Sprintf("?rid=%s", rid))
	if err != nil {
		return nil, err
	}

	return c.client.Get(u)
}
