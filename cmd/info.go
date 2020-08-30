package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

var infoCMD = &cobra.Command{
	Use:     "info",
	Aliases: []string{"info", "i"},
	Short:   "Info about current stock price",
	Run: func(cmd *cobra.Command, args []string) {
		info()
	},
}

var supportedCMD = &cobra.Command{
	Use:     "supported",
	Aliases: []string{"supported"},
	Short:   "Lists all supported currencies",
	Run: func(cmd *cobra.Command, args []string) {
		supported()
	},
}

func init() {
	rootCmd.AddCommand(infoCMD)
	rootCmd.AddCommand(supportedCMD)
}

// info displays current stock price
func info() {
	var rand = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(500)

	log.Printf("The current price for AAPL is %v USD", rand)
}

func supported() {

	for key, val := range Currencies {
		fmt.Printf("%v	-- %v\n", key, val.Currency)
	}
}
