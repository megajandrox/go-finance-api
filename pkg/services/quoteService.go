package services

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/piquette/finance-go/quote"
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
	q, err := quote.Get(symbol)
	if err != nil {
		if strings.Contains(err.Error(), "Can't find quote for symbol") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if q == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		return
	}

	response := QuoteResponse{
		Symbol:                     q.Symbol,
		ShortName:                  q.ShortName,
		RegularMarketPrice:         q.RegularMarketPrice,
		RegularMarketPreviousClose: q.RegularMarketPreviousClose,
		Open:                       q.RegularMarketOpen,
		High:                       q.RegularMarketDayHigh,
		Low:                        q.RegularMarketDayLow,
		RegularMarketVolume:        q.RegularMarketVolume,
		AverageDailyVolume10Day:    q.AverageDailyVolume10Day,
		AverageDailyVolume3Month:   q.AverageDailyVolume3Month,
	}

	c.JSON(http.StatusOK, response)
}
