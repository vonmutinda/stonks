package cmd

import (
	"log"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

var infoCMD = &cobra.Command{
	Use:     "Info",
	Aliases: []string{"info", "c"},
	Short:   "Info about current stock price",
	Run: func(cmd *cobra.Command, args []string) {
		info()
	},
}

func init() {
	rootCmd.AddCommand(infoCMD)
}

// info displays current stock price
func info() {

	var rand = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(500)

	log.Printf("The current price for AAPL is %v USD", rand)
}
