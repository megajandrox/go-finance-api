package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/megajandrox/go-finance-api/pkg/utils"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

// QuoteResponse represents the JSON structure for the quote response
type IndexResponse struct {
	Symbol                       string  `json:"symbol"`
	SMA40                        float64 `json:"sma_40"`
	SMA80                        float64 `json:"sma_80"`
	SMA200                       float64 `json:"sma_200"`
	SMAResult                    string  `json:"sma_result"`
	SMAAnalysis                  string  `json:"sma_analysis"`
	EMA12                        float64 `json:"ema_12"`
	EMA26                        float64 `json:"ema_26"`
	EMAResult                    string  `json:"ema_result"`
	EMAAnalysis                  string  `json:"ema_analysis"`
	MACDAnalysis                 string  `json:"macd_analysis"`
	RSIAnalysis                  string  `json:"rsi_analysis"`
	StochasticOscillatorAnalysis string  `json:"stochastic_oscillator_analysis"`
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
	var highs []float64
	var lows []float64
	for iter.Next() {
		p := iter.Bar()
		v, _ := p.Close.Float64()
		closes = append(closes, v)
		high, _ := p.High.Float64()
		low, _ := p.Low.Float64()
		highs = append(highs, high)
		lows = append(lows, low)
	}
	if err := iter.Err(); err != nil {
		fmt.Println(err)
	}

	sma40, err := utils.CalculateSMAN(closes, 40)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sma80, err := utils.CalculateSMAN(closes, 80)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sma200, err := utils.CalculateSMAN(closes, 200)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sma := &utils.SMAs{
		SMA40:  sma40,
		SMA80:  sma80,
		SMA200: sma200,
	}
	smaResult, smaMessage := utils.AnalyzeSMATrend(*sma)
	// Calculate EMA for 12-day and 26-day periods
	ema12, err := utils.CalculateEMA(closes, 12)
	if err != nil {
		log.Fatalf("Error calculating EMA 12: %v", err)
	}

	ema26, err := utils.CalculateEMA(closes, 26)
	if err != nil {
		log.Fatalf("Error calculating EMA 26: %v", err)
	}
	emaResult, emaAnalysis := utils.AnalyzeEMACrossover(ema12, ema26)

	// Calculate MACD and Signal line
	macd, signal, err := utils.CalculateMACD(closes)
	if err != nil {
		log.Fatalf("Error calculating MACD: %v", err)
	}

	// Analyze MACD
	macdAnalysis := utils.AnalyzeMACD(macd, signal)

	// Calculate RSI value
	// Calculate RSI for 14-day period
	rsi, err := utils.CalculateRSI(closes, 14)
	if err != nil {
		log.Fatalf("Error calculating RSI: %v", err)
	}
	// Get the last RSI value
	latestRSI := rsi[len(rsi)-1]

	// Analyze the RSI value
	rsiAnalysis := utils.AnalyzeRSI(latestRSI)

	// Calculate Stochastic Oscillator for 14-day period
	k, d, err := utils.CalculateStochasticOscillator(closes, highs, lows, 14)
	if err != nil {
		log.Fatalf("Error calculating Stochastic Oscillator: %v", err)
	}

	// Analyze the Stochastic Oscillator
	stochasticOscillatorAnalysis := utils.AnalyzeStochasticOscillator(k, d)

	response := IndexResponse{
		Symbol:                       symbol,
		SMA40:                        sma40,
		SMA80:                        sma80,
		SMA200:                       sma200,
		SMAResult:                    smaResult,
		SMAAnalysis:                  smaMessage,
		EMA12:                        ema12[len(ema12)-1],
		EMA26:                        ema26[len(ema26)-1],
		EMAResult:                    emaResult,
		EMAAnalysis:                  emaAnalysis,
		MACDAnalysis:                 macdAnalysis,
		RSIAnalysis:                  rsiAnalysis,
		StochasticOscillatorAnalysis: stochasticOscillatorAnalysis,
	}

	c.JSON(http.StatusOK, response)
}
