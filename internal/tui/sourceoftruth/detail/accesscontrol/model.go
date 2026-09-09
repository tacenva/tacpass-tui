package accesscontrol

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"

	coreAC "github.com/tacenva/tacpass-core/accesscontrol"
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-core/util/keyring"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/app/accesscontrol"
	userListTUI "github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail/accesscontrol/userlist"
	"github.com/tacenva/tacpass-tui/internal/tui/state"
)

var errFormPermissionNil = errors.New("form permission is nil")

type Focus int

const (
	FocusNone Focus = iota
	FocusContent
	FocusForm
	FocusUsers
	FocusCreatedKey
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
	Permission *entity.Permission
	KeyPair    *keyring.KeyPair
	Err        error
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

	UserList userListTUI.Model

	ScreenState *state.Async
	ActionState *state.Async

	CreatedPrivateKey string
	CreatedPublicKey  string
}

func New(
	appDeps *app.Deps,
	context *app.Context,
	coreACService *coreAC.Service,
	coreAuthService *auth.Service,
	ScreenState *state.Async,
	ActionState *state.Async,
) Model {
	service := accesscontrol.NewService(
		appDeps,
		context,
		coreACService,
		coreAuthService,
	)

	return Model{
		appDeps: appDeps,
		context: context,
		service: service,

		Focus:  FocusNone,
		Cursor: 0,
		Active: false,

		Permissions: nil,

		FormMode:      FormNone,
		FormCursor:    0,
		FormPrivilege: entity.PrivilegeRead,

		ScreenState: ScreenState,
		ActionState: ActionState,
	}
}

func (m *Model) Load() tea.Cmd {
	m.ScreenState.Start()

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

func (m *Model) openUsers(
	permission entity.Permission,
) tea.Cmd {
	m.UserList = userListTUI.New(
		m.service,
		permission.ID,
		m.ScreenState,
		m.ActionState,
	)

	m.UserList.Active = true
	m.UserList.ScreenState.Start()
	m.UserList.Focus = userListTUI.FocusContent

	m.Focus = FocusUsers

	return m.UserList.Load()
}
