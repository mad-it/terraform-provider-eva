package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

func (c *Client) Post(ctx context.Context, url string, requestBody interface{}, responseBody *interface{}) error {
	resp, err := c.restClient.R().
		SetBody(requestBody).
		Post(url)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp CreateAccountingRecipeResponse
	if err := json.Unmarshal([]byte(resp.Body()), &responseBody); err != nil || jsonResp.HasErrors {

		return errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return nil
}
