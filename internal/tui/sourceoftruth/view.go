package sourceoftruth

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/tacenva/tacpass-tui/internal/styles"
	"github.com/tacenva/tacpass-tui/internal/tui/components"
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

	if m.Form.Active {
		if m.Form.Editing {
			items = append(items, "Edit")
		} else {
			items = append(items, "New")
		}

		return styles.BreadcrumbItems(items...)
	}

	if m.Detail.Active {
		items = append(
			items,
			m.Detail.Breadcrumb()...,
		)
	}

	return styles.BreadcrumbItems(items...)
}

func (m Model) viewContent() string {
	if m.Form.Active {
		return m.Form.View()
	}

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

	items := []string{
		header,
	}

	for i, sot := range m.SoTList {
		hostname := sot.Hostname
		address := sot.Address

		if i == m.Cursor {
			hostname = styles.Selected.Render(
				"> " + hostname,
			)

			address = styles.Selected.Render(
				address,
			)
		} else {
			hostname = styles.Normal.Render(
				"  " + hostname,
			)

			address = styles.Normal.Render(
				address,
			)
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
			styles.Muted.Render(
				"  No Source of Truth found.",
			),
		)
	}

	newIndex := len(m.SoTList)

	if m.Cursor == newIndex {
		items = append(
			items,
			styles.Selected.PaddingTop(1).Render("> Add Node"),
		)
	} else {
		items = append(
			items,
			styles.Normal.PaddingTop(1).Render("+ New Node"),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		items...,
	)
}

func (m Model) viewNavigation() string {
	if m.Form.Active {
		return m.Form.Navigation(m.Width)
	}

	if m.Detail.Active {
		return m.Detail.Navigation()
	}

	content := components.CrudNavigation(
		m.Width,
		m.Cursor < len(m.SoTList),
	)

	return styles.Navigation.
		Width(m.Width).
		Render(content)
}
