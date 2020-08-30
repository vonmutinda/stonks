package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "cheap",
	Aliases: []string{"stock", "stonk"},
	Short:   "CLI tool for interacting with Cheap Stock services.",
	Long:    `CheapStock app for stock trading and services in Africa & Beyond`,
}

// Execute root command.
func Execute() error {
	return rootCmd.Execute()
}
