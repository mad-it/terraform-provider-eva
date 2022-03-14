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
	getUserRolePath                   = "/api/core/management/GetUserRoles"
	setUserRolePath                   = "/api/core/management/SetUserRoles"
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
	Name                  string              `json:"Name"`
	UserType              int64               `json:"UserType,omitempty"`
	Code                  string              `json:"Code,omitempty"`
	ScopedFunctionalities []RoleFunctionality `json:"ScopedFunctionalities"`
}

type RoleFunctionality struct {
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

func (c *Client) UpdateRole(ctx context.Context, req UpdateRoleRequest) (*EmptyResponse, error) {
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

	var jsonResp EmptyResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type DeleteRoleRequest struct {
	ID int64 `json:"ID"`
}

func (c *Client) DeleteRole(ctx context.Context, req DeleteRoleRequest) (*EmptyResponse, error) {
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

	var jsonResp EmptyResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type AttachFunctionalitiesToRoleRequest struct {
	RoleID                int64               `json:"RoleID"`
	ScopedFunctionalities []RoleFunctionality `json:"ScopedFunctionalities,omitempty"`
}

func (c *Client) AttachFunctionalitiesToRole(ctx context.Context, req AttachFunctionalitiesToRoleRequest) (*EmptyResponse, error) {
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

	var jsonResp EmptyResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type DetachFunctionalitiesFromRoleRequest struct {
	RoleID                int64               `json:"RoleID"`
	ScopedFunctionalities []RoleFunctionality `json:"ScopedFunctionalities,omitempty"`
}

func (c *Client) DetachFunctionalitiesFromRole(ctx context.Context, req DetachFunctionalitiesFromRoleRequest) (*EmptyResponse, error) {
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

	var jsonResp EmptyResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Error: %s \n Received: %s", err, resp.String()))
	}

	return &jsonResp, nil
}

type GetUserRoleRequest struct {
	UserId int64 `json:"UserID"`
}

type UserRole struct {
	RoleID             int64 `json:"RoleID"`
	OrganizationUnitID int64 `json:"OrganizationUnitID"`
	UserType           int64 `json:"UserType"`
}

type GetUserRoleResponse struct {
	Roles []UserRole `json:"Roles"`
}

func (c *Client) GetUserRole(ctx context.Context, req GetUserRoleRequest) (*GetUserRoleResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(getUserRolePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New("Request failed.")
	}

	tflog.Debug(ctx, "Request info", "Status code", resp.StatusCode(), "body", resp.String())

	var jsonResp GetUserRoleResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Error: %s \n Received: %s", err, resp.String()))
	}

	return &jsonResp, nil
}

type RoleOrganizationUnitSet struct {
	RoleID             int64 `json:"RoleID"`
	OrganizationUnitID int64 `json:"OrganizationUnitID"`
	UserType           int64 `json:"UserType"`
}

type SetUserRoleRequest struct {
	UserId int64                     `json:"UserID"`
	Roles  []RoleOrganizationUnitSet `json:"Roles"`
}

func (c *Client) SetUserRole(ctx context.Context, req SetUserRoleRequest) (*EmptyResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(setUserRolePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New("Request failed.")
	}

	tflog.Debug(ctx, "Request info", "Status code", resp.StatusCode(), "body", resp.String())

	var jsonResp EmptyResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Error: %s \n Received: %s", err, resp.String()))
	}

	return &jsonResp, nil
}
