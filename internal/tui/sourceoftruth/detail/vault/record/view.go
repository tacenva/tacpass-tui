package vaultrecord

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tacenva/tacpass-tui/internal/styles"
	"github.com/tacenva/tacpass-tui/internal/tui/components"
)

func (m Model) View() string {
	switch m.Focus {
	case FocusEdit:
		return m.viewEdit("Edit Credential")

	case FocusNew:
		return m.viewEdit("New Credential")

	default:
		return m.viewContent()
	}
}

func (m Model) viewContent() string {
	var rows []string

	nameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true)

	passwordStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10"))

	endpointStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))

	expiredStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))

	contentWidth := m.Width - 4
	if contentWidth < 40 {
		contentWidth = 40
	}

	for i, record := range m.Records {
		selected := i == m.Cursor

		prefix := "  "
		if selected {
			prefix = "> "
		}

		password := "••••••••"
		if selected && m.ShowPassword {
			password = record.Password
		}

		nameWidth := lipgloss.Width(record.Name)
		passwordWidth := lipgloss.Width(password)

		gap := contentWidth - nameWidth - passwordWidth
		if gap < 2 {
			gap = 2
		}

		header := prefix +
			nameStyle.Render(record.Name) +
			strings.Repeat(" ", gap) +
			passwordStyle.Render(password)

		endpoint := strings.TrimPrefix(record.Endpoint, "https://")
		endpoint = strings.TrimPrefix(endpoint, "http://")

		detail := "  " +
			endpointStyle.Render(endpoint) +
			"    " +
			expiredStyle.Render(
				"Exp: "+record.ExpiredAt.Format("02 Jan 2006"),
			)

		rows = append(
			rows,
			header,
			detail,
			"",
		)
	}

	if len(m.Records) == 0 {
		rows = append(
			rows,
			styles.Muted.Render("No credential found."),
			"",
		)
	}

	newCredential := "+ New Credential"

	if m.Cursor == len(m.Records) {
		newCredential = "> Add Credential"
	}

	rows = append(
		rows,
		"",
		styles.Selected.Render(newCredential),
	)

	return styles.MainContent.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			rows...,
		),
	)
}

func (m Model) viewEdit(title string) string {
	rows := []string{
		styles.Title.Render(title),
		"",
		m.viewEditField(
			"Name",
			m.EditName,
			FieldName,
		),
		"",
		m.viewEditField(
			"Endpoint",
			m.EditEndpoint,
			FieldEndpoint,
		),
		"",
		m.viewEditField(
			"Password",
			m.EditPassword,
			FieldPassword,
		),
		"",
		m.viewEditField(
			"Expired At",
			m.EditExpiredAt,
			FieldExpiredAt,
		),
	}

	return styles.MainContent.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			rows...,
		),
	)
}

func (m Model) viewEditField(
	label string,
	value string,
	field EditField,
) string {
	prefix := "  "

	if m.EditField == field {
		prefix = "> "
	}

	labelView := styles.Normal.Render(
		prefix + label,
	)

	valueView := styles.Normal.Render(
		"    " + value,
	)

	if m.EditField == field {
		valueView = styles.Selected.Render(
			"    " + value + "_",
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		labelView,
		valueView,
	)
}

func (m Model) BreadcrumbItems() []string {
	items := []string{
		"Records",
	}
	switch m.Focus {
	case FocusEdit:
		items = append(items, "Edit")
	case FocusNew:
		items = append(items, "New")
	}
	return items
}

func (m Model) Navigation() string {
	var content string
	if m.Focus == FocusEdit || m.Focus == FocusNew {
		content = components.FormNavigation(m.Width, "Ctrl+S")
	} else {
		content = styles.NavigationItems(
			m.Width,
			styles.Key("↑↓", "Navigate"),
			styles.Key("↵", "Edit"),
			styles.Key("p", "Show Password"),
			styles.Key("Esc", "Back"),
			styles.Key("q", "Quit"),
		)
	}

	return styles.Navigation.
		Width(m.Width).
		Render(content)
}
