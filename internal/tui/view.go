package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tacenva/tacpass-tui/internal/styles"
)

func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return ""
	}

	switch m.Screen {

	case ScreenLogin:
		return m.viewLogin()

	case ScreenSourceOfTruth:
		return m.SourceOfTruth.View()

	default:
		return ""
	}
}

func (m Model) viewLogin() string {
	title := styles.LoginTitle.Render("TACPASS")

	subtitle := styles.Help.Render(
		"Enter your master password to unlock",
	)

	masked := strings.Repeat(
		"•",
		len([]rune(m.Input)),
	)

	input := styles.LoginInput.
		Width(40).
		Render(masked + "█")

	errorText := ""

	if m.ErrorMessage != "" {
		errorText = styles.Error.Render(m.ErrorMessage)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		subtitle,
		"",
		input,
		"",
		errorText,
		"",
		styles.Help.Render("Enter unlock • Esc quit"),
	)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
