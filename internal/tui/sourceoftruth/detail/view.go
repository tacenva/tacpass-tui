package detail

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/tacenva/tacpass-tui/internal/styles"
	"github.com/tacenva/tacpass-tui/internal/tui/components"
)

func (m Model) View() string {
	sidebar := m.viewSidebar()
	var content string
	if m.ScreenState.Loading {
		content = m.viewLoading()
	} else {
		content = m.viewContent()
	}

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

		if i == m.SidebarCursor &&
			m.Focus == FocusSidebar {
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
		return m.accessControlTUI.View()

	case 2:
		return styles.Muted.Render("Setting")
	}

	return ""
}

func (m Model) viewLoading() string {
	return lipgloss.NewStyle().
		PaddingLeft(4).
		Render(
			styles.Normal.Render("Loading..."),
		)
}

func (m Model) Breadcrumb() []string {
	breadcrumbItems := []string{
		"Detail",
	}

	if m.Focus != FocusContent {
		return breadcrumbItems
	}

	switch m.SidebarCursor {
	case 0:
		breadcrumbItems = append(
			breadcrumbItems,
			m.vaultTUI.BreadcrumbItems()...,
		)

	case 1:
		breadcrumbItems = append(
			breadcrumbItems,
			m.accessControlTUI.BreadcrumbItems()...,
		)

	case 2:
		breadcrumbItems = append(
			breadcrumbItems,
			"Setting",
		)
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
			content = m.accessControlTUI.Navigation()

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
