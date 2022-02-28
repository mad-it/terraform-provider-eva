package eva

import (
	"github.com/go-resty/resty/v2"
)

const (
	contentType = "application/json"
	userAgent   = "terraform-provider-eva"
)

type Client struct {
	Client *resty.Client
}

func NewClient(apiURL string) *Client {

	client := resty.New().
		SetBaseURL(apiURL).
		SetHeader("Content-Type", contentType).
		SetHeader("EVA-User-Agent", userAgent)

	return &Client{
		Client: client,
	}
}
