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

	rows := []string{
		styles.Title.Render(title),
		"",
		m.viewField("Hostname", m.Hostname, FieldHostname),
		"",
		m.viewField("Address", m.Address, FieldAddress),
		"",
		m.viewField("Public Key", m.PublicKey, FieldPublicKey),
		"",
		m.viewField(
			"Private Key",
			m.privateKeyView(),
			FieldPrivateKey,
		),
	}

	return styles.MainContent.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			rows...,
		),
	)
}

func (m Model) viewField(
	label string,
	value string,
	field Field,
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

// func (m Model) privateKeyView() string {
// 	if m.PrivateKey == "" {
// 		return ""
// 	}

// 	runes := []rune(m.PrivateKey)
// 	result := make([]rune, len(runes))

// 	for i := range runes {
// 		result[i] = '•'
// 	}

// 	return string(result)
// }

func (m Model) privateKeyView() string {
	return m.PrivateKey
}

func (m Model) Navigation(width int) string {
	content := components.FormNavigation(m.Width, "Ctrl+S")

	return styles.Navigation.
		Width(width).
		Render(content)
}
