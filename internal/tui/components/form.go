package components

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/tacenva/tacpass-tui/internal/styles"
)

func FormField(
	label string,
	value string,
	selected bool,
) string {
	prefix := ""

	if selected {
		prefix = "> "
	}

	labelView := styles.Normal.Render(
		prefix + label,
	)

	valueView := styles.Normal.Render(
		value,
	)

	if selected {
		valueView = styles.Selected.Render(
			value + "_",
		)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().
			Width(16).
			Render(labelView),
		valueView,
	)
}
