package vault

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tacenva/tacpass-tui/internal/styles"
)

func (m Model) View() string {
	if m.Focus == FocusNewVault {
		return styles.MainContent.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				styles.Title.Render("Add Vault"),
				"",
				styles.Normal.Render("Name"),
				"> "+m.VaultName+"_",
			),
		)
	}

	var rows []string

	for i, vault := range m.VaultList {
		selected := i == m.Cursor && m.Focus == FocusContent

		prefix := ""
		if m.Focus == FocusContent {
			prefix = "  "
		}

		if selected {
			prefix = "> "
		}

		row := prefix + vault.Name

		if selected {
			row = styles.Selected.Render(row)
		} else {
			row = styles.Normal.Render(row)
		}

		rows = append(rows, row)
	}

	if len(rows) == 0 {
		rows = append(
			rows,
			styles.Muted.Render("No vault found."),
		)
	}

	addVault := "+ New Vault"

	if m.Cursor == len(m.VaultList) && m.Focus == FocusContent {
		addVault = styles.Selected.Render("> Add Vault")
	} else {
		addVault = styles.Muted.Render(addVault)
	}

	rows = append(rows, "", addVault)

	return styles.MainContent.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			rows...,
		),
	)
}
