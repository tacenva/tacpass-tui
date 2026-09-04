package styles

import "github.com/charmbracelet/lipgloss"

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1)

	Header = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1)

	Selected = lipgloss.NewStyle().
			Bold(true).
			PaddingLeft(1)

	Normal = lipgloss.NewStyle().
		PaddingLeft(1)

	Muted = lipgloss.NewStyle().
		Faint(true)

	Error = lipgloss.NewStyle().
		Bold(true)

	Border = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	Sidebar = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	Content = lipgloss.NewStyle().
		Padding(1, 2)
)
