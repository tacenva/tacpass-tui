package accesscontrol

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case PermissionsLoadedMsg:
		if msg.Err != nil {
			return m, nil
		}

		m.Permissions = msg.Permissions

		if len(m.Permissions) == 0 {
			m.Cursor = 0
			return m, nil
		}

		if m.Cursor >= len(m.Permissions) {
			m.Cursor = len(m.Permissions) - 1
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Permissions)-1 {
				m.Cursor++
			}

		case "esc":
			m.Focus = FocusNone
			m.Active = false

		case "enter":
			// Select permission nanti.
		}
	}

	return m, nil
}
