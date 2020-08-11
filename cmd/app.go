package cmd

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"strings"

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

// Global variables
var (
	Currency   string
	currencies = make(map[string]interface{})
)

// this url represents an up-to-date source of supported currencies for our stonks application
const (
	supportedCurrencyURL = "https://focusmobile-interview-materials.s3.eu-west-3.amazonaws.com/Cheap.Stocks.Internationalization.Currencies.csv"
)

// init function is executed upon start-up of the app
func init() {

	rootCmd.AddCommand(appCmd)

	appCmd.Flags().StringVarP(&Currency, "currency", "c", "", "Input currency is required. eg. KES")
	appCmd.MarkFlagRequired("currency")
}

// run CLI app
func run() {

	// We could check if the passed in value is a string but that won't be neccessary
	fmt.Printf("Processing %v please wait ...\n", Currency)

	// 1. download the csv file
	res, err := http.Get(supportedCurrencyURL)

	if err != nil {
		log.Fatalf("could not fetch supported currencies : %v", err)
	}

	defer res.Body.Close()

	// 2. process the csv file
	vals := csv.NewReader(res.Body)

	for {
		val, err := vals.Read()

		// if we've reached end of the csv, break
		if err == io.EOF {
			break
		}

		// any other internal error
		if err != nil {
			log.Fatalf("cannot read record from csv : %v", err)
		}

		currencies[val[2]] = struct {Country, Currency, Code string} { val[0], val[1], val[2] }
	}

	// 3. search for requested currency
	result, err := queryCurrency(Currency)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)

}

// queries from cache/db wether provided currency is supported
func queryCurrency(param string) (string, error) {

	param = strings.ToUpper(Currency) // convert to uppercase

	if _, ok := currencies[param]; !ok {

		return "", fmt.Errorf("currency %v not supported", param)
	}

	return fmt.Sprintf("Currency supported : %+v", currencies[param]), nil
}
