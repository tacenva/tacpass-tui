package accesscontrol

import (
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

	case FocusCreatedKey:
		return m.viewCreatedKey()

	default:
		return m.viewContent()
	}
}

func (m Model) viewCreatedKey() string {
	labelStyle := lipgloss.NewStyle().
		Width(12)

	valueStyle := lipgloss.NewStyle().
		Width(80)

	publicKey := lipgloss.JoinHorizontal(
		lipgloss.Left,
		labelStyle.Render("Public Key"),
		valueStyle.Render(m.CreatedPublicKey),
	)

	privateKey := lipgloss.JoinHorizontal(
		lipgloss.Left,
		labelStyle.Render("Private Key"),
		valueStyle.Render(m.CreatedPrivateKey),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		publicKey,
		"",
		privateKey,
	)

	return lipgloss.NewStyle().
		PaddingLeft(4).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				styles.Normal.Render("Access Control Created"),
				"",
				content,
				"",
				styles.Muted.Render(
					"Save both keys securely. The private key will not be shown again.",
				),
			),
		)
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

	return styles.MainContent.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			styles.Title.Render("Access Control"),
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

	rows := []string{
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
			"Use ←→ to change privilege",
		),
	}

	if actionState := m.viewActionState(); actionState != "" {
		rows = append(rows, "", actionState)
	}

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		rows...,
	)

	return lipgloss.NewStyle().
		PaddingLeft(4).
		Render(form)
}

func (m Model) viewActionState() string {
	return components.AsyncStateView(
		m.ActionState,
		"Saving...",
	)
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
	var content string
	switch m.Focus {
	case FocusForm:
		content = components.FormNavigation(
			m.Width,
			"↵",
		)

	case FocusUsers:
		content = m.UserList.Navigation()

	case FocusCreatedKey:
		content = styles.NavigationItems(
			m.Width,
			styles.Key("↑↓", "Navigate"),
			styles.Key("Esc", "Back"),
		)

	default:
		var extraItems []string
		if m.Cursor < len(m.Permissions) {
			extraItems = append(
				[]string{styles.Key("r", "Revoke")},
				components.DeleteEditItems...,
			)
		}
		content = components.Navigation(
			m.Width,
			extraItems...,
		)
	}
	return content
}
