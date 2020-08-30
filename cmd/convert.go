package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/buger/jsonparser"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var convertCMD = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"convert", "c"},
	Short:   "Currency conversion service",
	Run: func(cmd *cobra.Command, args []string) {
		convert()
	},
}

const (
	currencyLayerURL = "http://apilayer.net/api/convert?access_key=%v&from=%v&to=%v&amount=%v"
	fixerURL         = "http://data.fixer.io/api/convert?access_key=%v&from=%v&to=%v&amount=%v"
)

var (
	// From -
	From string
	// To -
	To string
	// Amount -
	Amount string
)

// init function is executed upon start-up of the app
func init() {

	rootCmd.AddCommand(convertCMD)

	convertCMD.Flags().StringVarP(&From, "from", "f", "", "From is required. eg. KES")
	convertCMD.Flags().StringVarP(&To, "to", "t", "", "To is required. eg. USD")
	convertCMD.Flags().StringVarP(&Amount, "amount", "a", "", "Amount is required. eg. 100")

	convertCMD.MarkFlagRequired("from")
	convertCMD.MarkFlagRequired("to")
	convertCMD.MarkFlagRequired("amount")

	// load .env variable
	if err := godotenv.Load(); err != nil {
		log.Fatalf("cannot load .env file. %v", err)
	}

	// fetch supported currencies
	if err := load(); err != nil {
		log.Fatalf("could not load supported currencies : %v", err)
	}
}

// Service -
type Service struct {
	URL, From, To, Amount, Result string
}

func convert() {

	var (
		cAPI     = os.Getenv("CLKey")
		fixerAPI = os.Getenv("FKey")

		done = make(chan Service)

		client = http.DefaultClient
	)

	payload := []Service{
		{URL: fmt.Sprintf(currencyLayerURL, cAPI, From, To, Amount), From: From, To: To, Amount: Amount},
		{URL: fmt.Sprintf(fixerURL, fixerAPI, From, To, Amount), From: From, To: To, Amount: Amount},
	}

	// use timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// Make two concurrent requests for currency conversion.
	// whichever returns result first we'll go with that.
	for _, k := range payload {

		go func(ct context.Context, val Service, result chan Service) {

			// check if requested currencies are supported
			if query(val.From) != true || query(val.To) != true {
				fmt.Printf("from %v. to %v\n", val.From, val.To)
				log.Fatalf("failed. conversion for the currencies not supported")
			}

			// prepare the request
			req, err := http.NewRequest(http.MethodGet, val.URL, nil)

			if err != nil {
				log.Fatalf("cannot prepare request. %v", err)
			}

			// make the request for currency conversion
			res, err := client.Do(req)

			if err != nil {
				log.Fatalf("cannot make request. %v", err)
			}

			defer res.Body.Close()

			// read data
			data, err := ioutil.ReadAll(res.Body)

			// fmt.Printf("data : %s\n\n", data)

			if err != nil {
				log.Fatalf("cannot read response data. %v", err)
			}

			// check status
			status, err := jsonparser.GetBoolean(data, "success")

			if err != nil {
				log.Fatalf("cannot read success status. %v", err)
			}

			// read error response
			erresp, err := jsonparser.GetString(data, "error", "info")

			if err != nil {
				log.Fatalf("cannot read info key in error body. %v", err)
			}

			if status != true {
				fmt.Printf("Conversion failed. %v\n", erresp)
			}

			// read data
			rslt, err := jsonparser.GetString(data, "result")

			if err != nil {
				fmt.Printf("cannot read conversion result. %v\n", err)
			}

			// if there's a returned result
			if len(rslt) > 0 && status {

				val.Result = rslt

				done <- val
			}

		}(ctx, k, done)

	}

	// listen for results
	for {
		select {
		case res := <-done:
			fmt.Printf("%v %v\n", res.To, res.Result)
			os.Exit(0)
		case <-ctx.Done():
			// both conversion services must have failed or took too long
			os.Exit(1)
		}
	}

}
