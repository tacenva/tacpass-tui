package vault

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tacenva/tacpass-tui/internal/styles"
	"github.com/tacenva/tacpass-tui/internal/tui/components"
)

func (m Model) View() string {
	if m.VaultRecordTUI.Active {
		return m.VaultRecordTUI.View()
	}

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
	if m.ErrorMessage != "" {
		rows = append(
			rows,
			styles.Error.Render(m.ErrorMessage),
			"",
		)
	}

	for i, vaultaccess := range m.VaultAccessList {
		selected := i == m.Cursor && m.Focus == FocusContent

		prefix := ""
		if m.Focus == FocusContent {
			prefix = "  "
		}

		if selected {
			prefix = "> "
		}

		row := prefix + vaultaccess.Vault.Name

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

	if m.Cursor == len(m.VaultAccessList) && m.Focus == FocusContent {
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
func (m Model) BreadcrumbItems() []string {
	items := []string{
		"Vault",
	}

	if m.VaultRecordTUI.Active {
		recordBreadcrumbItems := m.VaultRecordTUI.BreadcrumbItems()
		items = append(items, recordBreadcrumbItems...)
	}

	if m.Focus == FocusNewVault {
		items = append(items, "New")
	}

	return items
}

func (m Model) Navigation() string {
	if m.VaultRecordTUI.Active {
		return m.VaultRecordTUI.Navigation()
	}

	var content string

	switch m.Focus {
	case FocusContent:
		content = components.CrudNavigation(
			m.Width,
			m.Cursor < len(m.VaultAccessList),
		)

	case FocusNewVault:
		content = components.FormNavigation(
			m.Width,
			"↵",
		)
	}

	return content
}
