package eva

import (
	"fmt"
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

func (c *Client) Login(credentials LoginCredentials) error {

	resp, err := c.Client.R().
		SetBody(credentials).
		Post(path)

	if err != nil {
		return err
	}

	fmt.Println("  Body       :\n", resp)

	resp.Request.SetAuthToken("123")

	return nil
}
