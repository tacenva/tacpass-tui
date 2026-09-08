package accesscontrol

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/tacenva/tacpass-tui/internal/styles"
	"github.com/tacenva/tacpass-tui/internal/tui/components"
)

func (m Model) View() string {
	switch m.Focus {
	case FocusForm:
		return m.viewForm()

	case FocusUsers:
		return m.UserList.View()

	default:
		return m.viewContent()
	}
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

func (m Model) viewForm() string {
	title := "New Access Control"

	if m.FormMode == FormUpdate {
		title = "Update Access Control"
	}

	name := m.formName()

	if name == "" {
		name = "Enter access control name"
	}

	privilege := string(m.FormPrivilege)

	form := lipgloss.JoinVertical(
		lipgloss.Left,

		styles.Normal.Render(title),

		"",

		components.FormField(
			"Name",
			name,
			m.FormCursor == 0,
		),

		"",

		components.FormField(
			"Privilege",
			privilege,
			m.FormCursor == 1,
		),

		"",
		styles.Muted.Render(
			fmt.Sprintf("Use ←→ to change privilege"),
		),
	)

	return lipgloss.NewStyle().
		PaddingLeft(4).
		Render(form)
}

func (m Model) BreadcrumbItems() []string {
	items := []string{"Access Control"}

	if m.Focus == FocusForm {
		switch m.FormMode {
		case FormNew:
			items = append(items, "New")

		case FormUpdate:
			items = append(items, "Update")
		}
	}

	if m.Focus == FocusUsers {
		items = append(
			items,
			m.UserList.BreadcrumbItems()...,
		)
	}

	return items
}

func (m Model) Navigation() string {
	switch m.Focus {
	case FocusForm:
		return components.FormNavigation(
			m.Width,
			"↵",
		)

	case FocusUsers:
		return m.UserList.Navigation()

	default:
		return styles.NavigationItems(
			m.Width,
			styles.Key("↑↓", "Navigate"),
			styles.Key("↵", "Select"),
			styles.Key("e", "Edit"),
			styles.Key("Esc", "Back"),
			styles.Key("q", "Quit"),
		)
	}
}
