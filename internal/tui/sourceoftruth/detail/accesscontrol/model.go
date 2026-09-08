package accesscontrol

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/app/accesscontrol"
)

type Focus int

const (
	FocusNone Focus = iota
	FocusContent
)

type PermissionsLoadedMsg struct {
	Permissions []entity.Permission
	Err         error
}

type Model struct {
	appDeps *app.Deps
	context *app.Context
	service *accesscontrol.Service

	Width  int
	Height int

	Focus  Focus
	Cursor int
	Active bool

	Permissions []entity.Permission
}

func New(
	appDeps *app.Deps,
	context *app.Context,
) Model {
	return Model{
		appDeps: appDeps,
		context: context,

		service: accesscontrol.NewService(
			appDeps,
			context,
		),

		Focus:  FocusNone,
		Cursor: 0,
		Active: false,

		Permissions: nil,
	}
}

func (m Model) Load() tea.Cmd {
	return func() tea.Msg {
		permissions, err := m.service.List()

		return PermissionsLoadedMsg{
			Permissions: permissions,
			Err:         err,
		}
	}
}
