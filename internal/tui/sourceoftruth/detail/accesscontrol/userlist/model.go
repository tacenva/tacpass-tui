package userlist

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-tui/internal/app/accesscontrol"
)

type Focus int

const (
	FocusNone Focus = iota
	FocusContent
)

type UsersLoadedMsg struct {
	Users []entity.User
	Err   error
}

type UserUpdatedMsg struct {
	User *entity.User
	Err  error
}

type Model struct {
	service *accesscontrol.Service

	PermissionID string

	Width  int
	Height int

	Focus  Focus
	Cursor int
	Active bool

	Users []entity.User
}

func New(
	service *accesscontrol.Service,
	permissionID string,
) Model {
	return Model{
		service: service,

		PermissionID: permissionID,

		Focus:  FocusNone,
		Cursor: 0,
		Active: false,

		Users: nil,
	}
}

func (m Model) Load() tea.Cmd {
	return func() tea.Msg {
		users, err := m.service.UserList(
			m.PermissionID,
		)

		return UsersLoadedMsg{
			Users: users,
			Err:   err,
		}
	}
}
