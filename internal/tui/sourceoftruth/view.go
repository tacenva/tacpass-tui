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

	table := components.Table{
		Columns: []components.TableColumn{
			{
				Title: "Hostname",
				Width: 30,
			},
			{
				Title: "Address",
				Width: 50,
			},
		},
		Cursor:   m.Cursor,
		AddLabel: "New Node",
		OnFocus:  true,
	}

	for _, sot := range m.SoTList {
		table.Rows = append(
			table.Rows,
			components.TableRow{
				Values: []string{
					sot.Hostname,
					sot.Address,
				},
			},
		)
	}

	return table.View()
}

func (m Model) viewError() string {
	if m.ScreenState.Error == nil {
		return ""
	}

	return styles.Error.Render(
		"[ERROR] " + m.ScreenState.Error.Error(),
	)
}

func (m Model) viewNavigation() string {
	var navigation string

	switch {
	case m.Form.Active:
		navigation = m.Form.Navigation(m.Width)

	case m.Detail.Active:
		navigation = m.Detail.Navigation()

	default:
		var extraItems []string
		if m.Cursor < len(m.SoTList) {
			extraItems = components.DeleteEditItems
		}
		content := components.Navigation(
			m.Width,
			extraItems...,
		)

		navigation = styles.Navigation.
			Width(m.Width).
			Render(content)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.viewError(),
		navigation,
	)
}
