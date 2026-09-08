package vault

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/tacpass-core/entity"
)

type VaultCreatedMsg struct {
	VaultAccess *entity.VaultAccess
	Err         error
}

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

	case VaultsLoadedMsg:
		if msg.Err != nil {
			m.ScreenState.Fail(msg.Err)
			return m, nil
		}

		m.VaultAccessList = msg.VaultAccessList
		m.normalizeCursor()
		m.ScreenState.Success()

	case VaultCreatedMsg:
		if msg.Err != nil {
			m.ActionState.Fail(msg.Err)
			return m, nil
		}

		m.VaultAccessList = append(
			m.VaultAccessList,
			*msg.VaultAccess,
		)

		m.VaultName = ""
		m.Focus = FocusContent
		m.ActionState.Success()

	case tea.KeyMsg:
		if m.ScreenState.Loading ||
			m.ActionState.Loading {
			return m, nil
		}

		switch m.Focus {
		case FocusContent:
			return m.updateContent(msg)

		case FocusNewVault:
			return m.updateNewVault(msg)
		}
	}

	return m, nil
}

func (m *Model) normalizeCursor() {
	if len(m.VaultAccessList) == 0 {
		m.Cursor = 0
		return
	}

	if m.Cursor >= len(m.VaultAccessList) {
		m.Cursor = len(m.VaultAccessList) - 1
	}
}

func (m Model) updateContent(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if m.Cursor < len(m.VaultAccessList) {
			m.Cursor++
		}

	case "enter":
		if m.Cursor == len(m.VaultAccessList) {
			m.Focus = FocusNewVault
			m.VaultName = ""

			return m, nil
		}

		selectedVaultAccess := m.VaultAccessList[m.Cursor]

		if err := m.VaultRecordTUI.Load(
			&selectedVaultAccess,
		); err != nil {
			m.ScreenState.Fail(err)
			return m, nil
		}

		m.VaultRecordTUI.Active = true

	case "esc":
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

		m.ActionState.Start()

		name := m.VaultName

		return m, func() tea.Msg {
			vaultAccessData, err := m.VaultServiceTUI.CreateVault(name)

			return VaultCreatedMsg{
				VaultAccess: vaultAccessData,
				Err:         err,
			}
		}

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
