package userlist

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
				Title: "Hostname",
				Width: 30,
			},
			{
				Title: "Status",
				Width: 20,
			},
		},
		Cursor:  m.Cursor,
		OnFocus: m.Focus == FocusContent,
	}

	for _, user := range m.Users {
		table.Rows = append(
			table.Rows,
			components.TableRow{
				Values: []string{
					user.Hostname,
					string(user.Status),
				},
			},
		)
	}

	rows := []string{
		styles.Normal.Render("Users"),
		"",
		table.View(),
	}

	if actionState := m.viewActionState(); actionState != "" {
		rows = append(rows, "", actionState)
	}

	return lipgloss.NewStyle().
		PaddingLeft(4).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				rows...,
			),
		)
}

func (m Model) viewActionState() string {
	return components.AsyncStateView(
		m.ActionState,
		"Updating...",
	)
}

func (m Model) BreadcrumbItems() []string {
	return []string{"Users"}
}

func (m Model) Navigation() string {
	return styles.NavigationItems(
		m.Width,
		styles.Key("↑↓", "Navigate"),
		styles.Key("a", "Approve"),
		styles.Key("r", "Revoke"),
		styles.Key("Esc", "Back"),
		styles.Key("q", "Quit"),
	)
}
