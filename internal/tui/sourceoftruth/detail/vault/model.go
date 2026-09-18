package vault

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/tacenva-services/app"
	"github.com/tacenva/tacenva-services/app/sourceoftruth"
	"github.com/tacenva/tacenva-services/app/vault"
	coreApp "github.com/tacenva/tacpass-core/app"
	"github.com/tacenva/tacpass-core/entity"
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
	coreServices *coreApp.Services,
	sotService *sourceoftruth.Service,
	screenState *state.Async,
	actionState *state.Async,
) Model {
	vaultServiceTUI := vault.NewService(
		appDeps,
		context,
		masterKey,
		coreServices.Vault,
		coreServices.VaultRecordService,
		coreServices.Auth,
		sotService,
	)

	return Model{
		Active: true,

		context: context,

		VaultAccessList: []entity.VaultAccess{},
		VaultServiceTUI: vaultServiceTUI,

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
