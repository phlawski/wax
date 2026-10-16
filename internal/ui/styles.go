package ui

import "github.com/charmbracelet/lipgloss"

// Solarized Dark palette
var (
	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#002b36")).
			Background(lipgloss.Color("#268bd2")).
			PaddingLeft(1).PaddingRight(1)

	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#268bd2"))

	styleHelp = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#586e75"))

	styleStatus = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#93a1a1"))

	styleError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#dc322f"))

	styleFilter = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2aa198"))

	styleSelected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fdf6e3")).
			Background(lipgloss.Color("#268bd2")).
			Bold(true)

	styleRowEven = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#839496")).
			Background(lipgloss.Color("#002b36"))

	styleRowOdd = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#93a1a1")).
			Background(lipgloss.Color("#0a4452"))

	styleDim = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#586e75"))
)
