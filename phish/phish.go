package phish

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

type TrackOptions struct {
	UserAgent string
	ClientIP  string
}

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

func (c *Client) TrackOpen(rid string, opts ...TrackOptions) (*http.Response, error) {
	options := &TrackOptions{}
	for _, opt := range opts {
		if opt.ClientIP != "" {
			options.ClientIP = opt.ClientIP
		}
		if opt.UserAgent != "" {
			options.UserAgent = opt.UserAgent
		}
	}

	u, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}

	u = u.JoinPath("track")

	query := u.Query()
	query.Set("rid", rid)

	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if options.UserAgent != "" {
		req.Header.Set("User-Agent", options.UserAgent)
	}

	if options.ClientIP != "" {
		req.Header.Set("X-Forwarded-For", options.ClientIP)
	}

	return c.client.Do(req)
}

func (c *Client) TrackClick(rid string, opts ...TrackOptions) (*http.Response, error) {
	options := &TrackOptions{}
	for _, opt := range opts {
		if opt.ClientIP != "" {
			options.ClientIP = opt.ClientIP
		}
		if opt.UserAgent != "" {
			options.UserAgent = opt.UserAgent
		}
	}

	u, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Set("rid", rid)

	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if options.UserAgent != "" {
		req.Header.Set("User-Agent", options.UserAgent)
	}

	if options.ClientIP != "" {
		req.Header.Set("X-Forwarded-For", options.ClientIP)
	}

	return c.client.Do(req)
}
