package services

import (
	"fmt"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/quote"
)

// getQuote handles the retrieval of stock quotes
func FindQuote(symbol string) (*finance.Quote, error) {
	q, err := quote.Get(symbol)
	if err != nil {
		return nil, err
	}
	if q == nil {
		return nil, fmt.Errorf("Quote not found")
	}

	return q, nil
}
