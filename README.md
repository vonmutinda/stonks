## Stonks
- An application for managing Cheap Stonk Inc services. 

### Local Set Up  
+ Install [Go](https://golang.org/dl/)
+ Clone the [repo](https://github.com/vonmutinda/stonks.git)  
+ Get API keys from [Fixer] And [CurrencyLayer] and create a `.env` file.

### Usage 
```cmd
    
    - Help:             `go run main.go -h`
    - Currency Info     `go run main.go app -c KES`
    - Convert           `go run main.go convert -f KES -t USD -a 100`


```
## Technologies used:  
- [Golang version `go1.14.6`](https://golang.org) Go programming language.
- [Cobra](https://github.com/spf13/cobra) for drop-in CLI development. Simple and seamless.
