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
		components.FormField(
			"Hostname",
			m.Hostname,
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

	return styles.MainContent.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			rows...,
		),
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
