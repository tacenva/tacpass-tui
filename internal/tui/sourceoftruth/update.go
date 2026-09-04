package sourceoftruth

import tea "github.com/charmbracelet/bubbletea"

func InitCmd() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Vault.Collections)-1 {
				m.Cursor++
			}
		}
	}

	return m, nil
}
