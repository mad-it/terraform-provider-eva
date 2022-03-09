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

func NewClient(apiURL string, debugMode bool) *Client {

	restClient := resty.New().
		SetBaseURL(apiURL).
		SetHeader("Content-Type", contentType).
		SetHeader("EVA-User-Agent", userAgent).
		SetDebug(debugMode)

	return &Client{
		restClient: restClient,
	}
}

func (c *Client) SetAuthorizationHeader(token string) {
	c.restClient.SetHeader("authorization", token)
}
