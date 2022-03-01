package eva

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AuthenticationToken string `json:"AuthenticationToken"`
}

func (c *Client) Login(ctx context.Context, req LoginCredentials) error {
	const (
		path = "/api/core/Login"
	)

	resp, err := c.Client.R().
		SetBody(req).
		Post(path)

	if err != nil {
		tflog.Error(ctx, "An network error ocurred.", err)

		return err
	}

	if resp.StatusCode() != 200 {
		tflog.Debug(ctx, "Request failed.", "Status code", resp.StatusCode(), "body", resp.String())

		return errors.New("Login failed.")
	}

	var jsonResp LoginResponse
	if err := json.Unmarshal([]byte(resp.Body()), &jsonResp); err != nil {
		panic(err)
	}

	c.Client.SetHeader("authorization", jsonResp.AuthenticationToken)

	return nil
}
