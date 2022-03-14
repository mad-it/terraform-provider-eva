package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type CreateEmployeeResult int

const (
	CreatedNewUser CreateEmployeeResult = iota
	UpgradedExistingUser
	UpdatedExistingUser

	getUserPath        = "/api/core/GetUser"
	createEmployeePath = "/api/core/management/CreateEmployeeUser"
	updateUserPath     = "/api/core/UpdateUser"
	deleteUserPath     = "/api/core/DeleteUser"
)

type CreateEmployeeUserRequest struct {
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	EmailAddress string `json:"EmailAddress"`
	Password     string `json:"Password"`
}

type CreateEmployeeUserResponse struct {
	ID     int64                `json:"UserID"`
	Result CreateEmployeeResult `json:"Result"`
}

func (c *Client) CreateEmployee(ctx context.Context, req CreateEmployeeUserRequest) (*CreateEmployeeUserResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(createEmployeePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp CreateEmployeeUserResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type GetUserRequest struct {
	ID int64 `json:"ID"`
}

type GetEmployeeResponse struct {
	ID           int64  `json:"ID"`
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	EmailAddress string `json:"EmailAddress"`
}

func (c *Client) GetUser(ctx context.Context, req GetUserRequest) (*GetEmployeeResponse, error) {

	resp, err := c.restClient.R().
		SetBody(req).
		Post(getUserPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New("Request failed.")
	}

	tflog.Debug(ctx, "Request info", "Status code", resp.StatusCode(), "body", resp.String())

	var jsonResp GetEmployeeResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Error: %s \n Received: %s", err, resp.String()))
	}

	return &jsonResp, nil
}

type UpdateUserRequest struct {
	ID           int64  `json:"ID"`
	FirstName    string `json:"FirstName,omitempty"`
	LastName     string `json:"LastName,omitempty"`
	EmailAddress string `json:"EmailAddress,omitempty"`
}

func (c *Client) UpdateUser(ctx context.Context, req UpdateUserRequest) (*EmptyResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(updateUserPath)

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

type DeleteUserRequest struct {
	ID int64 `json:"ID"`
}

func (c *Client) DeleteUser(ctx context.Context, req DeleteUserRequest) (*EmptyResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(deleteUserPath)

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
