package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type CreateOrganizationUnitRequest struct {
	Name                string `json:"Name"`
	PhoneNumber         string `json:"PhoneNumber"`
	EmailAddress        string `json:"EmailAddress"`
	ParentID            int64  `json:"ParentID"`
	CurrencyID          string `json:"CurrencyID"`
	BackendID           string `json:"BackendID"`
	CostPriceCurrencyID string `json:"CostPriceCurrencyID"`
}

type createOrganizationUnitRequest struct {
	ToCreate CreateOrganizationUnitRequest `json:"ToCreate"`
}

type CreateOrganizationUnitResponse struct {
	ID int64 `json:"ID"`
}

func (c *Client) CreateOrganizationUnit(ctx context.Context, req CreateOrganizationUnitRequest) (*CreateOrganizationUnitResponse, error) {
	const (
		path = "/api/core/CreateOrganizationUnit"
	)

	requestBody := &createOrganizationUnitRequest{
		ToCreate: req,
	}

	resp, err := c.Client.R().
		SetBody(requestBody).
		Post(path)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	tflog.Debug(ctx, "Request info", "Status code", resp.StatusCode(), "body", resp.String())

	var jsonResp CreateOrganizationUnitResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		panic(err)
	}

	return &jsonResp, nil
}
