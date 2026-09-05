package vaultrecord

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Records)-1 {
				m.Cursor++
			}

		case "p":
			m.ShowPassword = !m.ShowPassword

		case "esc", "left", "h":
			m.ShowPassword = false
			m.Active = false
		}
	}

	return m, nil
}
