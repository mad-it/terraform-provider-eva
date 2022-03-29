package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	getSettingPath   = "/api/core/management/GetSetting"
	setSettingPath   = "/api/core/management/SetSetting"
	unsetSettingPath = "/api/core/management/UnsetSetting"
)

type SetSettingsRequest struct {
	Key                string `json:"Key"`
	Value              string `json:"Value"`
	OrganizationUnitID int64  `json:"OrganizationUnitID,omitempty"`
}

func (c *Client) SetSettings(ctx context.Context, req SetSettingsRequest) (*EmptyResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(setSettingPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp EmptyResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type GetSettingRequest struct {
	Key                string `json:"Key"`
	OrganizationUnitID int64  `json:"OrganizationUnitID,omitempty"`
}

type GetSettingResponse struct {
	Value string `json:"Value"`
}

func (c *Client) GetSetting(ctx context.Context, req GetSettingRequest) (*GetSettingResponse, error) {

	resp, err := c.restClient.R().
		SetBody(req).
		Post(getSettingPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New("Request failed.")
	}

	tflog.Debug(ctx, "Request info", "Status code", resp.StatusCode(), "body", resp.String())

	var jsonResp GetSettingResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type UnsetSettingsRequest struct {
	Key                string `json:"Key,omitempty"`
	OrganizationUnitID int64  `json:"OrganizationUnitID,omitempty"`
}

func (c *Client) UnsetSettings(ctx context.Context, req UnsetSettingsRequest) (*EmptyResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(unsetSettingPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp EmptyResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}
