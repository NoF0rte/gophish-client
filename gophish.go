package gophish

import (
	"github.com/NoF0rte/gophish-client/api"
	"github.com/NoF0rte/gophish-client/phish"
)

type Client struct {
	API   *api.Client
	Phish *phish.Client
}

func NewClient(phishURL string, adminURL string, apiKey string) *Client {
	return &Client{
		API:   api.NewClient(adminURL, apiKey),
		Phish: phish.NewClient(phishURL),
	}
}

func NewClientFromCredentials(phishURL string, adminURL string, username string, password string) (*Client, error) {
	apiClient, err := api.NewClientFromCredentials(adminURL, username, password)
	if err != nil {
		return nil, err
	}

	return &Client{
		API:   apiClient,
		Phish: phish.NewClient(phishURL),
	}, nil
}
