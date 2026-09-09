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
		return m.viewNewVault()
	}

	return m.viewContent()
}

func (m Model) viewContent() string {
	var rows []components.TableRow

	for _, vaultAccess := range m.VaultAccessList {
		rows = append(
			rows,
			components.TableRow{
				Values: []string{
					vaultAccess.Vault.Name,
				},
			},
		)
	}

	table := components.Table{
		Columns: []components.TableColumn{
			{
				Title: "Vault Name",
				Width: 50,
			},
		},
		Rows:     rows,
		Cursor:   m.Cursor,
		AddLabel: "New Vault",
		OnFocus:  m.Focus == FocusContent,
	}

	title := "Vault"
	if m.needSync {
		title = title + " outdated"
	}

	return styles.MainContent.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			styles.Title.Render(title),
			"",
			table.View(),
		),
	)
}

func (m Model) viewNewVault() string {
	title := "Add Vault"

	if m.Editing {
		title = "Edit Vault"
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.Title.Render(title),
		"",
		styles.Normal.Render("Name"),
		"> "+m.VaultName+"_",
	)

	if m.ActionState.Loading {
		action := "Creating..."

		if m.Editing {
			action = "Updating..."
		}

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			"",
			styles.Muted.Render(action),
		)
	}

	if m.ActionState.Error != nil {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			"",
			styles.Error.Render(
				"[ERROR] "+m.ActionState.Error.Error(),
			),
		)
	}

	return styles.MainContent.Render(content)
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
		var extraItems []string
		if m.Cursor < len(m.VaultAccessList) {
			extraItems = components.DeleteEditItems
		}

		extraItems = append(extraItems, styles.Key("s", "Sync"))
		content = components.Navigation(
			m.Width,
			extraItems...,
		)

	case FocusNewVault:
		content = components.FormNavigation(
			m.Width,
			"↵",
		)
	}

	return content
}
