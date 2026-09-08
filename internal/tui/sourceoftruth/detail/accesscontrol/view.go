package accesscontrol

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/tacenva/tacpass-tui/internal/styles"
	"github.com/tacenva/tacpass-tui/internal/tui/components"
)

func (m Model) View() string {
	return m.viewContent()
}

func (m Model) viewContent() string {
	table := components.Table{
		Columns: []components.TableColumn{
			{
				Title: "Name",
				Width: 25,
			},
			{
				Title: "Privilege",
				Width: 20,
			},
			{
				Title: "Status",
				Width: 20,
			},
		},
		Cursor:   m.Cursor,
		OnFocus:  m.Focus == FocusContent,
		AddLabel: "New Access Control",
	}

	for _, permission := range m.Permissions {
		status := "Active"

		if permission.Revoked {
			status = "Revoked"
		}

		table.Rows = append(
			table.Rows,
			components.TableRow{
				Values: []string{
					permission.Name,
					string(permission.Privilege),
					status,
				},
			},
		)
	}

	return lipgloss.NewStyle().
		PaddingLeft(4).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				styles.Normal.Render("Access Control"),
				"",
				table.View(),
			),
		)
}

func (m Model) BreadcrumbItems() []string {
	return []string{
		"Access Control",
	}
}

func (m Model) Navigation() string {
	return styles.NavigationItems(
		m.Width,
		styles.Key("↑↓", "Navigate"),
		styles.Key("↵", "Select"),
		styles.Key("Esc", "Back"),
		styles.Key("q", "Quit"),
	)
}
