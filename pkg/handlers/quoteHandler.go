package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/megajandrox/go-finance-api/pkg/services"
)

// QuoteResponse represents the JSON structure for the quote response
type QuoteResponse struct {
	Symbol                     string  `json:"symbol"`
	ShortName                  string  `json:"short_name"`
	RegularMarketPrice         float64 `json:"regular_market_price"`
	RegularMarketPreviousClose float64 `json:"regular_market_previous_close"`
	Open                       float64 `json:"open"`
	High                       float64 `json:"high"`
	Low                        float64 `json:"low"`
	RegularMarketVolume        int     `json:"regular_market_volume"`
	AverageDailyVolume10Day    int     `json:"average_daily_volume_10_day"`
	AverageDailyVolume3Month   int     `json:"average_daily_volume_3_month"`
}

// getQuote handles the retrieval of stock quotes
func GetQuote(c *gin.Context) {
	symbol := c.Param("symbol")
	quote, _ := services.FindQuote(symbol)

	response := QuoteResponse{
		Symbol:                     quote.Symbol,
		ShortName:                  quote.ShortName,
		RegularMarketPrice:         quote.RegularMarketPrice,
		RegularMarketPreviousClose: quote.RegularMarketPreviousClose,
		Open:                       quote.RegularMarketOpen,
		High:                       quote.RegularMarketDayHigh,
		Low:                        quote.RegularMarketDayLow,
		RegularMarketVolume:        quote.RegularMarketVolume,
		AverageDailyVolume10Day:    quote.AverageDailyVolume10Day,
		AverageDailyVolume3Month:   quote.AverageDailyVolume3Month,
	}

	c.JSON(http.StatusOK, response)
}
