package accesscontrol

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/tacenva/tacpass-core/entity"
)

type permissionsLoadedMsg struct {
	Permissions []entity.Permission
}

type usersLoadedMsg struct {
	Users []entity.User
}

type accessControlCreatedMsg struct {
	Permission *entity.Permission
}

type privilegeChangedMsg struct {
	ID        string
	Privilege entity.Privilege
}

type accessControlRevokedMsg struct {
	ID string
}

type userApprovedMsg struct {
	User *entity.User
}

type userRevokedMsg struct {
	User *entity.User
}

type errorMsg struct {
	err error
}

func (m Model) Init() tea.Cmd {
	return m.loadPermissions()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

	case permissionsLoadedMsg:
		m.Permissions = msg.Permissions

		if m.Cursor >= len(m.Permissions) {
			m.Cursor = len(m.Permissions) - 1
		}

		if m.Cursor < 0 {
			m.Cursor = 0
		}

		m.Error = nil
		m.Message = "permissions refreshed"

		return m, nil

	case usersLoadedMsg:
		m.Users = msg.Users

		if m.UserCursor >= len(m.Users) {
			m.UserCursor = len(m.Users) - 1
		}

		if m.UserCursor < 0 {
			m.UserCursor = 0
		}

		m.Error = nil
		m.Message = "users refreshed"

		return m, nil

	case accessControlCreatedMsg:
		if msg.Permission == nil {
			m.Error = fmt.Errorf(
				"daemon returned empty permission",
			)
			m.Message = ""

			return m, nil
		}

		m.Permissions = append(
			m.Permissions,
			*msg.Permission,
		)

		m.Cursor = len(m.Permissions) - 1
		m.Page = PageList

		m.Error = nil
		m.Message = fmt.Sprintf(
			"permission created: %s",
			msg.Permission.ID,
		)

		return m, nil

	case privilegeChangedMsg:
		for i := range m.Permissions {
			if m.Permissions[i].ID == msg.ID {
				m.Permissions[i].Privilege = msg.Privilege
				break
			}
		}

		if m.SelectedPermission != nil &&
			m.SelectedPermission.ID == msg.ID {
			m.SelectedPermission.Privilege = msg.Privilege
		}

		m.Error = nil
		m.Message = fmt.Sprintf(
			"privilege updated: %s",
			msg.Privilege,
		)

		return m, nil

	case accessControlRevokedMsg:
		for i := range m.Permissions {
			if m.Permissions[i].ID == msg.ID {
				m.Permissions[i].Revoked = true
				break
			}
		}

		if m.SelectedPermission != nil &&
			m.SelectedPermission.ID == msg.ID {
			m.SelectedPermission.Revoked = true
		}

		m.Error = nil
		m.Message = "permission revoked"

		return m, nil

	case userApprovedMsg:
		if msg.User == nil {
			return m, nil
		}

		for i := range m.Users {
			if m.Users[i].ID == msg.User.ID {
				m.Users[i] = *msg.User
				break
			}
		}

		m.Error = nil
		m.Message = fmt.Sprintf(
			"user approved: %s",
			msg.User.ID,
		)

		return m, nil

	case userRevokedMsg:
		if msg.User == nil {
			return m, nil
		}

		for i := range m.Users {
			if m.Users[i].ID == msg.User.ID {
				m.Users[i] = *msg.User
				break
			}
		}

		m.Error = nil
		m.Message = fmt.Sprintf(
			"user revoked: %s",
			msg.User.ID,
		)

		return m, nil

	case errorMsg:
		m.Error = msg.err
		m.Message = ""

		return m, nil
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.Page {
	case PageList:
		return m.updateList(msg)

	case PageUsers:
		return m.updateUsers(msg)

	case PageCreate:
		return m.updateCreate(msg)
	}

	return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if m.Cursor < len(m.Permissions)-1 {
			m.Cursor++
		}

	case "enter", "right", "l":
		if len(m.Permissions) == 0 {
			return m, nil
		}

		m.SelectedPermission = &m.Permissions[m.Cursor]
		m.Users = nil
		m.UserCursor = 0
		m.Error = nil
		m.Message = ""

		m.Page = PageUsers

		return m, m.loadUsers(
			m.SelectedPermission.ID,
		)

	case "n":
		m.CreatePrivilege = entity.PrivilegeRead
		m.Error = nil
		m.Message = ""
		m.Page = PageCreate

	case "r":
		m.Error = nil
		m.Message = ""

		return m, m.loadPermissions()

	case "esc", "left", "h", "q":
		m.Active = false
		m.Focus = FocusNone
	}

	return m, nil
}

func (m Model) updateUsers(msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.SelectedPermission == nil {
		m.Page = PageList
		return m, nil
	}

	switch msg.String() {
	case "up", "k":
		if m.UserCursor > 0 {
			m.UserCursor--
		}

	case "down", "j":
		if m.UserCursor < len(m.Users)-1 {
			m.UserCursor++
		}

	case "a":
		if len(m.Users) == 0 {
			return m, nil
		}

		user := m.Users[m.UserCursor]

		if user.Status == entity.UserStatusApproved {
			m.Message = "user already approved"
			return m, nil
		}

		m.Error = nil
		m.Message = ""

		return m, m.approveUser(user.ID)

	case "x":
		if len(m.Users) == 0 {
			return m, nil
		}

		user := m.Users[m.UserCursor]

		if user.Status == entity.UserStatusRevoked {
			m.Message = "user already revoked"
			return m, nil
		}

		m.Error = nil
		m.Message = ""

		return m, m.revokeUser(user.ID)

	case "r":
		m.Error = nil
		m.Message = ""

		return m, m.loadUsers(
			m.SelectedPermission.ID,
		)

	case "esc", "b", "left", "h":
		m.Page = PageList
		m.SelectedPermission = nil
		m.Users = nil
		m.Error = nil
		m.Message = ""

	case "q":
		m.Active = false
		m.Focus = FocusNone
	}

	return m, nil
}

func (m Model) updateCreate(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		m.CreatePrivilege = entity.PrivilegeAdmin
		m.Message = "selected privilege: admin"

	case "2":
		m.CreatePrivilege = entity.PrivilegeWrite
		m.Message = "selected privilege: write"

	case "3":
		m.CreatePrivilege = entity.PrivilegeRead
		m.Message = "selected privilege: read"

	case "enter":
		m.Error = nil
		m.Message = ""

		return m, m.createAccessControl(
			m.CreatePrivilege,
		)

	case "esc", "b", "left", "h":
		m.Page = PageList
		m.Error = nil
		m.Message = ""

	case "q":
		m.Active = false
		m.Focus = FocusNone
	}

	return m, nil
}

func (m Model) loadPermissions() tea.Cmd {
	return func() tea.Msg {
		permissions, err := m.Service.List()
		if err != nil {
			return errorMsg{
				err: err,
			}
		}

		return permissionsLoadedMsg{
			Permissions: permissions,
		}
	}
}

func (m Model) loadUsers(
	permissionID string,
) tea.Cmd {
	return func() tea.Msg {
		users, err := m.Service.ListUsers(
			permissionID,
		)
		if err != nil {
			return errorMsg{
				err: err,
			}
		}

		return usersLoadedMsg{
			Users: users,
		}
	}
}

func (m Model) createAccessControl(
	privilege entity.Privilege,
) tea.Cmd {
	return func() tea.Msg {
		response, err := m.Service.Create(
			privilege,
		)
		if err != nil {
			return errorMsg{
				err: err,
			}
		}

		if response == nil ||
			response.Permission == nil {
			return errorMsg{
				err: fmt.Errorf(
					"daemon returned empty permission",
				),
			}
		}

		return accessControlCreatedMsg{
			Permission: response.Permission,
		}
	}
}

func (m Model) changePrivilege(
	id string,
	privilege entity.Privilege,
) tea.Cmd {
	return func() tea.Msg {
		if err := m.Service.ChangePrivilege(
			id,
			privilege,
		); err != nil {
			return errorMsg{
				err: err,
			}
		}

		return privilegeChangedMsg{
			ID:        id,
			Privilege: privilege,
		}
	}
}

func (m Model) revokeAccessControl(
	id string,
) tea.Cmd {
	return func() tea.Msg {
		if err := m.Service.Revoke(
			id,
		); err != nil {
			return errorMsg{
				err: err,
			}
		}

		return accessControlRevokedMsg{
			ID: id,
		}
	}
}

func (m Model) approveUser(
	userID string,
) tea.Cmd {
	return func() tea.Msg {
		user, err := m.Service.ApproveUser(
			userID,
		)
		if err != nil {
			return errorMsg{
				err: err,
			}
		}

		return userApprovedMsg{
			User: user,
		}
	}
}

func (m Model) revokeUser(
	userID string,
) tea.Cmd {
	return func() tea.Msg {
		user, err := m.Service.RevokeUser(
			userID,
		)
		if err != nil {
			return errorMsg{
				err: err,
			}
		}

		return userRevokedMsg{
			User: user,
		}
	}
}
