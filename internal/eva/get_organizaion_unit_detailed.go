package eva

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type GetOrganizationUnitDetailedRequest struct {
	ID int64 `json:"ID"`
}

type GetOrganizationUnitDetailedResponse struct {
	ID           int64  `json:"ID"`
	Name         string `json:"Name"`
	PhoneNumber  string `json:"PhoneNumber"`
	EmailAddress string `json:"EmailAddress"`
	ParentID     int64  `json:"ParentID"`
	CurrencyID   string `json:"CurrencyID"`
	BackendID    string `json:"BackendID"`
}

func (c *Client) GetOrganizationUnitDetailed(ctx context.Context, req GetOrganizationUnitDetailedRequest) (*GetOrganizationUnitDetailedResponse, error) {
	const (
		path = "/api/core/GetOrganizationUnitDetailed"
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

		return nil, errors.New("Request failed.")
	}

	tflog.Debug(ctx, "Request info", "Status code", resp.StatusCode(), "body", resp.String())

	var jsonResp GetOrganizationUnitDetailedResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		panic(err)
	}

	return &jsonResp, nil
}
