package api

import (
	"fmt"

	coreEntity "github.com/tacenva/tacpass-core/entity"
)

type AccessControlCreateRequest struct {
	Privilege coreEntity.Privilege `json:"privilege"`
}

type AccessControlPrivilegeRequest struct {
	Privilege coreEntity.Privilege `json:"privilege"`
}

func (c *Client) ListAccessControl(
	address string,
) ([]coreEntity.Permission, error) {
	var response []coreEntity.Permission

	if err := c.Get(
		address,
		"/access-control",
		&response,
	); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetAccessControlByPublicKey(
	address string,
	publicKey string,
) (*coreEntity.Permission, error) {
	var response coreEntity.Permission

	path := fmt.Sprintf(
		"/access-control/public-key/%s",
		publicKey,
	)

	if err := c.Get(
		address,
		path,
		&response,
	); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) CreateAccessControl(
	address string,
	privilege coreEntity.Privilege,
) (*coreEntity.Permission, error) {
	var response coreEntity.Permission

	if err := c.Post(
		address,
		"/access-control",
		AccessControlCreateRequest{
			Privilege: privilege,
		},
		&response,
	); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) ChangeAccessControlPrivilege(
	address string,
	id string,
	privilege coreEntity.Privilege,
) error {
	path := fmt.Sprintf(
		"/access-control/%s/privilege",
		id,
	)

	return c.Patch(
		address,
		path,
		AccessControlPrivilegeRequest{
			Privilege: privilege,
		},
		nil,
	)
}

func (c *Client) RevokeAccessControl(
	address string,
	id string,
) error {
	path := fmt.Sprintf(
		"/access-control/%s",
		id,
	)

	return c.Delete(
		address,
		path,
		nil,
	)
}

func (c *Client) ListAccessControlUsers(
	address string,
	permissionID string,
) ([]coreEntity.User, error) {
	var response []coreEntity.User

	path := fmt.Sprintf(
		"/access-control/%s/users",
		permissionID,
	)

	if err := c.Get(
		address,
		path,
		&response,
	); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ApproveUser(
	address string,
	userID string,
) (*coreEntity.User, error) {
	var response coreEntity.User

	path := fmt.Sprintf(
		"/access-control/users/%s/approve",
		userID,
	)

	if err := c.Post(
		address,
		path,
		nil,
		&response,
	); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) RevokeUser(
	address string,
	userID string,
) (*coreEntity.User, error) {
	var response coreEntity.User

	path := fmt.Sprintf(
		"/access-control/users/%s/revoke",
		userID,
	)

	if err := c.Post(
		address,
		path,
		nil,
		&response,
	); err != nil {
		return nil, err
	}

	return &response, nil
}
