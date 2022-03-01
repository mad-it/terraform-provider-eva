package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type UpdateOrganizationUnitRequest struct {
	ID                  int64  `json:"ID"`
	Name                string `json:"Name,omitempty"`
	PhoneNumber         string `json:"PhoneNumber,omitempty"`
	EmailAddress        string `json:"EmailAddress,omitempty"`
	ParentID            int64  `json:"ParentID,omitempty"`
	CurrencyID          string `json:"CurrencyID,omitempty"`
	BackendID           string `json:"BackendID,omitempty"`
	CostPriceCurrencyID string `json:"CostPriceCurrencyID,omitempty"`
}

type UpdateOrganizationUnitResponse struct {
}

func (c *Client) UpdateOrganizationUnit(ctx context.Context, req UpdateOrganizationUnitRequest) (*UpdateOrganizationUnitResponse, error) {
	const (
		path = "/api/core/UpdateOrganizationUnit"
	)

	resp, err := c.Client.R().
		SetBody(req).
		Post(path)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp UpdateOrganizationUnitResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		panic(err)
	}

	return &jsonResp, nil
}
