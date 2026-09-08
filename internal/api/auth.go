package api

import (
	"fmt"
	"os"
)

type EnrollRequest struct {
	Hostname  string `json:"hostname"`
	PublicKey string `json:"public_key"`
}

type EnrollResponse struct {
	Status    string `json:"status"`
	AuthToken string `json:"auth_token"`
}

func (c *Client) Enroll(
	address string,
	request EnrollRequest,
) (*EnrollResponse, error) {
	debugFile, err := os.OpenFile(
		"debug.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err == nil {
		defer debugFile.Close()

		fmt.Fprintf(
			debugFile,
			"\n=== Enroll ===\n"+
				"Address: %s\n"+
				"Hostname: %s\n"+
				"PublicKey: %s\n",
			address,
			request.Hostname,
			request.PublicKey,
		)
	}

	var response EnrollResponse

	if err := c.PostPublic(
		address,
		"/auth/enroll",
		request,
		&response,
	); err != nil {
		if debugFile != nil {
			fmt.Fprintf(
				debugFile,
				"Error: %v\n",
				err,
			)
		}

		return nil, err
	}

	if debugFile != nil {
		fmt.Fprintf(
			debugFile,
			"Status: %s\n"+
				"AuthToken exists: %t\n",
			response.Status,
			response.AuthToken != "",
		)
	}

	return &response, nil
}

type LoginRequest struct {
	PublicKey string `json:"public_key"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (c *Client) Login(
	address string,
	request LoginRequest,
) (*LoginResponse, error) {
	var response LoginResponse

	if err := c.PostPublic(
		address,
		"/auth/login",
		request,
		&response,
	); err != nil {
		return nil, err
	}

	c.SetToken(response.Token)

	return &response, nil
}

func (c *Client) Logout() {
	c.ClearToken()
}
