package vault

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/entity"
	coreVault "github.com/tacenva/tacpass-core/vault"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/app/sourceoftruth"
	"github.com/tacenva/tacpass-tui/internal/app/vault"
	vaultrecord "github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail/vault/record"
	"github.com/tacenva/tacpass-tui/internal/tui/state"
)

type Focus int

const (
	FocusNone Focus = iota
	FocusContent
	FocusNewVault
)

type VaultsLoadedMsg struct {
	VaultAccessList []entity.VaultAccess
	Err             error
}

type Model struct {
	Width  int
	Height int

	Cursor int

	Focus  Focus
	Active bool

	VaultName string
	Editing   bool

	VaultAccessList []entity.VaultAccess

	VaultServiceTUI *vault.Service
	VaultRecordTUI  vaultrecord.Model

	context *app.Context

	ScreenState *state.Async
	ActionState *state.Async
}

func New(
	appDeps *app.Deps,
	context *app.Context,
	masterKey string,
	coreVaultService *coreVault.Service,
	authService *auth.Service,
	sotService *sourceoftruth.Service,
	screenState *state.Async,
	actionState *state.Async,
) Model {
	vaultServiceTUI := vault.NewService(
		appDeps,
		context,
		masterKey,
		coreVaultService,
		authService,
		sotService,
	)

	return Model{
		Active: true,

		context: context,

		VaultAccessList: []entity.VaultAccess{},
		VaultServiceTUI: vaultServiceTUI,

		VaultRecordTUI: vaultrecord.New(
			context,
			vaultServiceTUI,
			screenState,
			actionState,
		),

		ScreenState: screenState,
		ActionState: actionState,
	}
}

func (m Model) Load() tea.Cmd {
	m.ScreenState.Start()

	return func() tea.Msg {
		vaultAccessList, _, err := m.VaultServiceTUI.List()

		return VaultsLoadedMsg{
			VaultAccessList: vaultAccessList,
			Err:             err,
		}
	}
}
