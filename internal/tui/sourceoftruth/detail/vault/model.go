package vault

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/entity"
	coreVault "github.com/tacenva/tacpass-core/vault"
	"github.com/tacenva/tacpass-tui/internal/operations/app"
	"github.com/tacenva/tacpass-tui/internal/operations/app/sourceoftruth"
	"github.com/tacenva/tacpass-tui/internal/operations/app/vault"
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
	NeedSync        bool
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

	needSync bool
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

		// VaultRecordTUI: vaultrecord.New(
		// 	context,
		// 	vaultServiceTUI,
		// 	screenState,
		// 	actionState,
		// ),

		ScreenState: screenState,
		ActionState: actionState,
		needSync:    false,
	}
}

func (m Model) Load() tea.Cmd {
	m.ScreenState.Start()

	return func() tea.Msg {
		vaultAccessList, needSync, err := m.VaultServiceTUI.List()

		return VaultsLoadedMsg{
			VaultAccessList: vaultAccessList,
			NeedSync:        needSync,
			Err:             err,
		}
	}
}

func (m Model) Sync() tea.Cmd {
	m.ScreenState.Start()

	return func() tea.Msg {
		vaultAccessList, err := m.VaultServiceTUI.Sync()

		return VaultsLoadedMsg{
			VaultAccessList: vaultAccessList,
			NeedSync:        err != nil,
			Err:             err,
		}
	}
}
