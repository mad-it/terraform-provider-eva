package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	getOpenIDProviderPath    = "/api/authentication/openid/GetOpenIDProviderByID"
	createOpenIDProviderPath = "/api/authentication/openid/CreateOpenIDProvider"
	updateOpenIDProviderPath = "/api/authentication/openid/UpdateOpenIDProvider"
	deleteOpenIDProviderPath = "/api/authentication/openid/DeleteOpenIDProvider"
)

type CreateOpenIDProviderRequest struct {
	BaseUrl           string `json:"BaseUrl"`
	ClientID          string `json:"ClientID"`
	CreateUsers       bool   `json:"CreateUsers"`
	EmailAddressClaim string `json:"EmailAddressClaim,omitempty"`
	Enabled           bool   `json:"Enabled"`
	FirstNameClaim    string `json:"FirstNameClaim,omitempty"`
	LastNameClaim     string `json:"LastNameClaim,omitempty"`
	Name              string `json:"Name,omitempty"`
	NicknameClaim     string `json:"NicknameClaim,omitempty"`
	UserType          int64  `json:"UserType"`
}

type CreateOpenIDProviderResponse struct {
	ID int64 `json:"ID"`
}

func (c *Client) CreateOpenIDProvider(ctx context.Context, req CreateOpenIDProviderRequest) (*CreateOpenIDProviderResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(createOpenIDProviderPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp CreateOpenIDProviderResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type GetOpenIDProviderRequest struct {
	ID int64 `json:"ID"`
}

type GetOpenIDProviderResponse struct {
	ID                int64  `json:"ID"`
	BaseUrl           string `json:"BaseUrl"`
	ClientID          string `json:"ClientID"`
	CreateUsers       bool   `json:"CreateUsers"`
	EmailAddressClaim string `json:"EmailAddressClaim,omitempty"`
	Enabled           bool   `json:"Enabled"`
	FirstNameClaim    string `json:"FirstNameClaim,omitempty"`
	LastNameClaim     string `json:"LastNameClaim,omitempty"`
	Name              string `json:"Name,omitempty"`
	NicknameClaim     string `json:"NicknameClaim,omitempty"`
	UserType          int64  `json:"UserType"`
}

func (c *Client) GetOpenIDProvider(ctx context.Context, req GetOpenIDProviderRequest) (*GetOpenIDProviderResponse, error) {

	resp, err := c.restClient.R().
		SetBody(req).
		Post(getOpenIDProviderPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New("Request failed.")
	}

	tflog.Debug(ctx, "Request info", "Status code", resp.StatusCode(), "body", resp.String())

	var jsonResp GetOpenIDProviderResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Error: %s \n Received: %s", err, resp.String()))
	}

	return &jsonResp, nil
}

type UpdateOpenIDProviderRequest struct {
	ID                int64  `json:"ID"`
	BaseUrl           string `json:"BaseUrl,omitempty"`
	ClientID          string `json:"ClientID,omitempty"`
	CreateUsers       bool   `json:"CreateUsers,omitempty"`
	EmailAddressClaim string `json:"EmailAddressClaim,omitempty"`
	Enabled           bool   `json:"Enabled,omitempty"`
	FirstNameClaim    string `json:"FirstNameClaim,omitempty"`
	LastNameClaim     string `json:"LastNameClaim,omitempty"`
	Name              string `json:"Name,omitempty"`
	NicknameClaim     string `json:"NicknameClaim,omitempty"`
	UserType          int64  `json:"UserType"`
}

type UpdateOpenIDProviderResponse struct{}

func (c *Client) UpdateOpenIDProvider(ctx context.Context, req UpdateOpenIDProviderRequest) (*UpdateOpenIDProviderResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(updateOpenIDProviderPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp UpdateOpenIDProviderResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type DeleteOpenIDProviderRequest struct {
	ID int64 `json:"ID"`
}

type DeleteOpenIDProviderResponse struct{}

func (c *Client) DeleteOpenIDProvider(ctx context.Context, req DeleteOpenIDProviderRequest) (*DeleteOpenIDProviderResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(deleteOpenIDProviderPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp DeleteOpenIDProviderResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}
