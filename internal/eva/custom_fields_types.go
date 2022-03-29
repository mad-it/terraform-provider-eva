package eva

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	getCustomFieldTypesPath = "/api/core/management/GetCustomFieldTypes"
)

type CustomFieldType struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

type GetCustomFieldTypesResponse struct {
	CustomFieldTypes []CustomFieldType `json:"CustomFieldTypes"`
}

func (c *Client) GetCustomFieldTypes(ctx context.Context, req Empty) (*CreateCustomOrderStatusResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(getCustomFieldTypesPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, fmt.Errorf("request failed with error: %s", resp.String())
	}

	var jsonResp CreateCustomOrderStatusResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, fmt.Errorf("response could not be parsed. Received: %s", resp.String())
	}

	return &jsonResp, nil
}
