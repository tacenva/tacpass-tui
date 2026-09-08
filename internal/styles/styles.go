package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const paddingLeft = 4

var (
	Sidebar = lipgloss.NewStyle().
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		PaddingRight(2)

	MainContent = lipgloss.NewStyle().
			PaddingLeft(2)

	Header = lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(paddingLeft).
		PaddingTop(2)

	Breadcrumb = lipgloss.NewStyle().
			Faint(true).
			PaddingLeft(paddingLeft)

	Navigation = lipgloss.NewStyle().
			Faint(true).
			PaddingLeft(paddingLeft).
			PaddingTop(1)

	Selected = lipgloss.NewStyle().
			Bold(true).
			PaddingLeft(0)

	Normal = lipgloss.NewStyle().
		PaddingLeft(0)

	Muted = lipgloss.NewStyle().
		Faint(true)

	Error = lipgloss.NewStyle().
		PaddingLeft(paddingLeft).
		Bold(true)

	KeyStyle = lipgloss.NewStyle().
			Bold(true)

	Title = lipgloss.NewStyle().
		Bold(true).
		Padding(0)
)

func Key(key, action string) string {
	return KeyStyle.Render("["+key+"]") + " " + action
}

func BreadcrumbItems(items ...string) string {
	if len(items) == 0 {
		return ""
	}

	return Breadcrumb.Render(
		strings.Join(items, " / "),
	)
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
