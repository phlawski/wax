package table

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// Solarized Dark palette
const (
	solBase03 = lipgloss.Color("#002b36")
	solBase02 = lipgloss.Color("#073642")
	solBase01 = lipgloss.Color("#586e75")
	solBase0  = lipgloss.Color("#839496")
	solBase1  = lipgloss.Color("#93a1a1")
	solBlue   = lipgloss.Color("#268bd2")
	solCyan   = lipgloss.Color("#2aa198")
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(solBase03).
			Background(solBlue).
			PaddingLeft(1).
			PaddingRight(1)

	rowStyle = lipgloss.NewStyle().
			Foreground(solBase0).
			Background(solBase03).
			PaddingLeft(1).
			PaddingRight(1)

	altRowStyle = lipgloss.NewStyle().
			Foreground(solBase1).
			Background(lipgloss.Color("#0a4452")).
			PaddingLeft(1).
			PaddingRight(1)

	dimStyle = lipgloss.NewStyle().
			Foreground(solBase01)
)

// Column defines a table column.
type Column struct {
	Title string
	Width int
}

// Table renders a styled terminal table.
type Table struct {
	columns []Column
}

// New creates a Table with the given columns.
func New(columns []Column) *Table {
	return &Table{columns: columns}
}

// Render prints the table with headers and rows.
func (t *Table) Render(rows [][]string) {
	termWidth := terminalWidth()

	// Print header
	header := t.renderRow(headerStyle, termWidth)
	fmt.Println(header)

	// Print rows
	for i, row := range rows {
		style := rowStyle
		if i%2 == 1 {
			style = altRowStyle
		}
		fmt.Println(t.renderRowData(style, row, termWidth))
	}
}

func (t *Table) renderRow(style lipgloss.Style, termWidth int) string {
	cells := make([]string, len(t.columns))
	for i, col := range t.columns {
		cells[i] = style.Width(col.Width).Render(truncate(col.Title, col.Width-2))
	}
	return strings.Join(cells, "")
}

func (t *Table) renderRowData(style lipgloss.Style, data []string, termWidth int) string {
	cells := make([]string, len(t.columns))
	for i, col := range t.columns {
		val := ""
		if i < len(data) {
			val = data[i]
		}
		cells[i] = style.Width(col.Width).Render(truncate(val, col.Width-2))
	}
	return strings.Join(cells, "")
}

// PrintMeta prints a dimmed metadata line above the table.
func PrintMeta(format string, args ...any) {
	fmt.Println(dimStyle.Render(fmt.Sprintf(format, args...)))
	fmt.Println()
}

func truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	if max <= 3 {
		return string(runes[:max])
	}
	return string(runes[:max-1]) + "…"
}

func terminalWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w == 0 {
		return 120
	}
	return w
}
