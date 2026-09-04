package sourceoftruth

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tacenva/tacpass-tui/internal/styles"
)

func (m Model) View() string {
	header := m.viewHeader()
	breadcrumb := m.viewBreadcrumb()
	content := m.viewContent()
	footer := m.viewNavigation()

	contentHeight := m.Height -
		lipgloss.Height(header) -
		lipgloss.Height(breadcrumb) -
		lipgloss.Height(footer) -
		1

	if contentHeight < 1 {
		contentHeight = 1
	}

	content = lipgloss.NewStyle().
		Height(contentHeight).
		PaddingLeft(4).
		Render(content)

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		breadcrumb,
		"",
		content,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		body,
		footer,
	)
}

func (m Model) viewHeader() string {
	return styles.Header.Render("TACPASS")
}

func (m Model) viewBreadcrumb() string {
	items := []string{
		"Source of Truth",
	}
	if m.Detail.Active {
		items = append(items, m.Detail.Breadcrumb()...)
	}
	return styles.BreadcrumbItems(items...)
}

func (m Model) viewContent() string {
	if m.Detail.Active {
		return m.Detail.View()
	}

	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().
			Width(30).
			Render(styles.Muted.Render("Hostname")),
		lipgloss.NewStyle().
			Width(50).
			Render(styles.Muted.Render("Address")),
	)

	items := []string{header}

	for i, sot := range m.SoTList {
		hostname := sot.Hostname
		address := sot.Address

		if i == m.Cursor {
			hostname = styles.Selected.Render("> " + hostname)
			address = styles.Selected.Render(address)
		} else {
			hostname = styles.Normal.Render("  " + hostname)
			address = styles.Normal.Render(address)
		}

		item := lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().
				Width(30).
				Render(hostname),
			lipgloss.NewStyle().
				Width(50).
				Render(address),
		)

		items = append(items, item)
	}

	if len(m.SoTList) == 0 {
		items = append(
			items,
			styles.Muted.Render("  No Source of Truth found."),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		items...,
	)
}

func (m Model) viewNavigation() string {
	if m.Detail.Active {
		return m.Detail.Navigation()
	}

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
