package detail

import (
	"github.com/tacenva/tacenva-services/app"
	"github.com/tacenva/tacenva-services/app/sourceoftruth"

	coreApp "github.com/tacenva/tacpass-core/app"
	accessControlTUI "github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail/accesscontrol"
	vaultTUI "github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail/vault"
	"github.com/tacenva/tacpass-tui/internal/tui/state"
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

	context *app.Context

	vaultTUI         vaultTUI.Model
	accessControlTUI accessControlTUI.Model

	ScreenState *state.Async
	ActionState *state.Async
}

func New(
	appDeps *app.Deps,
	ctx *app.Context,
	masterKey string,
	coreServices *coreApp.Services,
	sotService *sourceoftruth.Service,
	screenState *state.Async,
	actionState *state.Async,
) Model {
	v := vaultTUI.New(
		appDeps,
		ctx,
		masterKey,
		coreServices,
		sotService,
		screenState,
		actionState,
	)

	a := accessControlTUI.New(
		appDeps,
		ctx,
		coreServices,
		screenState,
		actionState,
	)

	return Model{
		Active:        true,
		SidebarCursor: 0,
		Cursor:        0,

		context: ctx,

		vaultTUI:         v,
		accessControlTUI: a,

		ScreenState: screenState,
		ActionState: screenState,
	}
}
