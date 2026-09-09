package vault

import (
	"encoding/json"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-tui/util/debug"
)

type VaultCreatedMsg struct {
	VaultAccess *entity.VaultAccess
	Err         error
}

type VaultUpdatedMsg struct {
	VaultAccess *entity.VaultAccess
	Err         error
}

type VaultDeletedMsg struct {
	VaultID string
	Err     error
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.VaultRecordTUI.Active {
		updated, cmd := m.VaultRecordTUI.Update(msg)
		m.VaultRecordTUI = updated

		if !m.VaultRecordTUI.Active {
			m.Focus = FocusContent
		}

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

		data, err := json.Marshal(msg.VaultAccessList)
		if err != nil {
			debug.Error(fmt.Sprintf(
				"failed to marshal vault access list: %v",
				err,
			))
		} else {
			debug.Print(string(data))
		}

		m.VaultAccessList = msg.VaultAccessList
		m.normalizeCursor()
		m.needSync = msg.NeedSync
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

	case VaultUpdatedMsg:
		if msg.Err != nil {
			m.ActionState.Fail(msg.Err)
			return m, nil
		}

		if msg.VaultAccess != nil {
			for i := range m.VaultAccessList {
				if m.VaultAccessList[i].ID == msg.VaultAccess.ID {
					m.VaultAccessList[i] = *msg.VaultAccess
					break
				}
			}
		}

		m.VaultName = ""
		m.Focus = FocusContent
		m.ActionState.Success()

	case VaultDeletedMsg:
		if msg.Err != nil {
			m.ActionState.Fail(msg.Err)
			return m, nil
		}

		for i := range m.VaultAccessList {
			if m.VaultAccessList[i].VaultID == msg.VaultID {
				m.VaultAccessList = append(
					m.VaultAccessList[:i],
					m.VaultAccessList[i+1:]...,
				)
				break
			}
		}

		m.normalizeCursor()
		m.ActionState.Success()

		return m, nil

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

	case "e":
		if len(m.VaultAccessList) == 0 ||
			m.Cursor >= len(m.VaultAccessList) {
			return m, nil
		}

		selectedVaultAccess := m.VaultAccessList[m.Cursor]

		m.Focus = FocusNewVault
		m.VaultName = selectedVaultAccess.Vault.Name
		m.Editing = true

	case "delete":
		if len(m.VaultAccessList) == 0 ||
			m.Cursor >= len(m.VaultAccessList) {
			return m, nil
		}

		selectedVaultAccess := m.VaultAccessList[m.Cursor]

		m.ActionState.Start()

		return m, func() tea.Msg {
			err := m.VaultServiceTUI.DeleteVault(
				&selectedVaultAccess.Vault,
			)

			return VaultDeletedMsg{
				VaultID: selectedVaultAccess.VaultID,
				Err:     err,
			}
		}

	case "s":
		m.ActionState.Start()
		return m, m.Sync()

	case "enter":
		if m.Cursor == len(m.VaultAccessList) {
			m.Focus = FocusNewVault
			m.VaultName = ""
			m.Editing = false

			return m, nil
		}

		if len(m.VaultAccessList) == 0 {
			return m, nil
		}

		selectedVaultAccess := m.VaultAccessList[m.Cursor]

		m.VaultRecordTUI.Active = true

		return m, m.VaultRecordTUI.Load(
			&selectedVaultAccess,
		)

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
		m.Editing = false

	case "enter":
		if m.VaultName == "" {
			return m, nil
		}

		m.ActionState.Start()

		name := m.VaultName

		if m.Editing {
			if m.Cursor >= len(m.VaultAccessList) {
				m.ActionState.Fail(
					errors.New("invalid vault cursor"),
				)
				return m, nil
			}

			selectedVaultAccess := m.VaultAccessList[m.Cursor]

			return m, func() tea.Msg {
				vault := selectedVaultAccess.Vault
				vault.Name = name

				err := m.VaultServiceTUI.UpdateVault(
					&vault,
				)

				return VaultUpdatedMsg{
					VaultAccess: &entity.VaultAccess{
						ID:      selectedVaultAccess.ID,
						VaultID: selectedVaultAccess.VaultID,
						Vault:   vault,
					},
					Err: err,
				}
			}
		}

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
