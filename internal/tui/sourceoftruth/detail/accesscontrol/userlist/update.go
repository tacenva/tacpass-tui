package userlist

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		return m, nil

	case UsersLoadedMsg:
		if msg.Err != nil {
			return m, nil
		}

		m.Users = msg.Users

		if len(m.Users) == 0 {
			m.Cursor = 0
			return m, nil
		}

		if m.Cursor >= len(m.Users) {
			m.Cursor = len(m.Users) - 1
		}

		return m, nil

	case UserUpdatedMsg:
		if msg.Err != nil {
			return m, nil
		}

		return m, m.Load()

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Users)-1 {
				m.Cursor++
			}

		case "a":
			return m, m.approveUser()

		case "r":
			return m, m.revokeUser()

		case "enter":
			// User detail nanti.

		case "esc":
			m.Focus = FocusNone
			m.Active = false
		}
	}

	return m, nil
}

func (m Model) approveUser() tea.Cmd {
	if len(m.Users) == 0 {
		return nil
	}

	if m.Cursor < 0 || m.Cursor >= len(m.Users) {
		return nil
	}

	userID := m.Users[m.Cursor].ID

	return func() tea.Msg {
		user, err := m.service.ApproveUser(userID)

		return UserUpdatedMsg{
			User: user,
			Err:  err,
		}
	}
}

func (m Model) revokeUser() tea.Cmd {
	if len(m.Users) == 0 {
		return nil
	}

	if m.Cursor < 0 || m.Cursor >= len(m.Users) {
		return nil
	}

	userID := m.Users[m.Cursor].ID

	return func() tea.Msg {
		user, err := m.service.RevokeUser(userID)

		return UserUpdatedMsg{
			User: user,
			Err:  err,
		}
	}
}
