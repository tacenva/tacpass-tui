package detail

import (
	tea "github.com/charmbracelet/bubbletea"
	vaultTUI "github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail/vault"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch m.Focus {
		case FocusSidebar:
			return m.updateSidebar(msg)
		case FocusContent:
			updated, cmd := m.vaultTUI.Update(msg)
			m.vaultTUI = updated
			if m.vaultTUI.Focus == vaultTUI.FocusNone {
				m.Focus = FocusSidebar
				cmd = nil
			}
			return m, cmd
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
		m.vaultTUI.Focus = vaultTUI.FocusContent

	case "esc", "q":
		m.Active = false
	}

	return m, nil
}
