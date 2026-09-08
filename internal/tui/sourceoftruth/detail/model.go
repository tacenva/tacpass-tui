package detail

import (
	"github.com/tacenva/tacpass-tui/internal/app"
	accessControlTUI "github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail/accesscontrol"
	vaultTUI "github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail/vault"
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
}

func New(
	appDeps *app.Deps,
	ctx *app.Context,
	masterKey string,
) Model {
	v := vaultTUI.New(
		appDeps,
		ctx,
		masterKey,
	)

	v.Load()

	a := accessControlTUI.NewModel(appDeps, ctx)

	return Model{
		Active:        true,
		SidebarCursor: 0,
		Cursor:        0,

		context: ctx,

		vaultTUI:         v,
		accessControlTUI: a,
	}
}
