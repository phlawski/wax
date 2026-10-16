package cmd

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"wax/internal/ui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := newClient()
		m := ui.New(c)
		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
