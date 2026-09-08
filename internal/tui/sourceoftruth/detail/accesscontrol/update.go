package accesscontrol

import (
	tea "github.com/charmbracelet/bubbletea"
	userListTUI "github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail/accesscontrol/userlist"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case PermissionsLoadedMsg:
		if !m.ScreenState.Loading {
			return m, nil
		}

		if msg.Err != nil {
			m.ScreenState.Fail(msg.Err)
			return m, nil
		}

		m.Permissions = msg.Permissions

		if m.Cursor > len(m.Permissions) {
			m.Cursor = len(m.Permissions)
		}

		m.ScreenState.Success()

		return m, nil

	case PermissionSavedMsg:
		if msg.Err != nil {
			m.ActionState.Fail(msg.Err)
			return m, nil
		}

		m.ActionState.Success()
		m.closeForm()

		return m, m.Load()

	case userListTUI.UsersLoadedMsg:
		updated, cmd := m.UserList.Update(msg)
		m.UserList = updated

		return m, cmd

	case tea.KeyMsg:
		if m.ScreenState.Loading {
			return m, nil
		}

		switch m.Focus {
		case FocusContent:
			return m.updateContent(msg)

		case FocusForm:
			return m.updateForm(msg)

		case FocusUsers:
			updated, cmd := m.UserList.Update(msg)
			m.UserList = updated

			if m.UserList.Focus == userListTUI.FocusNone {
				m.UserList.Active = false
				m.Focus = FocusContent
				cmd = nil
			}

			return m, cmd
		}
	}

	return m, nil
}

func (m Model) updateContent(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if m.Cursor < len(m.Permissions) {
			m.Cursor++
		}

	case "enter":
		if m.Cursor == len(m.Permissions) {
			m.openNewForm()
			return m, nil
		}

		if m.Cursor >= 0 && m.Cursor < len(m.Permissions) {
			return m, m.openUsers(m.Permissions[m.Cursor])
		}

	case "e":
		if m.Cursor >= 0 && m.Cursor < len(m.Permissions) {
			m.openUpdateForm(m.Permissions[m.Cursor])
		}

	case "esc":
		m.Focus = FocusNone
		m.Active = false
	}

	return m, nil
}

func (m Model) updateForm(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.closeForm()
		return m, nil

	case "up", "k":
		if m.FormCursor > 0 {
			m.FormCursor--
		}

	case "down", "j":
		if m.FormCursor < 1 {
			m.FormCursor++
		}

	case "tab":
		if m.FormCursor < 1 {
			m.FormCursor++
		} else {
			m.FormCursor = 0
		}

	case "shift+tab":
		if m.FormCursor > 0 {
			m.FormCursor--
		} else {
			m.FormCursor = 1
		}

	case "enter":
		if m.FormCursor == 0 {
			m.FormCursor = 1
			return m, nil
		}

		return m, m.saveForm()

	case "ctrl+s":
		return m, m.saveForm()
	}

	if m.FormCursor == 0 {
		switch msg.Type {
		case tea.KeyRunes:
			m.FormName = append(m.FormName, msg.Runes...)
		}

		if msg.Type == tea.KeyBackspace ||
			msg.Type == tea.KeyDelete {
			if len(m.FormName) > 0 {
				m.FormName = m.FormName[:len(m.FormName)-1]
			}
		}

		return m, nil
	}

	switch msg.String() {
	case "left", "h":
		m.cyclePrivilege(-1)

	case "right", "l", " ":
		m.cyclePrivilege(1)
	}

	return m, nil
}

func (m *Model) saveForm() tea.Cmd {
	name := m.formName()
	privilege := m.FormPrivilege

	m.ActionState.Start()

	switch m.FormMode {
	case FormNew:
		return func() tea.Msg {
			_, _, err := m.service.Create(
				name,
				privilege,
			)

			return PermissionSavedMsg{
				Err: err,
			}
		}

	case FormUpdate:
		if m.FormPermission == nil {
			return func() tea.Msg {
				return PermissionSavedMsg{
					Err: errFormPermissionNil,
				}
			}
		}

		permissionID := m.FormPermission.ID
		oldName := m.FormPermission.Name
		oldPrivilege := m.FormPermission.Privilege

		return func() tea.Msg {
			if name != oldName {
				if err := m.service.ChangeName(
					permissionID,
					name,
				); err != nil {
					return PermissionSavedMsg{
						Err: err,
					}
				}
			}

			if privilege != oldPrivilege {
				if err := m.service.ChangePrivilege(
					permissionID,
					privilege,
				); err != nil {
					return PermissionSavedMsg{
						Err: err,
					}
				}
			}

			return PermissionSavedMsg{}
		}
	}

	return nil
}
