package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	Border = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	Sidebar = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	Content = lipgloss.NewStyle().
		Padding(1, 2)

	Header = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1)

	Breadcrumb = lipgloss.NewStyle().
			Faint(true).
			Padding(0, 1)

	Navigation = lipgloss.NewStyle().
			Faint(true).
			Padding(0, 1).
			PaddingTop(1)

	Title = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1)

	Selected = lipgloss.NewStyle().
			Bold(true).
			PaddingLeft(1)

	Normal = lipgloss.NewStyle().
		PaddingLeft(1)

	Muted = lipgloss.NewStyle().
		Faint(true).
		PaddingLeft(1)

	Error = lipgloss.NewStyle().
		Bold(true)

	KeyStyle = lipgloss.NewStyle().
			Bold(true)
)

func Key(key, action string) string {
	return KeyStyle.Render("["+key+"]") + " " + action
}

func NavigationItems(width int, items ...string) string {
	if width <= 0 {
		return strings.Join(items, "   ")
	}

	// Navigation punya padding kiri + kanan.
	maxWidth := width - Navigation.GetHorizontalFrameSize()

	if maxWidth <= 0 {
		return strings.Join(items, "   ")
	}

	var lines []string
	var current string

	for _, item := range items {
		if current == "" {
			current = item
			continue
		}

		next := current + "   " + item

		if lipgloss.Width(next) > maxWidth {
			lines = append(lines, current)
			current = item
			continue
		}

		current = next
	}

	if current != "" {
		lines = append(lines, current)
	}

	return strings.Join(lines, "\n")
}
