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
	Use:     "Convert",
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

	convertCMD.Flags().StringVarP(&From, "from", "f", "", "Input currency is required. eg. KES")
	convertCMD.Flags().StringVarP(&To, "to", "t", "", "Input currency is required. eg. KES")
	convertCMD.Flags().StringVarP(&Amount, "amount", "a", "", "Input currency is required. eg. KES")

	convertCMD.MarkFlagRequired("from")
	convertCMD.MarkFlagRequired("to")
	convertCMD.MarkFlagRequired("amount")
}

// Service -
type Service struct {
	URL    string
	From   string
	To     string
	Amount string
	Result string
}

func convert() {

	var done = make(chan Service)
	var client = http.DefaultClient

	// load .env variable
	if err := godotenv.Load(); err != nil {
		log.Fatalf("cannot load .env file. %v", err)
	}

	cAPI := os.Getenv("CLKey")
	fixerAPI := os.Getenv("FKey")

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

			// prepare the request
			req, err := http.NewRequest(http.MethodGet, val.URL, nil)

			if err != nil {
				log.Printf("cannot prepare request. %v", err)
			}

			// make the request for currency conversion
			res, err := client.Do(req)

			if err != nil {
				log.Printf("cannot make request. %v", err)
			}

			defer res.Body.Close()

			// read data
			data, err := ioutil.ReadAll(res.Body)

			// fmt.Printf("data : %s\n\n", data)

			if err != nil {
				log.Printf("cannot read response data. %v", err)
			}

			// check status
			status, err := jsonparser.GetBoolean(data, "success")

			if err != nil {
				log.Printf("cannot read success status. %v", err)
			}

			// read error response
			erresp, err := jsonparser.GetString(data, "error", "info")

			if err != nil {
				log.Printf("cannot read info key in error body. %v", err)
			}

			if status != true {
				log.Printf("Conversion failed. %v", erresp)
			}

			// read data
			rslt, err := jsonparser.GetString(data, "result")

			if err != nil {
				log.Printf("cannot read conversion result. %v", err)
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
			os.Exit(1)
		}
	}

}
