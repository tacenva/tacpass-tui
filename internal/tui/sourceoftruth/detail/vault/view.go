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

	var rows []components.TableRow

	for _, vaultaccess := range m.VaultAccessList {
		rows = append(
			rows,
			components.TableRow{
				Values: []string{
					vaultaccess.Vault.Name,
				},
			},
		)
	}

	table := components.Table{
		Columns: []components.TableColumn{
			{
				Title: "Vault",
				Width: 50,
			},
		},
		Rows:     rows,
		Cursor:   m.Cursor,
		AddLabel: "New Vault",
		OnFocus:  m.Focus == FocusContent,
	}

	return styles.MainContent.Render(
		table.View(),
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
