package detail

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		if m.Focus == FocusSidebar {
			return m.updateSidebar(msg)
		}

		if m.Focus == FocusContent {
			return m.updateContent(msg)
		}
	}

	return m, nil
}
func (m Model) updateSidebar(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	case "up", "k":
		if m.SidebarCursor > 0 {
			m.SidebarCursor--
		}

	case "down", "j":
		if m.SidebarCursor < 2 {
			m.SidebarCursor++
		}

	case "enter", "right", "l":
		m.Focus = FocusContent

	case "esc", "q":
		m.Active = false
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
		if m.Cursor < len(m.Vaults) {
			m.Cursor++
		}

	case "enter":
		if m.SidebarCursor == 0 {
			if m.Cursor == len(m.Vaults) {
				// TODO: Add Vault
				return m, nil
			}

			// TODO: Open selected Vault
		}

	case "left", "h", "esc":
		m.Cursor = 0
		m.Focus = FocusSidebar

	case "q":
		m.Active = false
	}

	return m, nil
}
