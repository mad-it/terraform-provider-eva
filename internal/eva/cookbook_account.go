package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	createCookbookAccountPath = "/api/core/management/CreateAccount"
	getCookbookAccountPath    = "/api/core/management/GetAccount"
	updateCookbookAccountPath = "/api/core/management/UpdateAccount"
	deleteCookbookAccountPath = "/api/core/management/DeleteAccount"
)

type CreateCookbookAccountRequest struct {
	Name          string `json:"Name"`
	ObjectAccount string `json:"ObjectAccount"`
	BookingFlags  int64  `json:"BookingFlags"`
	Type          int64  `json:"Type"`
}

type CreateCookbookAccountResponse struct {
	ID int64
}

func (c *Client) CreateCookbookAccount(ctx context.Context, req CreateCookbookAccountRequest) (*CreateCookbookAccountResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(createCookbookAccountPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp CreateCookbookAccountResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type GetCookbookAccountRequest struct {
	ID int64 `json:"ID"`
}

type GetCookbookAccountResponse struct {
	ID            int64  `json:"ID"`
	Name          string `json:"Name"`
	ObjectAccount string `json:"ObjectAccount"`
	BookingFlags  int64  `json:"BookingFlags"`
	Type          int64  `json:"Type"`
}

func (c *Client) GetCookbookAccount(ctx context.Context, req GetCookbookAccountRequest) (*GetCookbookAccountResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(getCookbookAccountPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp GetCookbookAccountResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type UpdateCookbookAccountRequest struct {
	ID            int64  `json:"ID"`
	Name          string `json:"Name"`
	ObjectAccount string `json:"ObjectAccount"`
	BookingFlags  int64  `json:"BookingFlags"`
	Type          int64  `json:"Type"`
}

func (c *Client) UpdateCookbookAccount(ctx context.Context, req UpdateCookbookAccountRequest) (*EmptyResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(updateCookbookAccountPath)

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

type DeleteCookbookAccountRequest struct {
	ID int64 `json:"ID"`
}

func (c *Client) DeleteCookbookAccount(ctx context.Context, req DeleteCookbookAccountRequest) (*EmptyResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(deleteCookbookAccountPath)

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
