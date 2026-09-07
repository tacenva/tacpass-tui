package api

import (
	"fmt"
	"os"
	"strings"
)

func (c *Client) ListVaults(
	address string,
	result any,
) error {
	debugFile, err := os.OpenFile(
		"debug.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err == nil {
		defer debugFile.Close()

		fmt.Fprintf(
			debugFile,
			"ListVaults: GET %s/vault\nToken: %s\n",
			strings.TrimRight(address, "/"),
			c.Token,
		)
	}

	return c.Get(
		address,
		"/vault",
		result,
	)
}

func (c *Client) GetVault(
	address string,
	id string,
	result any,
) error {
	path := fmt.Sprintf(
		"/vault/%s",
		id,
	)

	return c.Get(
		address,
		path,
		result,
	)
}

func (c *Client) CreateVault(
	address string,
	body any,
	result any,
) error {
	return c.Post(
		address,
		"/vault",
		body,
		result,
	)
}

func (c *Client) UpdateVault(
	address string,
	id string,
	body any,
	result any,
) error {
	path := fmt.Sprintf(
		"/vault/%s",
		id,
	)

	return c.Patch(
		address,
		path,
		body,
		result,
	)
}

func (c *Client) DeleteVault(
	address string,
	id string,
) error {
	path := fmt.Sprintf(
		"/vault/%s",
		id,
	)

	return c.Delete(
		address,
		path,
		nil,
	)
}
