package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	createMessageTemplatePath = "/api/core/management/CreateMessageTemplate"
	updateMessageTemplatePath = "/api/core/management/UpdateMessageTemplate"
	deleteMessageTemplatePath = "/api/core/management/DeleteMessageTemplate"
)

// WaitForNetworkIdle: boolean;
// WaitForJS: boolean;
// Size: PaperProperties.PaperSize;
// Format?: PaperFormats;
// Orientation?: PaperOrientations;
// Margin: PaperProperties.PaperMargin;
// ThermalPrinterTemplateType?: ThermalPrinterTemplateTypes;

// export interface PaperMargin {
// 	Top?: number;
// 	Left?: number;
// 	Bottom?: number;
// 	Right?: number;
//   }

//   export interface PaperSize {
// 	Width: string;
// 	Height: string;
// 	DeviceScaleFactor?: number;
//   }

type PaperMargin struct {
	Top    int64 `json:"Top,omitempty"`
	Left   int64 `json:"Left,omitempty"`
	Bottom int64 `json:"Bottom,omitempty"`
	Right  int64 `json:"Right,omitempty"`
}

type PaperSize struct {
	Width             string `json:"Width"`
	Height            string `json:"Height"`
	DeviceScaleFactor int64  `json:"Bottom,omitempty"`
}

type PaperProperties struct {
	WaitForNetworkIdle         bool        `json:"WaitForNetworkIdle"`
	WaitForJS                  bool        `json:"WaitForJS"`
	Size                       PaperSize   `json:"Size"`
	Format                     int64       `json:"Format,omitempty"`
	Orientation                int64       `json:"Orientation,omitempty"`
	Margin                     PaperMargin `json:"PaperMargin"`
	ThermalPrinterTemplateType int64       `json:"ThermalPrinterTemplateType,omitempty"`
}

type CreateMessageTemplateRequest struct {
	Name               string           `json:"Name"`
	OrganizationUnitID int64            `json:"OrganizationUnitID,omitempty"`
	LanguageID         string           `json:"LanguageID,omitempty"`
	CountryID          string           `json:"CountryID,omitempty"`
	Header             string           `json:"Header,omitempty"`
	Template           string           `json:"Template"`
	Footer             string           `json:"Footer,omitempty"`
	Helper             string           `json:"Helper,omitempty"`
	Type               int64            `json:"Type"`
	Layout             string           `json:"Layout,omitempty"`
	Destination        int64            `json:"Destination"`
	PaperProperties    *PaperProperties `json:"PaperProperties,omitempty"` //omitempty doesn't work for struct unless it is a pointer
	IsDisabled         bool             `json:"IsDisable,omitempty"`
}

type CreateMessageTemplateResponse struct {
	ID int64
}

func (c *Client) createMessageTemplate(ctx context.Context, req CreateMessageTemplateRequest) (*CreateMessageTemplateResponse, error) {
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
