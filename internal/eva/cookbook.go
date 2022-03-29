package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	getAccountingRecipePath    = "/api/cookbook/GetAccountingRecipe"
	createAccountingRecipePath = "/api/cookbook/CreateAccountingRecipe"
	updateAccountingRecipePath = "/api/cookbook/UpdateAccountingRecipe"
	deleteAccountingRecipePath = "/api/cookbook/DeleteAccountingRecipe"
)

type CreateAccountingRecipeRequest struct {
	Name     string `json:"Name"`
	Recipe   string `json:"Recipe"`
	IsActive bool   `json:"IsActive,omitempty"`
}

type CreateAccountingRecipeResponse struct {
	ID        int64 `json:"ID"`
	HasErrors bool  `json:"HasErrors"`
}

func (c *Client) CreateAccountingRecipe(ctx context.Context, req CreateAccountingRecipeRequest) (*CreateAccountingRecipeResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(createAccountingRecipePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp CreateAccountingRecipeResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil || jsonResp.HasErrors {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type GetAccountingRecipeRequest struct {
	ID int64 `json:"ID"`
}

type recipe struct {
	ID       int64  `json:"ID"`
	Name     string `json:"Name"`
	Recipe   string `json:"Recipe"`
	IsActive bool   `json:"IsActive"`
}

type GetAccountingRecipeResponse struct {
	Recipe recipe `json:"Recipe"`
}

func (c *Client) GetAccountingRecipe(ctx context.Context, req GetAccountingRecipeRequest) (*GetAccountingRecipeResponse, error) {

	resp, err := c.restClient.R().
		SetBody(req).
		Post(getAccountingRecipePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New("Request failed.")
	}

	tflog.Debug(ctx, "Request info", "Status code", resp.StatusCode(), "body", resp.String())

	var jsonResp GetAccountingRecipeResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Error: %s \n Received: %s", err, resp.String()))
	}

	return &jsonResp, nil
}

type UpdateAccountingRecipeRequest struct {
	ID       int64  `json:"ID"`
	Name     string `json:"Name,omitempty"`
	Recipe   string `json:"Recipe,omitempty"`
	IsActive bool   `json:"IsActive,omitempty"`
}

type UpdateAccountingRecipeResponse struct {
	HasErrors bool `json:"HasErrors"`
}

func (c *Client) UpdateAccountingRecipe(ctx context.Context, req UpdateAccountingRecipeRequest) (*UpdateAccountingRecipeResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(updateAccountingRecipePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp UpdateAccountingRecipeResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil || jsonResp.HasErrors {
		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type DeleteAccountingRecipeRequest struct {
	ID int64 `json:"ID"`
}

func (c *Client) DeleteAccountingRecipe(ctx context.Context, req DeleteAccountingRecipeRequest) (*EmptyResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(deleteAccountingRecipePath)

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
