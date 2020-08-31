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
	Name, URL, From, To, Amount, Result string
}

func convert() {

	var (
		cAPI     = os.Getenv("CLKey")
		fixerAPI = os.Getenv("FKey")

		done = make(chan Service)

		// a buffered channel.
		errChan = make(chan int, 2)

		client = http.DefaultClient
	)

	payload := []Service{
		{Name: "Currency Layer", URL: fmt.Sprintf(currencyLayerURL, cAPI, From, To, Amount), From: From, To: To, Amount: Amount},
		{Name: "Fixer.io", URL: fmt.Sprintf(fixerURL, fixerAPI, From, To, Amount), From: From, To: To, Amount: Amount},
	}

	// use timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)

	defer cancel()

	// Make two concurrent requests for currency conversion.
	// whichever returns result first we'll go with that.
	for _, k := range payload {

		go func(ct context.Context, val Service, result chan Service, eChan chan int) {
			// log behavior of the service

			// check if requested currencies are supported
			if query(val.From) != true || query(val.To) != true {
				fmt.Printf("[%v] failed. conversion for the currencies not supported\n", val.Name)
			}

			// prepare the request
			req, err := http.NewRequest(http.MethodGet, val.URL, nil)

			if err != nil {
				fmt.Printf("[%v] cannot prepare request. %v\n", val.Name, err)
			}

			// make the request for currency conversion
			res, err := client.Do(req)

			if err != nil {
				fmt.Printf("[%v] cannot make request. %v", val.Name, err)
			}

			defer res.Body.Close()

			// read data
			data, err := ioutil.ReadAll(res.Body)

			// fmt.Printf("data : %s\n\n", data)

			if err != nil {
				fmt.Printf("[%v] cannot read response data. %v", val.Name, err)
			}

			// check status
			status, err := jsonparser.GetBoolean(data, "success")

			if err != nil {
				fmt.Printf("[%v] cannot read success status. %v", val.Name, err)
			}

			if status != true {

				// read error response
				erresp, err := jsonparser.GetString(data, "error", "info")

				if err != nil {
					fmt.Printf("[%v] cannot read info key in error body. %v", val.Name, err)
				}

				fmt.Printf("[%v] Conversion failed. %v\n", val.Name, erresp)

				eChan <- 1
			}

			// read data
			rslt, err := jsonparser.GetString(data, "result")

			if err != nil {
				fmt.Printf("[%v] cannot read conversion result. %v\n", val.Name, err)
			}

			// if there's a returned result
			if len(rslt) > 0 && status {

				val.Result = rslt

				done <- val
			}

		}(ctx, k, done, errChan)

	}

	// listen for results
	var counter = 0
	for {
		select {
		case res := <-done:
			fmt.Printf("%v %v\n", res.To, res.Result)
			os.Exit(0)
		case <-ctx.Done():
			// both conversion services must have failed or took too long
			log.Fatalf("Conversion services timed out!")
		case _, ok := <-errChan:
			if ok {
				counter++
				if counter == 2 {
					log.Fatalf("All currency conversion services failed") // os.Exit(1)
				}
			}
		}
	}

}
