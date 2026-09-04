package sourceoftruth

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/tacpass-tui/internal/entity"
)

type SoTLoadedMsg struct {
	SoTList []entity.SourceOfTruth
}

func InitCmd() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case SoTLoadedMsg:
		m.SoTList = msg.SoTList

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
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
			// nanti buka detail/edit
		}
	}

	return m, nil
}

func (m *Model) Load() error {
	SoTList, err := m.SoTService.List()
	if err != nil {
		return err
	}

	m.SoTList = SoTList

	return nil
}
