package eva

import (
	"github.com/go-resty/resty/v2"
)

const (
	contentType = "application/json"
	userAgent   = "terraform-provider-eva"
)

type Client struct {
	restClient *resty.Client
}

func NewClient(apiURL string) *Client {

	restClient := resty.New().
		SetBaseURL(apiURL).
		SetHeader("Content-Type", contentType).
		SetHeader("EVA-User-Agent", userAgent).
		SetDebug(true)

	return &Client{
		restClient: restClient,
	}
}

func (c *Client) SetAuthorizationHeader(token string) {
	c.restClient.SetHeader("authorization", token)
}

type Empty struct{}
