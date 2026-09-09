package form

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/tacenva/tacpass-tui/internal/styles"
	"github.com/tacenva/tacpass-tui/internal/tui/components"
)

func (m Model) View() string {
	title := "New Source of Truth"

	if m.Editing {
		title = "Edit Source of Truth"
	}

	hostname := m.Hostname

	if hostname == "" && !m.Editing {
		hostname = "Optional, used as alias"
	}

	rows := []string{
		styles.Title.Render(title),
		"",
		components.FormField(
			"Hostname",
			hostname,
			m.EditField == FieldHostname,
		),
		"",
		components.FormField(
			"Address",
			m.Address,
			m.EditField == FieldAddress,
		),
		"",
		components.FormField(
			"Public Key",
			m.PublicKey,
			m.EditField == FieldPublicKey,
		),
		"",
		components.FormField(
			"Private Key",
			m.privateKeyView(),
			m.EditField == FieldPrivateKey,
		),
	}

	if content := m.viewActionState(); content != "" {
		rows = append(rows, "", content)
	}

	return styles.MainContent.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			rows...,
		),
	)
}

func (m Model) viewActionState() string {
	return components.AsyncStateView(
		m.ActionState,
		"Saving...",
	)
}

func (m Model) privateKeyView() string {
	return m.PrivateKey
}

func (m Model) Navigation(width int) string {
	content := components.FormNavigation(
		m.Width,
		"Ctrl+S",
	)

	return styles.Navigation.
		Width(width).
		Render(content)
}
