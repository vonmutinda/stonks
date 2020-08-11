package cmd

import (
	"github.com/spf13/cobra"

	"fmt"
)

var appCmd = &cobra.Command{
	Use:     "App",
	Aliases: []string{"app", "a"},
	Short:   "Start App",
	Long:    `Initiate a session with Stonk`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

// Currency variable for currency flag
var Currency string

// init function is executed upon start-up of the app
func init() {

	rootCmd.AddCommand(appCmd)

	appCmd.Flags().StringVarP(&Currency, "currency", "c", "", "Input currency is required. eg. KES")
	appCmd.MarkFlagRequired("currency")
}

// run CLI app
func run() {

	fmt.Println("Looking up support for currency : ", Currency) 
}