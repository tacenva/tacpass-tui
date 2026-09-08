package accesscontrol

import (
	"fmt"

	coreAC "github.com/tacenva/tacpass-core/accesscontrol"
	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-core/util/keyring"

	"github.com/tacenva/tacpass-tui/internal/app"
)

type Service struct {
	appDeps *app.Deps
	context *app.Context

	coreACService *coreAC.Service
}

func NewService(
	appDeps *app.Deps,
	context *app.Context,
	coreACService *coreAC.Service,
) *Service {
	return &Service{
		appDeps:       appDeps,
		context:       context,
		coreACService: coreACService,
	}
}

func (s *Service) List() ([]entity.Permission, error) {
	if s.context.IsRemote {
		if s.context.SelectedSoT == nil {
			return nil, fmt.Errorf("selected source of truth is nil")
		}

		return s.appDeps.Client.ListAccessControls(
			s.context.SelectedSoT.Address,
		)
	}

	return s.coreACService.List()
}

func (s *Service) Get(
	id string,
) (*entity.Permission, error) {
	if s.context.IsRemote {
		if s.context.SelectedSoT == nil {
			return nil, fmt.Errorf("selected source of truth is nil")
		}

		return s.appDeps.Client.GetAccessControl(
			s.context.SelectedSoT.Address,
			id,
		)
	}

	return s.coreACService.Get(id)
}

func (s *Service) Create(
	name string,
	privilege entity.Privilege,
) (*entity.Permission, *keyring.KeyPair, error) {
	if s.context.IsRemote {
		if s.context.SelectedSoT == nil {
			return nil, nil, fmt.Errorf(
				"selected source of truth is nil",
			)
		}

		result, err := s.appDeps.Client.CreateAccessControl(
			s.context.SelectedSoT.Address,
			name,
			privilege,
		)
		if err != nil {
			return nil, nil, err
		}

		return result.Permission, result.KeyPair, nil
	}

	return s.coreACService.Create(
		name,
		privilege,
	)
}

func (s *Service) ChangeName(
	id string,
	name string,
) error {
	if s.context.IsRemote {
		if s.context.SelectedSoT == nil {
			return fmt.Errorf("selected source of truth is nil")
		}

		return s.appDeps.Client.ChangeAccessControlName(
			s.context.SelectedSoT.Address,
			id,
			name,
		)
	}

	return s.coreACService.ChangeName(id, name)
}

func (s *Service) ChangePrivilege(
	id string,
	privilege entity.Privilege,
) error {
	if s.context.IsRemote {
		if s.context.SelectedSoT == nil {
			return fmt.Errorf(
				"selected source of truth is nil",
			)
		}

		return s.appDeps.Client.ChangeAccessControlPrivilege(
			s.context.SelectedSoT.Address,
			id,
			privilege,
		)
	}

	return s.coreACService.ChangePrivilege(
		id,
		privilege,
	)
}

func (s *Service) Revoke(
	id string,
) error {
	if s.context.IsRemote {
		if s.context.SelectedSoT == nil {
			return fmt.Errorf(
				"selected source of truth is nil",
			)
		}

		return s.appDeps.Client.RevokeAccessControl(
			s.context.SelectedSoT.Address,
			id,
		)
	}

	return s.coreACService.Revoke(id)
}

func (s *Service) UserList(
	permissionID string,
) ([]entity.User, error) {
	if s.context.IsRemote {
		if s.context.SelectedSoT == nil {
			return nil, fmt.Errorf(
				"selected source of truth is nil",
			)
		}

		return s.appDeps.Client.ListAccessControlUsers(
			s.context.SelectedSoT.Address,
			permissionID,
		)
	}

	return s.coreACService.UserList(
		permissionID,
	)
}

func (s *Service) ApproveUser(
	userID string,
) (*entity.User, error) {
	if s.context.IsRemote {
		if s.context.SelectedSoT == nil {
			return nil, fmt.Errorf(
				"selected source of truth is nil",
			)
		}

		return s.appDeps.Client.ApproveAccessControlUser(
			s.context.SelectedSoT.Address,
			userID,
		)
	}

	return s.coreACService.ApproveUser(
		userID,
	)
}

func (s *Service) RevokeUser(
	userID string,
) (*entity.User, error) {
	if s.context.IsRemote {
		if s.context.SelectedSoT == nil {
			return nil, fmt.Errorf(
				"selected source of truth is nil",
			)
		}

		return s.appDeps.Client.RevokeAccessControlUser(
			s.context.SelectedSoT.Address,
			userID,
		)
	}

	return s.coreACService.RevokeUser(
		userID,
	)
}
