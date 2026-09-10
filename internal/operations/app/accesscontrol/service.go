package accesscontrol

import (
	"fmt"

	coreAC "github.com/tacenva/tacpass-core/accesscontrol"
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-core/util/keyring"
	coreVault "github.com/tacenva/tacpass-core/vault"

	"github.com/tacenva/tacpass-tui/internal/operations/app"
)

type Service struct {
	appDeps *app.Deps
	context *app.Context

	coreACService    *coreAC.Service
	coreVaultService *coreVault.Service
	authService      *auth.Service
}

func NewService(
	appDeps *app.Deps,
	context *app.Context,
	coreACService *coreAC.Service,
	authService *auth.Service,
) *Service {
	return &Service{
		appDeps:       appDeps,
		context:       context,
		coreACService: coreACService,
		authService:   authService,
	}
}

func (s *Service) authUser() (*entity.User, error) {
	if s.context.SelectedSoT == nil {
		return nil, fmt.Errorf("selected source of truth is nil")
	}

	return s.authService.GetUserData(
		s.context.SelectedSoT.AuthToken,
	)
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

	authUser, err := s.authUser()
	if err != nil {
		return nil, err
	}

	return s.coreACService.List(authUser)
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

	authUser, err := s.authUser()
	if err != nil {
		return nil, err
	}

	return s.coreACService.Get(
		authUser,
		id,
	)
}

func (s *Service) Create(
	name string,
	privilege entity.Privilege,
) ([]entity.VaultAccess, *entity.Permission, *keyring.KeyPair, error) {
	if s.context.IsRemote {
		if s.context.SelectedSoT == nil {
			return nil, nil, nil, fmt.Errorf(
				"selected source of truth is nil",
			)
		}

		result, err := s.appDeps.Client.CreateAccessControl(
			s.context.SelectedSoT.Address,
			name,
			privilege,
		)
		if err != nil {
			return nil, nil, nil, err
		}

		return result.VaultAccessList, result.Permission, result.KeyPair, nil
	}

	authUser, err := s.authUser()
	if err != nil {
		return nil, nil, nil, err
	}

	return s.coreACService.Create(
		authUser,
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

	authUser, err := s.authUser()
	if err != nil {
		return err
	}

	return s.coreACService.ChangeName(
		authUser,
		id,
		name,
	)
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

	authUser, err := s.authUser()
	if err != nil {
		return err
	}

	return s.coreACService.ChangePrivilege(
		authUser,
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

	authUser, err := s.authUser()
	if err != nil {
		return err
	}

	return s.coreACService.Revoke(
		authUser,
		id,
	)
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

	authUser, err := s.authUser()
	if err != nil {
		return nil, err
	}

	return s.coreACService.UserList(
		authUser,
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

	authUser, err := s.authUser()
	if err != nil {
		return nil, err
	}

	return s.coreACService.ApproveUser(
		authUser,
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

	authUser, err := s.authUser()
	if err != nil {
		return nil, err
	}

	return s.coreACService.RevokeUser(
		authUser,
		userID,
	)
}

func (s *Service) DeleteAccessControl(
	id string,
) error {
	if s.context.IsRemote {
		if s.context.SelectedSoT == nil {
			return fmt.Errorf(
				"selected source of truth is nil",
			)
		}

		return s.appDeps.Client.DeleteAccessControl(
			s.context.SelectedSoT.Address,
			id,
		)
	}

	authUser, err := s.authUser()
	if err != nil {
		return err
	}

	return s.coreACService.DeletePermission(
		authUser,
		id,
	)
}

func (s *Service) GrantPrivilege(
	id string,
	myVaultAccessList []entity.VaultAccess,
	keyPair *keyring.KeyPair,
) error {
	newVaultAccessList := make([]entity.VaultAccess, 0, len(myVaultAccessList))

	for _, vaultAccess := range myVaultAccessList {
		vaultKey, err := s.context.SelectedSoT.KeyPair.Open(vaultAccess.VaultKey)
		if err != nil {
			return err
		}

		encryptedVaultKey, err := keyPair.Seal(vaultKey)
		if err != nil {
			return err
		}

		newVaultAccessList = append(
			newVaultAccessList,
			entity.VaultAccess{
				VaultID:      vaultAccess.VaultID,
				PermissionID: id,
				VaultKey:     encryptedVaultKey,
			},
		)
	}

	if s.context.IsRemote {
		if s.context.SelectedSoT == nil {
			return fmt.Errorf(
				"selected source of truth is nil",
			)
		}

		return s.appDeps.Client.GrantVaultAccess(
			s.context.SelectedSoT.Address,
			newVaultAccessList,
		)
	}

	authUser, err := s.authUser()
	if err != nil {
		return err
	}

	return s.coreACService.GrantPrivilege(authUser, newVaultAccessList)
}
