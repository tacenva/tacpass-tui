package detail

import (
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-core/permission"
	"github.com/tacenva/tacpass-core/vault"
	"github.com/tacenva/tacpass-core/vaultaccess"
	"github.com/tacenva/tacpass-tui/internal/app"
	tuiEntity "github.com/tacenva/tacpass-tui/internal/entity"
)

type Focus int

const (
	FocusSidebar Focus = iota
	FocusContent
	FocusNewVault
)

type Model struct {
	Width  int
	Height int

	SidebarCursor int
	Cursor        int

	Focus  Focus
	Active bool

	VaultName string

	VaultList    []entity.Vault
	VaultService *vault.Service

	SelectedSoT *tuiEntity.SourceOfTruth
}

func New(deps *app.Dependencies, SelectedSoT *tuiEntity.SourceOfTruth, authService *auth.Service, permissionService *permission.Service) Model {
	vaultRepository := vault.NewRepository(deps.SqliteDB)
	vaultaccessRepository := vaultaccess.NewRepository(deps.SqliteDB)
	vaultaccessService := vaultaccess.NewService(vaultaccessRepository)

	return Model{
		Active:        true,
		SidebarCursor: 0,
		Cursor:        0,

		SelectedSoT: SelectedSoT,

		VaultList:    []entity.Vault{},
		VaultService: vault.NewService(vaultRepository, deps.TacenvaDB, authService, permissionService, vaultaccessService),
	}
}

func (m *Model) Load() error {
	vaultList, err := m.VaultService.List()
	if err != nil {
		return err
	}

	m.VaultList = vaultList
	return nil
}
