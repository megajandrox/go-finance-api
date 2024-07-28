package dto

import "github.com/megajandrox/go-finance-api/pkg/models"

type BuyPosition struct {
	Symbol     string            `json:"symbol"`      // Financial asset symbol
	Price      float64           `json:"price"`       // price
	Quantity   int               `json:"quantity"`    // Quantity of shares/buys
	MarketType models.MarketType `json:"market_type"` // Market (Equities, ETFs)
}

type SellPosition struct {
	Symbol   string  `json:"symbol"`   // Financial asset symbol
	Price    float64 `json:"price"`    // price
	Quantity int     `json:"quantity"` // Quantity of shares/buys
}
