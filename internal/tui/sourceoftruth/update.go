package sourceoftruth

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail"
)

func InitCmd() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.Detail.Active {
		updated, cmd := m.Detail.Update(msg)
		m.Detail = updated
		return m, cmd
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.SoTList)-1 {
				m.Cursor++
			}

		case "enter":
			if len(m.SoTList) == 0 {
				return m, nil
			}
			selectedSoT := m.SoTList[m.Cursor]
			m.Detail = detail.New(m.appDeps, &app.Context{
				SelectedSoT:       &selectedSoT,
				AuthService:       m.authService,
				PermissionService: m.permissionService,
			})
		}
	}

	return m, nil
}
