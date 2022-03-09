package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	createRolePath                    = "/api/core/management/CreateRole"
	getRolePath                       = "/api/core/management/GetRole"
	updateRolePath                    = "/api/core/management/UpdateRole"
	deleteRolePath                    = "/api/core/management/DeleteRole"
	attachFunctionalitiesToRolePath   = "/api/core/management/AttachFunctionalitiesToRole"
	detachFunctionalitiesFromRolePath = "/api/core/management/DetachFunctionalitiesFromRole"
)

type CreateRoleRequest struct {
	Name     string `json:"Name"`
	UserType int64  `json:"UserType,omitempty"`
	Code     string `json:"Code,omitempty"`
}
type CreateRoleResponse struct {
	ID int64 `json:"ID"`
}

func (c *Client) CreateRole(ctx context.Context, req CreateRoleRequest) (*CreateRoleResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(createRolePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp CreateRoleResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type GetRoleRequest struct {
	ID int64 `json:"ID"`
}

type Role struct {
	Name                  string                    `json:"Name"`
	UserType              int64                     `json:"UserType,omitempty"`
	Code                  string                    `json:"Code,omitempty"`
	ScopedFunctionalities []RoleScopedFunctionality `json:"ScopedFunctionalities"`
}

type RoleScopedFunctionality struct {
	Functionality     string
	Scope             int64
	RequiresElevation bool
}

type GetRoleResponse struct {
	Result Role `json:"Result"`
}

func (c *Client) GetRole(ctx context.Context, req GetRoleRequest) (*GetRoleResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(getRolePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New("Request failed.")
	}

	tflog.Debug(ctx, "Request info", "Status code", resp.StatusCode(), "body", resp.String())

	var jsonResp GetRoleResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Error: %s \n Received: %s", err, resp.String()))
	}

	return &jsonResp, nil
}

type UpdateRoleRequest struct {
	ID       int64  `json:"ID"`
	Name     string `json:"Name,omitempty"`
	UserType int64  `json:"UserType,omitempty"`
	Code     string `json:"Code,omitempty"`
}

type UpdateRoleResponse struct {
}

func (c *Client) UpdateRole(ctx context.Context, req UpdateRoleRequest) (*UpdateRoleResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(updateRolePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp UpdateRoleResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type DeleteRoleRequest struct {
	ID int64 `json:"ID"`
}

type DeleteRoleResponse struct {
}

func (c *Client) DeleteRole(ctx context.Context, req DeleteRoleRequest) (*DeleteRoleResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(deleteRolePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp DeleteRoleResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type AttachFunctionalitiesToRoleRequest struct {
	RoleID                int64                     `json:"RoleID"`
	ScopedFunctionalities []RoleScopedFunctionality `json:"ScopedFunctionalities,omitempty"`
}

type AttachFunctionalitiesToRoleResponse struct {
}

func (c *Client) AttachFunctionalitiesToRole(ctx context.Context, req AttachFunctionalitiesToRoleRequest) (*AttachFunctionalitiesToRoleResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(attachFunctionalitiesToRolePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp AttachFunctionalitiesToRoleResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type DetachFunctionalitiesFromRoleRequest struct {
	RoleID                int64                     `json:"RoleID"`
	ScopedFunctionalities []RoleScopedFunctionality `json:"ScopedFunctionalities,omitempty"`
}

type DetachFunctionalitiesFromRoleResponse struct {
}

func (c *Client) DetachFunctionalitiesFromRole(ctx context.Context, req DetachFunctionalitiesFromRoleRequest) (*DetachFunctionalitiesFromRoleResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(detachFunctionalitiesFromRolePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New("Request failed.")
	}

	tflog.Debug(ctx, "Request info", "Status code", resp.StatusCode(), "body", resp.String())

	var jsonResp DetachFunctionalitiesFromRoleResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Error: %s \n Received: %s", err, resp.String()))
	}

	return &jsonResp, nil
}
