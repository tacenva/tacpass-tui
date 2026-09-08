package components

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/tacenva/tacpass-tui/internal/styles"
)

func CrudNavigation(
	width int,
	canEdit bool,
) string {
	items := []string{
		styles.Key("↑↓", "Navigate"),
		styles.Key("↵", "Select"),
	}

	if canEdit {
		items = append(
			items,
			styles.Key("e", "Edit"),
			styles.Key("Del", "Delete"),
		)
	}

	items = append(
		items,
		styles.Key("Esc", "Back"),
		styles.Key("q", "Quit"),
	)

	content := styles.NavigationItems(
		width,
		items...,
	)

	return lipgloss.NewStyle().
		Width(width).
		Render(content)
}

func FormNavigation(
	width int,
	saveKey string,
) string {
	content := styles.NavigationItems(
		width,
		styles.Key(saveKey, "Save"),
		styles.Key("Esc", "Cancel"),
	)

	return lipgloss.NewStyle().
		Width(width).
		Render(content)
}
