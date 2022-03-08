package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	attachFunctionalitiesToRolePath   = "/api/core/management/AttachFunctionalitiesToRole"
	detachFunctionalitiesFromRolePath = "/api/core/management/DetachFunctionalitiesFromRole"
)

type AttachFunctionalitiesToRoleRequest struct {
	RoleID                int64                 `json:"RoleID"`
	ScopedFunctionalities []ScopedFunctionality `json:"ScopedFunctionalities,omitempty"`
}

type ScopedFunctionality struct {
	Functionality     string
	Scope             int64
	RequiresElevation bool
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
	RoleID                int64                 `json:"RoleID"`
	ScopedFunctionalities []ScopedFunctionality `json:"ScopedFunctionalities,omitempty"`
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
