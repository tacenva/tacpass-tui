package vaultrecord

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/tacenva/tacpass-tui/internal/styles"
)

func (m Model) View() string {
	var rows []string

	for i, record := range m.Records {
		selected := i == m.Cursor

		prefix := "  "
		if selected {
			prefix = "> "
		}

		name := prefix + record.Name

		if selected {
			name = styles.Selected.Render(name)

			password := "********"
			if m.ShowPassword {
				password = record.Password
			}

			detail := lipgloss.JoinVertical(
				lipgloss.Left,
				styles.Normal.Render(fmt.Sprintf("Endpoint : %s", record.Endpoint)),
				styles.Normal.Render(fmt.Sprintf("Password : %s", password)),
				styles.Normal.Render(fmt.Sprintf("Expired  : %s", record.ExpiredAt)),
			)

			rows = append(
				rows,
				name,
				"  "+detail,
				"",
			)

			continue
		}

		rows = append(
			rows,
			styles.Normal.Render(name),
		)
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
		m.SelectedVault.Name,
	}
}

func (m Model) Navigation() string {
	content := styles.NavigationItems(
		m.Width,
		styles.Key("↑↓", "Navigate"),
		styles.Key("↵", "Select"),
		styles.Key("p", "Show Password"),
		styles.Key("Esc", "Back"),
		styles.Key("q", "Quit"),
	)

	return styles.Navigation.
		Width(m.Width).
		Render(content)
}
