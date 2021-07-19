package finnhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Quote reprsents a stock quote
type Quote struct {
	symbol    string
	open      float64
	high      float64
	low       float64
	cur       float64
	prevClose float64
	time      time.Time // time of information
}

func (q Quote) String() string {
	return fmt.Sprintf(`%s:
	As of  : %s
	Current: %.2f
	Low    : %.2f
	High   : %.2f
	Open   : %.2f
	Previous Close: %.2f`,
		q.symbol, q.time.String(), q.cur, q.low, q.high, q.open, q.prevClose,
	)
}

// quoteResponseBody holds the results of a quote api request
type quoteResponseBody struct {
	Time      float64 `json:"t"`
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Cur       float64 `json:"c"`
	PrevClose float64 `json:"pc"`
}

func (q *quoteResponseBody) toQuote(symbol string) Quote {
	return Quote{
		symbol:    symbol,
		open:      q.Open,
		high:      q.High,
		low:       q.Low,
		cur:       q.Cur,
		prevClose: q.PrevClose,
		time:      time.Unix(int64(q.Time), int64(0)),
	}
}

// An API is a wrapper around the Finnhub.io api.
//
// Before use, you must set your api token via API.SetToken(token).
type API struct {
	token string
}

// SetToken sets your finnhub api token, enabling api calls
func (api *API) SetToken(token string) {
	api.token = token
}

func (api *API) noToken() bool {
	return api.token == ""
}

const (
	baseURL = "https://finnhub.io/api/v1/"

	errTokenNotSet = "token not set, use API.SetToken()"
)

func (api *API) Quote(symbol string) (Quote, error) {
	// current location for more informative error messages
	const errCaller = "API.Quote()"

	if api.noToken() {
		return Quote{}, errorFrom(errCaller, errTokenNotSet)
	}
	if len(symbol) == 0 {
		return Quote{}, errorFrom(errCaller, "no symbol given")
	}

	reqURL := baseURL + "quote"

	client := &http.Client{}
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return Quote{}, fmt.Errorf("%s: %v", errCaller, err)
	}

	query := url.Values{}
	query.Set("symbol", symbol)
	req.URL.RawQuery = query.Encode()
	req.Header.Add("X-Finnhub-Token", api.token)

	var quote quoteResponseBody

	resp, err := client.Do(req)
	if err != nil {
		return Quote{}, fmt.Errorf("%s: %v", errCaller, err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// defer resp.Body.Close() // is this needed?
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&quote)
		if err != nil {
			return Quote{}, fmt.Errorf("%s error decoding response body: %v", errCaller, err)
		}
	} else {
		return Quote{}, fmt.Errorf(`%s request for "%s" failed\n`, errCaller, symbol)
	}
	return quote.toQuote(symbol), nil
}

// errorFrom is a helper function to construct error messages with
// location information
func errorFrom(from string, msg string) error {
	return fmt.Errorf("%s %s", from, msg)
}
