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

func New(dbDeps *app.Deps, context *app.Context) Model {
	vaultRepository := vault.NewRepository(dbDeps.SqliteDB)
	vaultaccessRepository := vaultaccess.NewRepository(dbDeps.SqliteDB)
	vaultaccessService := vaultaccess.NewService(vaultaccessRepository)
	vaultService := vault.NewService(
		vaultRepository,
		context.NodeDB,
		vaultaccessService,
	)

	return Model{
		Active: true,
		Cursor: 0,

		context: context,

		VaultList:    []entity.Vault{},
		VaultService: vaultService,

		VaultRecordTUI: vaultrecord.New(context, vaultService),
	}
}

func (m *Model) Load() error {
	vaultList, err := m.VaultService.VaultList(m.context.AuthUser)
	if err != nil {
		return err
	}

	m.VaultList = vaultList
	return nil
}
