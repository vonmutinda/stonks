package cmd

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"fmt"
)

var appCmd = &cobra.Command{
	Use:     "app",
	Aliases: []string{"app", "a"},
	Short:   "Start App",
	Long:    `Initiate a session with Stonk`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

// CCC - country currency code
type CCC struct {
	Country, Currency, Code string
}

// Global variables
var (
	currency   string
	Currencies = make(map[string]CCC)
)

// this url represents an up-to-date source of supported currencies for our stonks application
const (
	supportedCurrencyURL = "https://focusmobile-interview-materials.s3.eu-west-3.amazonaws.com/Cheap.Stocks.Internationalization.Currencies.csv"
)

// init function is executed upon start-up of the app
func init() {

	rootCmd.AddCommand(appCmd)

	appCmd.Flags().StringVarP(&currency, "currency", "c", "", "Input currency is required. eg. USD")
	appCmd.MarkFlagRequired("currency")

	// load .env variable
	if err := godotenv.Load(); err != nil {
		log.Fatalf("cannot load .env file. %v", err)
	}

	// fetch supported currencies
	if err := load(); err != nil {
		log.Fatalf("could not load supported currencies : %v", err)
	}
}

// run CLI app
func run() {

	// We could check if the passed in value is a string but that won't be neccessary
	fmt.Printf("Processing %v please wait ...\n\n", currency)

	// 3. search for requested currency
	if query(currency) != true {
		log.Fatalf("Currency %v NOT supported.\n", currency)
	}

	log.Printf("Currency %v is supported.\n", currency)
}

// load supported currencies
func load() error {

	// 1. download the csv file
	res, err := http.Get(supportedCurrencyURL)

	if err != nil {
		return err
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
			return err
		}

		if len(val[2]) == 3 {
			Currencies[val[2]] = CCC{val[0], val[1], val[2]}
		}
	}

	return nil
}

// queries from cache/db wether provided currency is supported
func query(param string) bool {

	param = strings.ToUpper(param) // convert to uppercase

	if _, ok := Currencies[param]; !ok {

		return false
	}

	return true
}
