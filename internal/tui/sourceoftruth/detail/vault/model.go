package vault

import (
	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-core/vault"
	"github.com/tacenva/tacpass-core/vaultaccess"
	"github.com/tacenva/tacpass-tui/internal/app"
	vaultrecord "github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail/vault/record"
)

type Focus int

const (
	FocusNone Focus = iota
	FocusContent
	FocusNewVault
)

type Model struct {
	Width  int
	Height int

	Cursor int

	Focus  Focus
	Active bool

	VaultName string

	VaultList    []entity.Vault
	VaultService *vault.Service

	context        *app.Context
	VaultRecordTUI vaultrecord.Model
}

func New(deps *app.DatabaseDeps, context *app.Context) Model {
	vaultRepository := vault.NewRepository(deps.SqliteDB)
	vaultaccessRepository := vaultaccess.NewRepository(deps.SqliteDB)
	vaultaccessService := vaultaccess.NewService(vaultaccessRepository)

	return Model{
		Active: true,
		Cursor: 0,

		context: context,

		VaultList: []entity.Vault{},
		VaultService: vault.NewService(
			vaultRepository,
			deps.TacenvaDB,
			context.AuthService,
			context.PermissionService,
			vaultaccessService,
		),

		VaultRecordTUI: vaultrecord.New(),
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
