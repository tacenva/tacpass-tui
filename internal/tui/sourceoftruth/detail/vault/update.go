package vault

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.VaultRecordTUI.Active {
		updated, cmd := m.VaultRecordTUI.Update(msg)
		m.VaultRecordTUI = updated

		return m, cmd
	}
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch m.Focus {
		case FocusContent:
			return m.updateContent(msg)
		case FocusNewVault:
			return m.updateNewVault(msg)
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
		if m.Cursor < len(m.VaultList) {
			m.Cursor++
		}

	case "enter":
		if m.Cursor == len(m.VaultList) {
			m.Focus = FocusNewVault
			m.VaultName = ""

			return m, nil
		}

		selectedVault := m.VaultList[m.Cursor]
		m.VaultRecordTUI.SelectedVault = &selectedVault
		m.VaultRecordTUI.Active = true

		return m, nil

	case "left", "h", "esc":
		m.Focus = FocusNone

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

		vaultData, err := m.VaultService.Create(m.VaultName, m.context.SelectedSoT.AuthToken)
		if err != nil {
			return m, nil
		}

		m.VaultList = append(m.VaultList, *vaultData)

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
