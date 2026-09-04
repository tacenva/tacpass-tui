package detail

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/tacpass-core/entity"
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
			return m.updateContent(msg)
		case FocusNewVault:
			return m.updateNewVault(msg)
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
	if m.Focus == FocusNewVault {

	}

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
				m.Focus = FocusNewVault
				m.VaultName = ""
				return m, nil
			}

			// TODO: Open selected Vault
		}

	case "left", "h", "esc":
		m.Focus = FocusSidebar

	case "q":
		m.Active = false
	}

	return m, nil
}

func (m Model) updateNewVault(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	case "esc", "ctrl+c":
		m.Focus = FocusContent
		m.VaultName = ""

	case "enter":
		if m.VaultName == "" {
			return m, nil
		}

		m.Vaults = append(m.Vaults, entity.Vault{
			ID:   "vault-new",
			Name: m.VaultName,
		})

		m.Cursor = len(m.Vaults) - 1
		m.VaultName = ""
		m.Focus = FocusContent

	case "backspace":
		if len(m.VaultName) > 0 {
			m.VaultName = m.VaultName[:len(m.VaultName)-1]
		}

	default:
		if len(msg.Runes) > 0 {
			m.VaultName += string(msg.Runes)
		}
	}

	return m, nil
}
