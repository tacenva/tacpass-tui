package detail

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/tacenva/tacpass-tui/internal/styles"
	"github.com/tacenva/tacpass-tui/internal/tui/components"
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
		"Access Control",
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
		return m.vaultTUI.View()

	case 1:
		// return m.accessControlTUI.View()

	case 2:
		return styles.Muted.Render("Setting")
	}

	return ""
}

func (m Model) Breadcrumb() []string {
	breadcrumbItems := []string{
		"Detail",
	}
	if m.Focus == FocusContent {
		switch m.SidebarCursor {
		case 0:
			items := m.vaultTUI.BreadcrumbItems()

			breadcrumbItems = append(
				breadcrumbItems,
				items...,
			)

		case 1:
			// ...

		case 2:
			// ...
		}
	}

	return breadcrumbItems
}

func (m Model) Navigation() string {
	var content string

	switch m.Focus {
	case FocusSidebar:
		content = components.CrudNavigation(
			m.Width,
			false,
		)

	case FocusContent:
		switch m.SidebarCursor {
		case 0:
			content = m.vaultTUI.Navigation()

		case 1:
			content = styles.NavigationItems(
				m.Width,
				styles.Key("↑↓", "Navigate"),
				styles.Key("↵", "Select"),
				styles.Key("Esc", "Back"),
				styles.Key("q", "Quit"),
			)

		case 2:
			content = styles.NavigationItems(
				m.Width,
				styles.Key("Esc", "Back"),
				styles.Key("q", "Quit"),
			)
		}

	case FocusNewVault:
		content = styles.NavigationItems(
			m.Width,
			styles.Key("↵", "Save"),
			styles.Key("Esc", "Back"),
		)
	}

	return styles.Navigation.
		Width(m.Width).
		Render(content)
}
