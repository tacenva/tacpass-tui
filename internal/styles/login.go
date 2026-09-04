package styles

import "github.com/charmbracelet/lipgloss"

var (
	LoginTitle = lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center).
			Padding(0, 1)

	Help = lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Align(lipgloss.Center)

	LoginInput = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Align(lipgloss.Left)
)
