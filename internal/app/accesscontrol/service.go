package accesscontrol

import (
	"fmt"

	coreAccessControl "github.com/tacenva/tacpass-core/accesscontrol"
	"github.com/tacenva/tacpass-core/entity"

	"github.com/tacenva/tacpass-tui/internal/api"
)

type Service struct {
	api           *api.Client
	address       string
	acCoreService *coreAccessControl.Service
	isRemote      bool
}

func NewService(
	apiClient *api.Client,
	address string,
	acCoreService *coreAccessControl.Service,
	isRemote bool,
) *Service {
	return &Service{
		api:           apiClient,
		address:       address,
		acCoreService: acCoreService,
		isRemote:      isRemote,
	}
}

func (s *Service) List() ([]entity.Permission, error) {
	if !s.isRemote {
		if s.acCoreService == nil {
			return nil, fmt.Errorf(
				"access control core service is not initialized",
			)
		}

		return s.acCoreService.List()
	}

	response, err := s.api.ListAccessControls(
		s.address,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list access controls: %w",
			err,
		)
	}

	return response, nil
}

func (s *Service) Get(
	id string,
) (*entity.Permission, error) {
	if !s.isRemote {
		if s.acCoreService == nil {
			return nil, fmt.Errorf(
				"access control core service is not initialized",
			)
		}

		response, err := s.acCoreService.Get(id)
		if err != nil {
			return nil, fmt.Errorf(
				"get access control: %w",
				err,
			)
		}

		return response, nil
	}

	response, err := s.api.GetAccessControl(
		s.address,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get access control: %w",
			err,
		)
	}

	return response, nil
}

func (s *Service) Create(
	privilege entity.Privilege,
) (*api.AccessControlCreateResponse, error) {
	if !s.isRemote {
		if s.acCoreService == nil {
			return nil, fmt.Errorf(
				"access control core service is not initialized",
			)
		}

		permission, keyPair, err :=
			s.acCoreService.Create(privilege)

		if err != nil {
			return nil, fmt.Errorf(
				"create access control: %w",
				err,
			)
		}

		return &api.AccessControlCreateResponse{
			Permission: permission,
			KeyPair:    keyPair,
		}, nil
	}

	response, err := s.api.CreateAccessControl(
		s.address,
		privilege,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create access control: %w",
			err,
		)
	}

	return response, nil
}

func (s *Service) ChangePrivilege(
	id string,
	privilege entity.Privilege,
) error {
	if !s.isRemote {
		if s.acCoreService == nil {
			return fmt.Errorf(
				"access control core service is not initialized",
			)
		}

		if err := s.acCoreService.ChangePrivilege(
			id,
			privilege,
		); err != nil {
			return fmt.Errorf(
				"change access control privilege: %w",
				err,
			)
		}

		return nil
	}

	if err := s.api.ChangeAccessControlPrivilege(
		s.address,
		id,
		privilege,
	); err != nil {
		return fmt.Errorf(
			"change access control privilege: %w",
			err,
		)
	}

	return nil
}

func (s *Service) Revoke(
	id string,
) error {
	if !s.isRemote {
		if s.acCoreService == nil {
			return fmt.Errorf(
				"access control core service is not initialized",
			)
		}

		if err := s.acCoreService.Revoke(id); err != nil {
			return fmt.Errorf(
				"revoke access control: %w",
				err,
			)
		}

		return nil
	}

	if err := s.api.RevokeAccessControl(
		s.address,
		id,
	); err != nil {
		return fmt.Errorf(
			"revoke access control: %w",
			err,
		)
	}

	return nil
}

func (s *Service) ListUsers(
	permissionID string,
) ([]entity.User, error) {
	if !s.isRemote {
		if s.acCoreService == nil {
			return nil, fmt.Errorf(
				"access control core service is not initialized",
			)
		}

		response, err := s.acCoreService.UserList(
			permissionID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"list access control users: %w",
				err,
			)
		}

		return response, nil
	}

	response, err := s.api.ListAccessControlUsers(
		s.address,
		permissionID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list access control users: %w",
			err,
		)
	}

	return response, nil
}

func (s *Service) ApproveUser(
	userID string,
) (*entity.User, error) {
	if !s.isRemote {
		if s.acCoreService == nil {
			return nil, fmt.Errorf(
				"access control core service is not initialized",
			)
		}

		response, err := s.acCoreService.ApproveUser(
			userID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"approve user: %w",
				err,
			)
		}

		return response, nil
	}

	response, err := s.api.ApproveAccessControlUser(
		s.address,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"approve user: %w",
			err,
		)
	}

	return response, nil
}

func (s *Service) RevokeUser(
	userID string,
) (*entity.User, error) {
	if !s.isRemote {
		if s.acCoreService == nil {
			return nil, fmt.Errorf(
				"access control core service is not initialized",
			)
		}

		response, err := s.acCoreService.RevokeUser(
			userID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"revoke user: %w",
				err,
			)
		}

		return response, nil
	}

	response, err := s.api.RevokeAccessControlUser(
		s.address,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"revoke user: %w",
			err,
		)
	}

	return response, nil
}
