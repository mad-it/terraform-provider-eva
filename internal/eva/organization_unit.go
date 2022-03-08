package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	getOrganizationUnitPath    = "/api/core/GetOrganizationUnitDetailed"
	createOrganizationUnitPath = "/api/core/CreateOrganizationUnit"
	deleteOrganizationUnitPath = "/api/core/DeleteOrganizationUnit"
	updateOrganizationUnitPath = "/api/core/UpdateOrganizationUnit"
)

type address struct {
	HouseNumber string  `json:"HouseNumber,omitempty"`
	Address1    string  `json:"Address1,omitempty"`
	Address2    string  `json:"Address2,omitempty"`
	ZipCode     string  `json:"ZipCode,omitempty"`
	City        string  `json:"City,omitempty"`
	CountryID   string  `json:"CountryID,omitempty"`
	Latitude    float64 `json:"Latitude,omitempty"`
	Longitude   float64 `json:"Longitude,omitempty"`
}

type CreateOrUpdateOrganizationUnitRequest struct {
	ID                  int64   `json:"ID,omitempty"`
	Name                string  `json:"Name,omitempty"`
	PhoneNumber         string  `json:"PhoneNumber,omitempty"`
	EmailAddress        string  `json:"EmailAddress,omitempty"`
	ParentID            int64   `json:"ParentID,omitempty"`
	CurrencyID          string  `json:"CurrencyID,omitempty"`
	BackendID           string  `json:"BackendID,omitempty"`
	CostPriceCurrencyID string  `json:"CostPriceCurrencyID,omitempty"`
	TypeID                  int64   `json:"ID,omitempty"`
	Address             address `json:"Address,omitempty"`
}

type createOrganizationUnitRequest struct {
	ToCreate CreateOrUpdateOrganizationUnitRequest `json:"ToCreate"`
}

type CreateOrganizationUnitResponse struct {
	ID int64 `json:"ID"`
}

func (c *Client) CreateOrganizationUnit(ctx context.Context, req CreateOrUpdateOrganizationUnitRequest) (*CreateOrganizationUnitResponse, error) {
	requestBody := &createOrganizationUnitRequest{
		ToCreate: req,
	}

	resp, err := c.restClient.R().
		SetBody(requestBody).
		Post(createOrganizationUnitPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp CreateOrganizationUnitResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type UpdateOrganizationUnitResponse struct {
}

func (c *Client) UpdateOrganizationUnit(ctx context.Context, req CreateOrUpdateOrganizationUnitRequest) (*UpdateOrganizationUnitResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(updateOrganizationUnitPath)

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

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

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

	resp, err := c.restClient.R().
		SetBody(req).
		Post(getOrganizationUnitPath)

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

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type DeleteOrganizationUnitRequest struct {
	ID int64 `json:"OrganizationUnitID"`
}

type DeleteOrganizationUnitResponse struct {
}

func (c *Client) DeleteOrganizationUnit(ctx context.Context, req DeleteOrganizationUnitRequest) (*DeleteOrganizationUnitResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(deleteOrganizationUnitPath)

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

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}
