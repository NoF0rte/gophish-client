package phish

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
)

type TrackOption func(*trackOptions)
type trackOptions struct {
	userAgent string
	clientIP  string
}

func WithUserAgent(userAgent string) TrackOption {
	return func(to *trackOptions) {
		to.userAgent = userAgent
	}
}

func WithClientIP(ip string) TrackOption {
	return func(to *trackOptions) {
		to.clientIP = ip
	}
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

func (c *Client) TrackOpen(rid string, opts ...TrackOption) (*http.Response, error) {
	options := &trackOptions{}
	for _, opt := range opts {
		opt(options)
	}

	u, err := url.JoinPath(c.url, fmt.Sprintf("track?rid=%s", rid))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	if options.userAgent != "" {
		req.Header.Set("User-Agent", options.userAgent)
	}

	if options.clientIP != "" {
		req.Header.Set("X-Forwarded-For", options.clientIP)
	}

	return c.client.Do(req)
}

func (c *Client) TrackClick(rid string, opts ...TrackOption) (*http.Response, error) {
	options := &trackOptions{}
	for _, opt := range opts {
		opt(options)
	}

	u, err := url.JoinPath(c.url, fmt.Sprintf("?rid=%s", rid))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	if options.userAgent != "" {
		req.Header.Set("User-Agent", options.userAgent)
	}

	if options.clientIP != "" {
		req.Header.Set("X-Forwarded-For", options.clientIP)
	}

	return c.client.Do(req)
}
