package accesscontrol

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/tacenva/tacpass-core/entity"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true)

	mutedStyle = lipgloss.NewStyle().
			Faint(true)

	errorStyle = lipgloss.NewStyle().
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Bold(true)
)

func (m *Model) View() string {
	if !m.Active {
		return ""
	}

	var content string

	switch m.Page {
	case PageList:
		content = m.viewList()

	case PageUsers:
		content = m.viewUsers()

	case PageCreate:
		content = m.viewCreate()
	}

	return titleStyle.Render("Access Control") +
		"\n\n" +
		content +
		"\n\n" +
		m.viewFooter()
}

func (m *Model) viewList() string {
	if len(m.Permissions) == 0 {
		return "No access control found."
	}

	var b strings.Builder

	b.WriteString(
		fmt.Sprintf(
			"%-32s  %-10s  %s",
			"NAME",
			"PRIVILEGE",
			"STATUS",
		),
	)

	b.WriteString("\n")
	b.WriteString(strings.Repeat("-", 58))
	b.WriteString("\n")

	for i, permission := range m.Permissions {
		cursor := "  "

		if i == m.Cursor {
			cursor = "> "
		}

		status := "active"

		if permission.Revoked {
			status = "revoked"
		}

		name := permission.ID

		line := fmt.Sprintf(
			"%s%-32s  %-10s  %s",
			cursor,
			name,
			permission.Privilege,
			status,
		)

		if i == m.Cursor {
			line = selectedStyle.Render(line)
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	return strings.TrimRight(
		b.String(),
		"\n",
	)
}

func (m *Model) viewUsers() string {
	if m.SelectedPermission == nil {
		return "Permission not found."
	}

	var b strings.Builder

	b.WriteString(
		fmt.Sprintf(
			"Permission: %s",
			m.SelectedPermission.ID,
		),
	)

	b.WriteString("\n\n")

	if len(m.Users) == 0 {
		b.WriteString("No users found.")
		return b.String()
	}

	b.WriteString(
		fmt.Sprintf(
			"%-22s  %-10s  %s",
			"ID",
			"STATUS",
			"HOSTNAME",
		),
	)

	b.WriteString("\n")
	b.WriteString(strings.Repeat("-", 55))
	b.WriteString("\n")

	for i, user := range m.Users {
		cursor := "  "

		if i == m.UserCursor {
			cursor = "> "
		}

		line := fmt.Sprintf(
			"%s%-22s  %-10s  %s",
			cursor,
			user.ID,
			user.Status,
			user.Hostname,
		)

		if i == m.UserCursor {
			line = selectedStyle.Render(line)
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	return strings.TrimRight(
		b.String(),
		"\n",
	)
}

func (m *Model) viewCreate() string {
	var b strings.Builder

	b.WriteString("Create Access Control")
	b.WriteString("\n\n")

	privileges := []struct {
		key       string
		privilege entity.Privilege
	}{
		{
			key:       "1",
			privilege: entity.PrivilegeAdmin,
		},
		{
			key:       "2",
			privilege: entity.PrivilegeWrite,
		},
		{
			key:       "3",
			privilege: entity.PrivilegeRead,
		},
	}

	for _, item := range privileges {
		cursor := "  "

		if item.privilege == m.CreatePrivilege {
			cursor = "> "
		}

		line := fmt.Sprintf(
			"%s[%s] %s",
			cursor,
			item.key,
			item.privilege,
		)

		if item.privilege == m.CreatePrivilege {
			line = selectedStyle.Render(line)
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	return strings.TrimRight(
		b.String(),
		"\n",
	)
}

func (m *Model) viewFooter() string {
	var footer string

	switch m.Page {
	case PageList:
		footer = "↑/k ↓/j navigate • enter users • n new • r refresh • esc back"

	case PageUsers:
		footer = "↑/k ↓/j navigate • a approve • x revoke • r refresh • esc back"

	case PageCreate:
		footer = "1 admin • 2 write • 3 read • enter create • esc back"
	}

	var status string

	if m.Error != nil {
		status = errorStyle.Render(
			"Error: " + m.Error.Error(),
		)
	} else if m.Message != "" {
		status = successStyle.Render(
			m.Message,
		)
	}

	if status != "" {
		footer += "\n\n" + status
	}

	return mutedStyle.Render(footer)
}
