package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	createMessageTemplatePath  = "/api/core/management/CreateMessageTemplate"
	getMessageTemplateByIDPath = "/api/core/management/GetMessageTemplateByID"
	updateMessageTemplatePath  = "/api/core/management/UpdateMessageTemplate"
	deleteMessageTemplatePath  = "/api/core/management/DeleteMessageTemplate"
)

type PaperMargin struct {
	Top    int64 `json:"Top"`
	Left   int64 `json:"Left"`
	Bottom int64 `json:"Bottom"`
	Right  int64 `json:"Right"`
}

type PaperSize struct {
	Width             string  `json:"Width"`
	Height            string  `json:"Height"`
	DeviceScaleFactor float64 `json:"DeviceScaleFactor"`
}

type PaperProperties struct {
	WaitForNetworkIdle         bool        `json:"WaitForNetworkIdle"`
	WaitForJS                  bool        `json:"WaitForJS"`
	Size                       PaperSize   `json:"Size"`
	Format                     int64       `json:"Format"`
	Orientation                int64       `json:"Orientation"`
	Margin                     PaperMargin `json:"PaperMargin"`
	ThermalPrinterTemplateType int64       `json:"ThermalPrinterTemplateType"`
}

type CreateMessageTemplateRequest struct {
	Name               string           `json:"Name"`
	OrganizationUnitID int64            `json:"OrganizationUnitID,omitempty"`
	LanguageID         string           `json:"LanguageID,omitempty"`
	CountryID          string           `json:"CountryID,omitempty"`
	Header             string           `json:"Header,omitempty"`
	Template           string           `json:"Template"`
	Footer             string           `json:"Footer,omitempty"`
	Helpers            string           `json:"Helpers,omitempty"`
	Type               int64            `json:"Type"`
	Layout             string           `json:"Layout,omitempty"`
	Destination        int64            `json:"Destination"`
	PaperProperties    *PaperProperties `json:"PaperProperties,omitempty"` //omitempty doesn't work for struct unless it is a pointer
	IsDisabled         bool             `json:"IsDisable,omitempty"`
}

type CreateMessageTemplateResponse struct {
	ID int64
}

func (c *Client) CreateMessageTemplate(ctx context.Context, req CreateMessageTemplateRequest) (*CreateMessageTemplateResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(createMessageTemplatePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp CreateMessageTemplateResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type GetMessageTemplateByIDRequest struct {
	ID int64
}

type GetMessageTemplateByIDResponse struct {
	ID                 int64            `json:"ID"`
	Name               string           `json:"Name"`
	OrganizationUnitID int64            `json:"OrganizationUnitID,omitempty"`
	LanguageID         string           `json:"LanguageID"`
	CountryID          string           `json:"CountryID"`
	Header             string           `json:"Header"`
	Template           string           `json:"Template"`
	Footer             string           `json:"Footer"`
	Helpers            string           `json:"Helpers"`
	Type               int64            `json:"Type"`
	Layout             string           `json:"Layout"`
	Destination        int64            `json:"Destination"`
	PaperProperties    *PaperProperties `json:"PaperProperties"`
	IsDisabled         bool             `json:"IsDisable,omitempty"`
}

func (c *Client) GetMessageTemplateByID(ctx context.Context, req GetMessageTemplateByIDRequest) (*GetMessageTemplateByIDResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(getMessageTemplateByIDPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp GetMessageTemplateByIDResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", err.Error()))
	}

	return &jsonResp, nil
}

type UpdateMessageTemplateRequest struct {
	ID                 int64            `json:"ID"`
	Name               string           `json:"Name"`
	OrganizationUnitID int64            `json:"OrganizationUnitID,omitempty"`
	LanguageID         string           `json:"LanguageID,omitempty"`
	CountryID          string           `json:"CountryID,omitempty"`
	Header             string           `json:"Header,omitempty"`
	Template           string           `json:"Template"`
	Footer             string           `json:"Footer,omitempty"`
	Helpers            string           `json:"Helpers,omitempty"`
	Layout             string           `json:"Layout,omitempty"`
	Destination        int64            `json:"Destination"`
	PaperProperties    *PaperProperties `json:"PaperProperties,omitempty"` //omitempty doesn't work for struct unless it is a pointer
	IsDisabled         bool             `json:"IsDisable,omitempty"`
}

type UpdateMessageTemplateResponse struct {
}

func (c *Client) UpdateMessageTemplate(ctx context.Context, req UpdateMessageTemplateRequest) (*UpdateMessageTemplateResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(updateMessageTemplatePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp UpdateMessageTemplateResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}

type DeleteMessageTemplateRequesst struct {
	ID int64
}

type DeleteMessageTemplateResponse struct {
}

func (c *Client) DeleteMessageTemplate(ctx context.Context, req DeleteMessageTemplateRequesst) (*DeleteMessageTemplateResponse, error) {
	resp, err := c.restClient.R().
		SetBody(req).
		Post(deleteMessageTemplatePath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return nil, err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed", "Status code", resp.StatusCode(), "body", resp.String())

		return nil, errors.New(fmt.Sprintf("Request failed with error: %s", resp.String()))
	}

	var jsonResp DeleteMessageTemplateResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {

		return nil, errors.New(fmt.Sprintf("Response could not be parsed. Received: %s", resp.String()))
	}

	return &jsonResp, nil
}
