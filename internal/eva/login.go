package eva

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	loginPath = "/api/core/Login"
)

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AuthenticationToken string `json:"AuthenticationToken"`
}

func (c *Client) Login(ctx context.Context, req LoginCredentials) error {

	resp, err := c.Client.R().
		SetBody(req).
		Post(loginPath)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return err
	}

	if resp.StatusCode() != 200 {
		tflog.Info(ctx, "Request failed.", "Status code", resp.StatusCode(), "body", resp.String())

		return errors.New("Login failed.")
	}

	var jsonResp LoginResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		return errors.New(fmt.Sprintf("Response could not be parsed. Error received: %s \nResponse received: %s", err, resp.String()))
	}

	c.Client.SetHeader("authorization", jsonResp.AuthenticationToken)

	return nil
}
