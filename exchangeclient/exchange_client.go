package exchangeclient

import (
	"fmt"
	"github.com/ervitis/exchangerateapp/exchangeclient/client"
	"github.com/ervitis/exchangerateapp/exchangeclient/client/rates"
	"github.com/ervitis/exchangerateapp/exchangeclient/client/symbols"
	"github.com/go-openapi/strfmt"
	"os"
)

type (
	ExchangeApi struct {
		URL        string
		APIKEY     string
		client     *client.FixerIo
		currencies []string
	}
)

func NewClient() (*ExchangeApi, error) {
	s := os.Getenv("API_URL")
	if s == "" {
		return nil, fmt.Errorf("API_URL env variable not set")
	}

	c := &ExchangeApi{}
	c.URL = s

	s = os.Getenv("API_KEY")
	if s == "" {
		return nil, fmt.Errorf("API_KEY env variable not set")
	}
	c.APIKEY = s

	c.client = client.NewHTTPClient(strfmt.Default)

	if err := c.getCurrencies(); err != nil {
		return nil, fmt.Errorf("newClient: error creating client: %w", err)
	}

	return c, nil
}

func (e *ExchangeApi) getCurrencies() error {
	if e.currencies != nil || len(e.currencies) != 0 {
		return nil
	}

	resp, err := e.client.Symbols.GetSymbols(&symbols.GetSymbolsParams{AccessKey: e.APIKEY})
	if err != nil {
		return fmt.Errorf("getCurrencies: error retrieving symbols from api: %w", err)
	}

	if !resp.GetPayload().Success {
		return fmt.Errorf("getCurrencies: the request was not success")
	}

	e.currencies = make([]string, 0)
	for cu := range resp.GetPayload().Symbols {
		e.currencies = append(e.currencies, cu)
	}
	return nil
}

func (e *ExchangeApi) validateInputFrom(from string) error {
	if !contains(e.currencies, from) {
		return fmt.Errorf("validateInputFrom: the currency %s is not valid", from)
	}
	return nil
}

func (e *ExchangeApi) validateInputTo(to []string) error {
	for _, v := range to {
		if !contains(e.currencies, v) {
			return fmt.Errorf("validateInputTo: the currencies %s is not valid", v)
		}
	}
	return nil
}

func contains(l []string, e string) bool {
	for _, v := range l {
		if v == e {
			return true
		}
	}
	return false
}

func (e *ExchangeApi) ConvertCurrency(from string, to []string) (ConvertedCurrency, error) {
	if err := e.validateInputTo(to); err != nil {
		return nil, fmt.Errorf("ConvertCurrency: error validating: %w", err)
	}

	if err := e.validateInputFrom(from); err != nil {
		return nil, fmt.Errorf("ConvertCurrency: error validating: %w", err)
	}

	okResponse, _, _, err := e.client.Rates.GetLatest(&rates.GetLatestParams{AccessKey: e.APIKEY, Base: &from, Symbols: to})
	if err != nil {
		return nil, fmt.Errorf("ConvertCurrency: error GetLatest: %w", err)
	}

	cc := make(ConvertedCurrency, 0)
	for k, v := range okResponse.GetPayload().Rates {
		cc[k] = v
	}
	return cc, nil
}
