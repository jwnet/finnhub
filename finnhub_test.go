package finnhub

import (
	"fmt"
	"testing"
)

var apiToken string

func TestMain(m *testing.M) {
	token, ok := os.LookupEnv("FINNHUB_API_TOKEN")
	if ! ok {
		log.Fatal("TestMain(): FINNHUB_API_TOKEN environment variable not set")
	}

	// set global apiToken
	apiToken = token

	os.Exit(m.Run())
}

func TestQuote(t *testing.T) {
	api := API{}
	api.SetToken(apiToken)
	symbol := "GME"
	quote, err := api.Quote(symbol)
	if err != nil {
		t.Logf("API.Quote(\"%s\") failed: %v", symbol, err)
	}
	fmt.Println(quote)
}
