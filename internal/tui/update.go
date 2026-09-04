package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		m.SourceOfTruth.Width = msg.Width
		m.SourceOfTruth.Height = msg.Height

		return m, nil

	case tea.KeyMsg:
		switch m.Screen {

		case ScreenLogin:
			return m.updateLogin(msg)

		case ScreenSourceOfTruth:
			var cmd tea.Cmd

			m.SourceOfTruth, cmd = m.SourceOfTruth.Update(msg)

			return m, cmd
		}
	}

	return m, nil
}

func (m Model) updateLogin(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	case "ctrl+c", "esc":
		return m, tea.Quit

	case "enter":
		err := m.sotService.Access(m.Input)
		if err == nil {
			m.Input = ""
			m.ErrorMessage = ""
			m.Screen = ScreenSourceOfTruth

			return m, sourceoftruth.InitCmd()
		}

		m.ErrorMessage = "Invalid master password"
		m.Input = ""

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
