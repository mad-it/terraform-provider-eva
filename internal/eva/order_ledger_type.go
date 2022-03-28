package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	createOrderLedgerTypePath = "/api/core/management/CreateOrderLedgerType"
	listOrderLedgerTypePath   = "/api/core/management/ListOrderLedgerTypes"
	updateOrderLedgerTypePath = "/api/core/management/UpdateOrderLedgerType"
	deleteOrderLedgerTypePath = "/api/core/management/DeleteOrderLedgerType"
)

type CreateOrderLedgerTypeRequest struct {
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
}

type CreateOrderLedgerTypeResponse struct {
	ID int64
}

func (c *Client) CreateOrderLedgerType(ctx context.Context, req CreateOrderLedgerTypeRequest) (*CreateOrderLedgerTypeResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(createOrderLedgerTypePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp CreateOrderLedgerTypeResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type OrderLedgerType struct {
	ID          int64  `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
}

type ListOrderLedgerTypeResponse struct {
	Result []OrderLedgerType `json:"OrderLedgerTypes"`
}

func (c *Client) ListOrderLedgerTypes(ctx context.Context) (*ListOrderLedgerTypeResponse, error) {
	resp, err := c.restClient.R().
		Post(listOrderLedgerTypePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New("Request failed.")
	}

	tflog.Debug(ctx, "Request info", "Status code", resp.StatusCode(), "body", resp.String())

	var jsonResp ListOrderLedgerTypeResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type UpdateOrderLedgerTypeRequest struct {
	ID          int64  `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

type UpdateOrderLedgerTypeResponse struct {
}

func (c *Client) UpdateOrderLedgerType(ctx context.Context, req UpdateOrderLedgerTypeRequest) (*UpdateOrderLedgerTypeResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(updateOrderLedgerTypePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp UpdateOrderLedgerTypeResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type DeleteOrderLedgerTypeRequest struct {
	ID int64 `json:"ID"`
}

type DeleteOrderLedgerTypeResponse struct {
}

func (c *Client) DeleteOrderLedgerType(ctx context.Context, req DeleteOrderLedgerTypeRequest) (*DeleteOrderLedgerTypeResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(deleteOrderLedgerTypePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp DeleteOrderLedgerTypeResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}
