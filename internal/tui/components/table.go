package components

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/tacenva/tacpass-tui/internal/styles"
)

type TableColumn struct {
	Title string
	Width int
}

type TableRow struct {
	Values []string
}

type Table struct {
	Columns  []TableColumn
	Rows     []TableRow
	Cursor   int
	AddLabel string

	OnFocus bool
}

func (t Table) View() string {
	var items []string

	items = append(
		items,
		t.header(),
	)

	for i, row := range t.Rows {
		items = append(
			items,
			t.row(i, row),
		)
	}

	if len(t.Rows) == 0 {
		items = append(
			items,
			styles.Muted.Render("  No data found."),
		)
	}

	if t.AddLabel != "" {
		items = append(
			items,
			t.addRow(),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		items...,
	)
}

func (t Table) header() string {
	var columns []string

	for _, column := range t.Columns {
		columns = append(
			columns,
			lipgloss.NewStyle().
				Width(column.Width).
				Render(
					styles.Muted.Render(column.Title),
				),
		)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		columns...,
	)
}

func (t Table) row(
	index int,
	row TableRow,
) string {
	selected := t.OnFocus && index == t.Cursor

	prefix := ""
	if t.OnFocus {
		prefix = "  "
	}

	if selected {
		prefix = "> "
	}

	var columns []string

	for i, value := range row.Values {
		if i == 0 {
			value = prefix + value
		}

		if selected {
			value = styles.Selected.Render(value)
		} else {
			value = styles.Normal.Render(value)
		}

		width := 0

		if i < len(t.Columns) {
			width = t.Columns[i].Width
		}

		columns = append(
			columns,
			lipgloss.NewStyle().
				Width(width).
				Render(value),
		)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		columns...,
	)
}

func (t Table) addRow() string {
	selected := t.OnFocus && t.Cursor == len(t.Rows)

	prefix := "+ "

	if selected {
		prefix = "> "
	}

	value := prefix + t.AddLabel

	if selected {
		return styles.Selected.
			PaddingTop(1).
			Render(value)
	}

	return styles.Normal.
		PaddingTop(1).
		Render(value)
}
