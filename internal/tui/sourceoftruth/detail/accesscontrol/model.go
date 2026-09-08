package accesscontrol

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/app/accesscontrol"
)

var errFormPermissionNil = errors.New("form permission is nil")

type Focus int

const (
	FocusNone Focus = iota
	FocusContent
	FocusForm
)

type FormMode int

const (
	FormNone FormMode = iota
	FormNew
	FormUpdate
)

type PermissionsLoadedMsg struct {
	Permissions []entity.Permission
	Err         error
}

type PermissionSavedMsg struct {
	Err error
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

	FormMode       FormMode
	FormCursor     int
	FormName       []rune
	FormPrivilege  entity.Privilege
	FormPermission *entity.Permission
}

func New(
	appDeps *app.Deps,
	context *app.Context,
) Model {
	return Model{
		appDeps: appDeps,
		context: context,
		service: accesscontrol.NewService(appDeps, context),

		Focus:  FocusNone,
		Cursor: 0,
		Active: false,

		Permissions: nil,

		FormMode:      FormNone,
		FormCursor:    0,
		FormPrivilege: entity.PrivilegeRead,
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

func (m *Model) openNewForm() {
	m.FormMode = FormNew
	m.FormCursor = 0
	m.FormName = nil
	m.FormPrivilege = entity.PrivilegeRead
	m.FormPermission = nil
	m.Focus = FocusForm
}

func (m *Model) openUpdateForm(permission entity.Permission) {
	m.FormMode = FormUpdate
	m.FormCursor = 0
	m.FormName = []rune(permission.Name)
	m.FormPrivilege = permission.Privilege
	m.FormPermission = &permission
	m.Focus = FocusForm
}

func (m *Model) closeForm() {
	m.FormMode = FormNone
	m.FormCursor = 0
	m.FormName = nil
	m.FormPrivilege = entity.PrivilegeRead
	m.FormPermission = nil
	m.Focus = FocusContent
}

func (m Model) formName() string {
	return string(m.FormName)
}

func (m *Model) cyclePrivilege(direction int) {
	privileges := []entity.Privilege{
		entity.PrivilegeAdmin,
		entity.PrivilegeWrite,
		entity.PrivilegeRead,
	}

	current := 0

	for i, privilege := range privileges {
		if privilege == m.FormPrivilege {
			current = i
			break
		}
	}

	current += direction

	if current < 0 {
		current = len(privileges) - 1
	}

	if current >= len(privileges) {
		current = 0
	}

	m.FormPrivilege = privileges[current]
}
