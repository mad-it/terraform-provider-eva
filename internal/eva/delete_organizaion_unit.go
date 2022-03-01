package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type DeleteOrganizationUnitRequest struct {
	ID int64 `json:"OrganizationUnitID"`
}

type DeleteOrganizationUnitResponse struct {
}

func (c *Client) DeleteOrganizationUnit(ctx context.Context, req DeleteOrganizationUnitRequest) (*DeleteOrganizationUnitResponse, error) {
	const (
		path = "/api/core/DeleteOrganizationUnit"
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

	var jsonResp DeleteOrganizationUnitResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		panic(err)
	}

	return &jsonResp, nil
}
