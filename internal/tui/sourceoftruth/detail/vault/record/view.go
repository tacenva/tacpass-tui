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
	var rows []components.TableRow

	for i, record := range m.Records {
		password := "••••••••"

		if m.ShowPassword &&
			i == m.Cursor {
			password = record.Password
		}

		endpoint := strings.TrimPrefix(
			record.Endpoint,
			"https://",
		)

		endpoint = strings.TrimPrefix(
			endpoint,
			"http://",
		)

		rows = append(
			rows,
			components.TableRow{
				Values: []string{
					record.Name,
					password,
					endpoint,
					record.ExpiredAt.Format("02 Jan 2006"),
				},
			},
		)
	}

	table := components.Table{
		Columns: []components.TableColumn{
			{
				Title: "Name",
				Width: 30,
			},
			{
				Title: "Password",
				Width: 25,
			},
			{
				Title: "Endpoint",
				Width: 40,
			},
			{
				Title: "Expired",
				Width: 20,
			},
		},
		Rows:     rows,
		Cursor:   m.Cursor,
		AddLabel: "New Credential",
		OnFocus:  true,
	}

	return styles.MainContent.Render(
		table.View(),
	)
}

func (m Model) viewEdit(title string) string {
	expiredAt := m.EditExpiredAt

	if expiredAt == "" {
		expiredAt = "YYYY-MM-DD"
	}

	rows := []string{
		styles.Title.Render(title),
		"",
		components.FormField(
			"Name",
			m.EditName,
			m.EditField == FieldName,
		),
		"",
		components.FormField(
			"Endpoint",
			m.EditEndpoint,
			m.EditField == FieldEndpoint,
		),
		"",
		components.FormField(
			"Password",
			m.EditPassword,
			m.EditField == FieldPassword,
		),
		"",
		components.FormField(
			"Expired At",
			expiredAt,
			m.EditField == FieldExpiredAt,
		),
	}

	return styles.MainContent.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			rows...,
		),
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
		content = components.CrudNavigation(m.Width, true)
	}

	return styles.Navigation.
		Width(m.Width).
		Render(content)
}
