package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	createCustomOrderStatusPath = "/api/core/CreateCustomOrderStatus"
	listCustomOrderStatusPath   = "/api/core/ListCustomOrderStatus"
	updateCustomOrderStatusPath = "/api/core/UpdateCustomOrderStatus"
	deleteCustomOrderStatusPath = "/api/core/DeleteCustomOrderStatus"
)

type CreateCustomOrderStatusRequest struct {
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
}

type CreateCustomOrderStatusResponse struct {
	ID int64
}

func (c *Client) CreateCustomOrderStatus(ctx context.Context, req CreateCustomOrderStatusRequest) (*CreateCustomOrderStatusResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(createCustomOrderStatusPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp CreateCustomOrderStatusResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type CustomOrderStatus struct {
	ID          int64  `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
}

type ListCustomOrderStatusResponse struct {
	Result []CustomOrderStatus `json:"Result"`
}

func (c *Client) ListCustomOrderStatus(ctx context.Context) (*ListCustomOrderStatusResponse, error) {
	resp, err := c.restClient.R().
		Post(listCustomOrderStatusPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New("Request failed.")
	}

	tflog.Debug(ctx, "Request info", "Status code", resp.StatusCode(), "body", resp.String())

	var jsonResp ListCustomOrderStatusResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type UpdateCustomOrderStatusRequest struct {
	ID          int64  `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
}

func (c *Client) UpdateCustomOrderStatus(ctx context.Context, req UpdateCustomOrderStatusRequest) (*Empty, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(updateCustomOrderStatusPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp Empty
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type DeleteCustomOrderStatusRequest struct {
	ID int64 `json:"ID"`
}

func (c *Client) DeleteCustomOrderStatus(ctx context.Context, req DeleteCustomOrderStatusRequest) (*Empty, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(deleteCustomOrderStatusPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp Empty
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}
