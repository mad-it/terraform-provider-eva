package eva

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type CustomFieldDataTypes int

const (
	String CustomFieldDataTypes = iota
	Bool
	Integer
	Decimal
	Enum
	Text
	DateTime
	Date

	createCustomFieldPath = "/api/core/management/CreateCustomField"
	updateCustomFieldPath = "/api/core/management/UpdateCustomField"
	deleteCustomFieldPath = "/api/core/management/DeleteCustomField"
	getCustomFieldPath    = "/api/core/management/GetCustomFieldByID"
)

type CustomFieldOptions struct {
	IsArray       bool              `json:"IsArray,omitempty"`
	IsRequired    bool              `json:"IsRequired,omitempty"`
	MinimumValue  int64             `json:"MinimumValue,omitempty"`
	MaximumValue  int64             `json:"MaximumValue,omitempty"`
	MinimumLength int64             `json:"MinimumLength,omitempty"`
	MaximumLength int64             `json:"MaximumLength,omitempty"`
	MinimumDate   string            `json:"MinimumDate,omitempty"`
	MaximumDate   string            `json:"MaximumDate,omitempty"`
	DefaultValue  string            `json:"DefaultValue,omitempty"`
	EnumValues    map[string]string `json:"EnumValues,omitempty"`
}

type CreateCustomFieldRequest struct {
	TypeID              int64                `json:"TypeID,omitempty"`
	TypeKey             string               `json:"TypeKey,omitempty"`
	Name                string               `json:"Name,omitempty"`
	DisplayName         string               `json:"DisplayName,omitempty"`
	DataType            CustomFieldDataTypes `json:"DataType,omitempty"`
	Order               int64                `json:"Order,omitempty"`
	Options             CustomFieldOptions   `json:"Options,omitempty"`
	BackendID           string               `json:"BackendID,omitempty"`
	VisibleByUserTypes  int64                `json:"VisibleByUserTypes,omitempty"`
	EditableByUserTypes int64                `json:"EditableByUserTypes,omitempty"`
}

type CreateCustomFieldResponse struct {
	ID int64 `json:"ID"`
}

func (c *Client) CreateCustomField(ctx context.Context, req CreateCustomFieldRequest) (*CreateCustomFieldResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(createCustomFieldPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, fmt.Errorf("request failed with error: %s", resp.String())
	}

	var jsonResp CreateCustomFieldResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, fmt.Errorf("response could not be parsed. Received: %s", resp.String())
	}

	return &jsonResp, nil
}

type UpdateCustomFieldRequest struct {
	ID                  int64                `json:"ID"`
	TypeID              int64                `json:"TypeID,omitempty"`
	TypeKey             string               `json:"TypeKey,omitempty"`
	Name                string               `json:"Name,omitempty"`
	DisplayName         string               `json:"DisplayName,omitempty"`
	DataType            CustomFieldDataTypes `json:"DataType,omitempty"`
	Order               int64                `json:"Order,omitempty"`
	Options             CustomFieldOptions   `json:"Options,omitempty"`
	BackendID           string               `json:"BackendID,omitempty"`
	VisibleByUserTypes  int64                `json:"VisibleByUserTypes,omitempty"`
	EditableByUserTypes int64                `json:"EditableByUserTypes,omitempty"`
}

func (c *Client) UpdateCustomField(ctx context.Context, req UpdateCustomFieldRequest) (*Empty, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(updateCustomFieldPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, fmt.Errorf("request failed with error: %s", resp.String())
	}

	var jsonResp Empty
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, fmt.Errorf("response could not be parsed. Received: %s", resp.String())
	}

	return &jsonResp, nil
}

type GetCustomFieldByIDRequest struct {
	ID int64 `json:"ID"`
}

func (c *Client) GetCustomField(ctx context.Context, req GetCustomFieldByIDRequest) (*UpdateCustomFieldRequest, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(getCustomFieldPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, fmt.Errorf("request failed with error: %s", resp.String())
	}

	var jsonResp UpdateCustomFieldRequest
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, fmt.Errorf("response could not be parsed. Received: %s", resp.String())
	}

	return &jsonResp, nil
}

type DeleteCustomFieldRequest struct {
	ID int64 `json:"ID"`
}

func (c *Client) DeleteCustomField(ctx context.Context, req DeleteCustomFieldRequest) (*Empty, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(getCustomFieldPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, fmt.Errorf("request failed with error: %s", resp.String())
	}

	var jsonResp Empty
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, fmt.Errorf("response could not be parsed. Received: %s", resp.String())
	}

	return &jsonResp, nil
}
