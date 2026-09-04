package sourceoftruth

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tacenva/tacpass-tui/internal/styles"
)

func (m Model) View() string {
	var content strings.Builder

	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.Title.Render("TACPASS"),
		"  ",
		styles.Muted.Render(m.Vault.Name),
	)

	content.WriteString(header)
	content.WriteString("\n\n")

	for i, collection := range m.Vault.Collections {
		if i == m.Cursor {
			content.WriteString(
				styles.Selected.Render("> " + collection.Name),
			)
		} else {
			content.WriteString(
				styles.Normal.Render("  " + collection.Name),
			)
		}

		content.WriteString("\n")

		content.WriteString(
			styles.Muted.Render(
				fmt.Sprintf("    %d items", len(collection.Items)),
			),
		)

		content.WriteString("\n\n")
	}

	content.WriteString(
		styles.Muted.Render("j/k navigate • enter select • q quit"),
	)

	return content.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
