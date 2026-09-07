package vault

import (
	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/app/vault"
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

	VaultAccessList []entity.VaultAccess
	VaultServiceTUI *vault.Service

	context        *app.Context
	VaultRecordTUI vaultrecord.Model
}

func New(appDeps *app.Deps, context *app.Context, masterKey string) Model {
	vaultServiceTUI := vault.NewService(appDeps, context, masterKey)
	return Model{
		Active: true,
		Cursor: 0,

		context: context,

		VaultAccessList: []entity.VaultAccess{},
		VaultServiceTUI: vaultServiceTUI,

		VaultRecordTUI: vaultrecord.New(context, vaultServiceTUI),
	}
}

func (m *Model) Load() error {
	vaultAccessList, err := m.VaultServiceTUI.List()
	if err != nil {
		return err
	}

	m.VaultAccessList = vaultAccessList
	return nil
}
