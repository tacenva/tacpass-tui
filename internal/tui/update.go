package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		m.SourceOfTruth.Width = msg.Width
		m.SourceOfTruth.Height = msg.Height

		return m, nil
	}

	switch m.Screen {

	case ScreenLogin:
		keyMsg, ok := msg.(tea.KeyMsg)
		if !ok {
			return m, nil
		}

		return m.updateLogin(keyMsg)

	case ScreenSourceOfTruth:
		var cmd tea.Cmd

		m.SourceOfTruth, cmd = m.SourceOfTruth.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m Model) updateLogin(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	case "ctrl+c", "esc":
		return m, tea.Quit

	case "enter":
		masterKey := m.Input
		if masterKey == "" {
			m.ErrorMessage = "Invalid master password"
			return m, nil
		}

		err := m.sotService.Access(m.Input)
		if err != nil {
			m.ErrorMessage = "Invalid master password"
			m.Input = ""
			return m, nil
		}

		m.Input = ""
		m.ErrorMessage = ""

		if err := m.SourceOfTruth.Load(masterKey); err != nil {
			m.ErrorMessage = err.Error()
			return m, nil
		}

		m.Screen = ScreenSourceOfTruth

		return m, nil

	case "backspace":
		runes := []rune(m.Input)

		if len(runes) > 0 {
			m.Input = string(runes[:len(runes)-1])
		}

	default:
		if len(msg.Runes) > 0 {
			m.Input += string(msg.Runes)
			m.ErrorMessage = ""
		}
	}

	return m, nil
}
