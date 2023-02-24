package phish

import (
	"crypto/tls"
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

	if options.userAgent != "" {
		req.Header.Set("User-Agent", options.userAgent)
	}

	if options.clientIP != "" {
		req.Header.Set("X-Forwarded-For", options.clientIP)
	}

	return c.client.Do(req)
}
