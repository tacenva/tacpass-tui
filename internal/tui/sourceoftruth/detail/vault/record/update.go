package vaultrecord

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch m.Focus {
		case FocusContent:
			return m.updateContent(msg)

		case FocusEdit:
			return m.updateEdit(msg)

		case FocusNew:
			return m.updateNew(msg)
		}
	}

	return m, nil
}

func (m Model) updateContent(msg tea.KeyMsg) (Model, tea.Cmd) {
	newCredentialCursor := len(m.Records)

	switch msg.String() {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
			m.ShowPassword = false
		}

	case "down", "j":
		if m.Cursor < newCredentialCursor {
			m.Cursor++
			m.ShowPassword = false
		}

	case "enter":
		if m.Cursor == newCredentialCursor {
			m.startNew()
			return m, nil
		}

		m.startEdit()

	case "p":
		if m.Selected() != nil {
			m.ShowPassword = !m.ShowPassword
		}

	case "esc", "left", "h":
		m.ShowPassword = false
		m.Focus = Unfocus

		// case "q":
		m.Active = false
	}

	return m, nil
}

func (m Model) updateEdit(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	case "tab", "down":
		m.nextField()

	case "shift+tab", "up":
		m.previousField()

	case "enter", "ctrl+s":
		m.saveEdit()

	case "esc":
		m.cancelForm()

	case "p":
		if m.EditField == FieldPassword {
			m.ShowPassword = !m.ShowPassword
		}

	case "backspace":
		m.removeEditCharacter()

	default:
		if len(msg.Runes) > 0 {
			m.appendEditCharacter(string(msg.Runes))
		}
	}

	return m, nil
}

func (m Model) updateNew(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	case "tab", "down":
		m.nextField()

	case "shift+tab", "up":
		m.previousField()

	case "ctrl+s":
		m.saveNew()

	case "esc":
		m.cancelForm()

	case "backspace":
		m.removeEditCharacter()

	default:
		if len(msg.Runes) > 0 {
			m.appendEditCharacter(string(msg.Runes))
		}
	}

	return m, nil
}

func (m *Model) nextField() {
	switch m.EditField {

	case FieldName:
		m.EditField = FieldEndpoint

	case FieldEndpoint:
		m.EditField = FieldPassword

	case FieldPassword:
		m.EditField = FieldExpiredAt

	case FieldExpiredAt:
		m.EditField = FieldName
	}
}

func (m *Model) previousField() {
	switch m.EditField {

	case FieldName:
		m.EditField = FieldExpiredAt

	case FieldEndpoint:
		m.EditField = FieldName

	case FieldPassword:
		m.EditField = FieldEndpoint

	case FieldExpiredAt:
		m.EditField = FieldPassword
	}
}

func (m *Model) appendEditCharacter(value string) {
	switch m.EditField {

	case FieldName:
		m.EditName += value

	case FieldEndpoint:
		m.EditEndpoint += value

	case FieldPassword:
		m.EditPassword += value

	case FieldExpiredAt:
		m.EditExpiredAt += value
	}
}

func (m *Model) removeEditCharacter() {
	switch m.EditField {

	case FieldName:
		if len(m.EditName) > 0 {
			m.EditName = m.EditName[:len(m.EditName)-1]
		}

	case FieldEndpoint:
		if len(m.EditEndpoint) > 0 {
			m.EditEndpoint = m.EditEndpoint[:len(m.EditEndpoint)-1]
		}

	case FieldPassword:
		if len(m.EditPassword) > 0 {
			m.EditPassword = m.EditPassword[:len(m.EditPassword)-1]
		}

	case FieldExpiredAt:
		if len(m.EditExpiredAt) > 0 {
			m.EditExpiredAt = m.EditExpiredAt[:len(m.EditExpiredAt)-1]
		}
	}
}
