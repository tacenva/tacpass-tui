package detail

import (
	tea "github.com/charmbracelet/bubbletea"

	accessControlTUI "github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail/accesscontrol"
	vaultTUI "github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail/vault"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		return m, nil

	case accessControlTUI.PermissionsLoadedMsg:
		updated, cmd := m.accessControlTUI.Update(msg)

		m.accessControlTUI = updated

		return m, cmd

	case tea.KeyMsg:
		switch m.Focus {
		case FocusSidebar:
			return m.updateSidebar(msg)

		case FocusContent:
			switch m.SidebarCursor {
			case 0:
				updated, cmd := m.vaultTUI.Update(msg)

				m.vaultTUI = updated

				if m.vaultTUI.Focus == vaultTUI.FocusNone {
					m.Focus = FocusSidebar
					cmd = nil
				}

				return m, cmd

			case 1:
				updated, cmd := m.accessControlTUI.Update(msg)

				m.accessControlTUI = updated

				if m.accessControlTUI.Focus == accessControlTUI.FocusNone {
					m.accessControlTUI.Active = false
					m.Focus = FocusSidebar
					cmd = nil
				}

				return m, cmd
			}
		}

		return m, nil
	}

	switch m.Focus {
	case FocusContent:
		switch m.SidebarCursor {
		case 0:
			updated, cmd := m.vaultTUI.Update(msg)

			m.vaultTUI = updated

			if m.vaultTUI.Focus == vaultTUI.FocusNone {
				m.Focus = FocusSidebar
				cmd = nil
			}

			return m, cmd

		case 1:
			updated, cmd := m.accessControlTUI.Update(msg)

			m.accessControlTUI = updated

			if m.accessControlTUI.Focus == accessControlTUI.FocusNone {
				m.accessControlTUI.Active = false
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
	case "up":
		if m.SidebarCursor > 0 {
			m.SidebarCursor--
		}

		if m.SidebarCursor == 1 {
			return m, m.accessControlTUI.Load()
		}

	case "down":
		if m.SidebarCursor < 2 {
			m.SidebarCursor++
		}

		if m.SidebarCursor == 1 {
			return m, m.accessControlTUI.Load()
		}

	case "enter":
		m.Focus = FocusContent

		switch m.SidebarCursor {
		case 0:
			m.vaultTUI.Focus = vaultTUI.FocusContent

		case 1:
			m.accessControlTUI.Active = true
			m.accessControlTUI.Focus = accessControlTUI.FocusContent

		case 2:
			// Setting nanti.
		}

	case "esc", "q":
		m.Active = false
	}

	return m, nil
}
