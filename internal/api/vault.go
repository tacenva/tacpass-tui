package api

import "fmt"

func (c *Client) ListVaults(
	address string,
	result any,
) error {
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
