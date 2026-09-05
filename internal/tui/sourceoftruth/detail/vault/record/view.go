package vault

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tacenva/tacpass-tui/internal/styles"
)

func (m Model) View() string {
	var rows []string

	for i, record := range m.Records {
		prefix := "  "

		if i == m.Cursor {
			prefix = "> "
		}

		row := prefix + record.Name

		if i == m.Cursor {
			row = styles.Selected.Render(row)
		} else {
			row = styles.Normal.Render(row)
		}

		rows = append(rows, row)
	}

	if len(rows) == 0 {
		rows = append(
			rows,
			styles.Muted.Render("No record found."),
		)
	}

	return styles.MainContent.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			rows...,
		),
	)
}

func (m Model) Breadcrumb() []string {
	return []string{
		"Vault",
		m.Vault.Name,
	}
}

func (m Model) Navigation() string {
	content := styles.NavigationItems(
		m.Width,
		styles.Key("↑↓", "Navigate"),
		styles.Key("↵", "Select"),
		styles.Key("Esc", "Back"),
		styles.Key("q", "Quit"),
	)

	return styles.Navigation.
		Width(m.Width).
		Render(content)
}
