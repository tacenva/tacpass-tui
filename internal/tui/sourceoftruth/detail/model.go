package detail

import (
	"github.com/tacenva/tacpass-tui/internal/app"
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

	vaultTUI vaultTUI.Model
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

	return Model{
		Active:        true,
		SidebarCursor: 0,
		Cursor:        0,

		context: ctx,

		vaultTUI: v,
	}
}
