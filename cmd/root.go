package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"wax/internal/config"
	"wax/pkg/output"
)

var rootCmd = &cobra.Command{
	Use:   "wax",
	Short: "Discogs CLI for browsing your collection and the marketplace",
	Long: `
 ██╗    ██╗ █████╗ ██╗  ██╗
 ██║    ██║██╔══██╗╚██╗██╔╝
 ██║ █╗ ██║███████║ ╚███╔╝
 ██║███╗██║██╔══██║ ██╔██╗
 ╚███╔███╔╝██║  ██║██╔╝ ██╗
  ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝  ╚═╝

Discogs CLI — browse your vinyl collection and the marketplace.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := config.Load(); err != nil {
			log.Fatalf("config error: %s", err)
		}
		if err := viper.Unmarshal(&config.AppConfig); err != nil {
			log.Fatalf("config decode error: %s", err)
		}
		output.Debug("config loaded: token=%s... username=%s",
			maskToken(config.AppConfig.Discogs.Token), config.AppConfig.Discogs.Username)
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func maskToken(t string) string {
	if len(t) <= 4 {
		return "****"
	}
	return t[:4] + "****"
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&output.DebugEnabled, "debug", "d", false, "Enable debug output")
	rootCmd.AddCommand(collectionCmd, searchCmd, marketplaceCmd)
}
