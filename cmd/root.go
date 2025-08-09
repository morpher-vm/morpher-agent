package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"morpher-agent/cmd/version"
)

var rootCmd = &cobra.Command{
	Use:           "morpher-agent",
	Short:         "Lightweight agent that collects VM information for migration",
	Long:          `morpher-agent is a lightweight agent that collects VM information for migration.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(version.VersionCmd)
}
