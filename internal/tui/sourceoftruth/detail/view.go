package detail

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/tacenva/tacpass-tui/internal/styles"
)

func (m Model) View() string {
	sidebar := m.viewSidebar()
	content := m.viewContent()

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		content,
	)
}

func (m Model) viewSidebar() string {
	items := []string{
		"Vault",
		"Permission",
		"Setting",
	}

	var rows []string

	for i, item := range items {
		prefix := "  "

		if i == m.SidebarCursor && m.Focus == FocusSidebar {
			prefix = "> "
		}

		row := prefix + item

		if i == m.SidebarCursor {
			row = styles.Selected.Render(row)
		} else {
			row = styles.Normal.Render(row)
		}

		rows = append(rows, row)
	}

	return styles.Sidebar.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			rows...,
		),
	)
}

func (m Model) viewContent() string {
	switch m.SidebarCursor {

	case 0:
		return m.viewVault()

	case 1:
		return styles.Muted.Render("Permission")

	case 2:
		return styles.Muted.Render("Setting")
	}

	return ""
}

func (m Model) viewVault() string {
	var rows []string

	for i, vault := range m.Vaults {
		prefix := "  "

		if i == m.Cursor && m.Focus == FocusContent {
			prefix = "> "
		}

		row := prefix + vault.Name

		if i == m.Cursor && m.Focus == FocusContent {
			row = styles.Selected.Render(row)
		} else {
			row = styles.Normal.Render(row)
		}

		rows = append(rows, row)
	}

	if len(rows) == 0 {
		rows = append(
			rows,
			styles.Muted.Render("No vault found."),
		)
	}

	addVault := "+ Add Vault"

	if m.Cursor == len(m.Vaults) {
		addVault = styles.Selected.Render("> Add Vault")
	} else {
		addVault = styles.Muted.Render(addVault)
	}

	rows = append(rows, "", addVault)

	return styles.MainContent.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			rows...,
		))
}

func (m Model) Breadcrumb() []string {
	// if m.Detail.Active {
	// 	items = append(items, m.Detail.BreadcrumbItems()...)
	// }
	return []string{"Detail"}
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
