package eva

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	path = "/api/core/Login"
)

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AuthenticationToken string `json:"AuthenticationToken"`
}

func (c *Client) Login(ctx context.Context, credentials LoginCredentials) error {

	resp, err := c.Client.R().
		SetBody(credentials).
		Post(path)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return err
	}

	if resp.StatusCode() != 200 {
		tflog.Trace(ctx, "Login failed.", "Status code", resp.StatusCode(), "body", resp.String())

		return errors.New("Login failed.")
	}

	var loginResponse LoginResponse
	if err := json.Unmarshal([]byte(resp.Body()), &loginResponse); err != nil {
		panic(err)
	}

	resp.Request.SetAuthToken(loginResponse.AuthenticationToken)

	return nil
}
