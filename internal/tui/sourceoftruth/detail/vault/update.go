package vault

import (
	"fmt"

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
		fmt.Printf(
			"cursor=%d len=%d\n",
			m.Cursor,
			len(m.VaultAccessList),
		)
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		fmt.Printf(
			"cursor=%d len=%d\n",
			m.Cursor,
			len(m.VaultAccessList),
		)
		if m.Cursor < len(m.VaultAccessList) {
			m.Cursor++
		}

	case "enter":
		fmt.Printf(
			"cursor=%d len=%d\n",
			m.Cursor,
			len(m.VaultAccessList),
		)
		if m.Cursor == len(m.VaultAccessList) {
			m.Focus = FocusNewVault
			m.VaultName = ""

			return m, nil
		}

		selectedVaultAccess := m.VaultAccessList[m.Cursor]

		if err := m.VaultRecordTUI.Load(&selectedVaultAccess); err != nil {
			m.ErrorMessage = err.Error()
			return m, nil
		}
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

		vaultAccessData, err := m.VaultServiceTUI.CreateVault(m.VaultName)
		if err != nil {
			m.ErrorMessage = err.Error()
			return m, nil
		}

		m.VaultAccessList = append(m.VaultAccessList, *vaultAccessData)

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
