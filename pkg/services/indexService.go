package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

type SMAs struct {
	SMA40  float64
	SMA80  float64
	SMA200 float64
}

// QuoteResponse represents the JSON structure for the quote response
type IndexResponse struct {
	Symbol      string  `json:"symbol"`
	SMA40       float64 `json:"sma_40"`
	SMA80       float64 `json:"sma_80"`
	SMA200      float64 `json:"sma_200"`
	SMAResult   string  `json:"sma_result"`
	SMAAnalysis string  `json:"sma_analysis"`
	EMA12       float64 `json:"ema_12"`
	EMA26       float64 `json:"ema_26"`
	EMAResult   string  `json:"ema_result"`
	EMAAnalysis string  `json:"ema_analysis"`
}

// getQuote handles the retrieval of stock quotes
func GetIndex(c *gin.Context) {
	symbol := c.Param("symbol")
	startDateTime := &datetime.Datetime{
		Month: int(time.Now().AddDate(-1, 0, 0).Month()),
		Day:   1,
		Year:  time.Now().AddDate(-1, 0, 0).Year(),
	}
	endDateTime := &datetime.Datetime{
		Month: int(time.Now().Month()),
		Day:   time.Now().Day(),
		Year:  time.Now().Year(),
	}
	params := &chart.Params{
		Symbol:   symbol,
		Interval: datetime.OneHour,
		Start:    startDateTime,
		End:      endDateTime,
	}
	iter := chart.Get(params)

	var closes []float64
	for iter.Next() {
		p := iter.Bar()
		v, _ := p.Close.Float64()
		closes = append(closes, v)
	}
	if err := iter.Err(); err != nil {
		fmt.Println(err)
	}

	sma40, err := calculateSMAN(closes, 40)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sma80, err := calculateSMAN(closes, 80)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sma200, err := calculateSMAN(closes, 200)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sma := &SMAs{
		SMA40:  sma40,
		SMA80:  sma80,
		SMA200: sma200,
	}
	smaResult, smaMessage := analyzeSMATrend(*sma)
	// Calculate EMA for 12-day and 26-day periods
	ema12, err := calculateEMA(closes, 12)
	if err != nil {
		log.Fatalf("Error calculating EMA 12: %v", err)
	}

	ema26, err := calculateEMA(closes, 26)
	if err != nil {
		log.Fatalf("Error calculating EMA 26: %v", err)
	}
	emaResult, emaAnalysis := analyzeEMACrossover(ema12, ema26)

	response := IndexResponse{
		Symbol:      symbol,
		SMA40:       sma40,
		SMA80:       sma80,
		SMA200:      sma200,
		SMAResult:   smaResult,
		SMAAnalysis: smaMessage,
		EMA12:       ema12[len(ema12)-1],
		EMA26:       ema26[len(ema26)-1],
		EMAResult:   emaResult,
		EMAAnalysis: emaAnalysis,
	}

	c.JSON(http.StatusOK, response)
}

func calculateSMAN(closes []float64, n int) (float64, error) {
	// Calculate SMA for the last N days
	if len(closes) < n {
		return 0, fmt.Errorf("not enough data to calculate SMA%d", n)
	}

	sum := 0.0
	for i := len(closes) - n; i < len(closes); i++ {
		sum += closes[i]
	}
	sma := sum / float64(n)
	return sma, nil
}

// analyzeSMATrend analyzes the relationship between the SMAs to determine if it's an uptrend
func analyzeSMATrend(smas SMAs) (string, string) {
	if smas.SMA40 > smas.SMA80 && smas.SMA80 > smas.SMA200 {
		return "Uptrend", "SMA 40 > SMA 80 > SMA 200: Esta relación sugiere que el precio de la acción está en una tendencia alcista. Las SMAs más cortas (40 días) están por encima de las SMAs más largas (80 y 200 días), lo que indica que los precios recientes están más altos que los precios pasados."
	} else if smas.SMA40 < smas.SMA80 && smas.SMA80 < smas.SMA200 {
		return "Downtrend", "SMA 40 < SMA 80 < SMA 200: Esta relación sugiere que el precio de la acción está en una tendencia bajista. Las SMAs más cortas (40 días) están por debajo de las SMAs más largas (80 y 200 días), lo que indica que los precios recientes están más bajos que los precios pasados."
	} else if smas.SMA40 > smas.SMA80 && smas.SMA80 < smas.SMA200 {
		return "Shor-term Uptrend Long-term Downtrend", "SMA 40 > SMA 80 < SMA 200: Esta relación sugiere que el precio de la acción puede estar en una recuperación a corto plazo, pero sigue en una tendencia bajista a largo plazo."
	} else if smas.SMA40 < smas.SMA80 && smas.SMA80 > smas.SMA200 {
		return "Short term Downtrend Long-term Uptrend", "SMA 40 < SMA 80 > SMA 200: Esta relación sugiere que el precio de la acción puede estar en una corrección a corto plazo, pero sigue en una tendencia alcista a largo plazo."
	}
	return "Neutral", "Las SMAs no están en un orden claro para confirmar una tendencia específica."

}

// calculateEMA calculates the Exponential Moving Average (EMA)
func calculateEMA(prices []float64, period int) ([]float64, error) {
	if len(prices) < period {
		return nil, fmt.Errorf("not enough data to calculate EMA for the given period")
	}

	// Initialize the EMA array
	ema := make([]float64, len(prices))

	// Calculate the initial SMA to start the EMA
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema[period-1] = sum / float64(period)

	// Calculate the smoothing factor
	alpha := 2.0 / float64(period+1)

	// Calculate the EMA for each subsequent price
	for i := period; i < len(prices); i++ {
		ema[i] = (prices[i] * alpha) + (ema[i-1] * (1 - alpha))
	}

	return ema[period-1:], nil // Return only the valid EMA values
}

// analyzeEMACrossover analyzes the crossover between EMA12 and EMA26
func analyzeEMACrossover(ema12, ema26 []float64) (string, string) {
	if len(ema12) == 0 || len(ema26) == 0 {
		return "Neutral", "Not enough data for analysis."
	}

	// Determine the most recent values
	latestEMA12 := ema12[len(ema12)-1]
	latestEMA26 := ema26[len(ema26)-1]

	// Determine the previous values
	prevEMA12 := ema12[len(ema12)-2]
	prevEMA26 := ema26[len(ema26)-2]

	// Analyze the crossover
	if latestEMA12 > latestEMA26 && prevEMA12 <= prevEMA26 {
		return "Potential uptrend", "Bullish crossover detected. EMA12 has crossed above EMA26, indicating a potential uptrend."
	} else if latestEMA12 < latestEMA26 && prevEMA12 >= prevEMA26 {
		return "Potential downtrend", "Bearish crossover detected. EMA12 has crossed below EMA26, indicating a potential downtrend."
	} else if latestEMA12 > latestEMA26 {
		return "Potential uptrend", "EMA12 is above EMA26, indicating a potential uptrend."
	} else if latestEMA12 < latestEMA26 {
		return "Potential downtrend", "EMA12 is below EMA26, indicating a potential downtrend."
	}
	return "Neutral", "No significant crossover detected."
}
