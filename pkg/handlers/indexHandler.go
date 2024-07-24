package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/megajandrox/go-finance-api/pkg/services"
)

// QuoteResponse represents the JSON structure for the quote response
type IndexResponse struct {
	Symbol                       string `json:"symbol"`
	SMAResult                    string `json:"sma_result"`
	SMAAnalysis                  string `json:"sma_analysis"`
	EMAResult                    string `json:"ema_result"`
	EMAAnalysis                  string `json:"ema_analysis"`
	MACDResult                   string `json:"macd_result"`
	MACDAnalysis                 string `json:"macd_analysis"`
	RSIResult                    string `json:"rsi_result"`
	RSIAnalysis                  string `json:"rsi_analysis"`
	StochasticOscillatorResult   string `json:"stochastic_oscillator_result"`
	StochasticOscillatorAnalysis string `json:"stochastic_oscillator_analysis"`
	VolumeAnalysis               string `json:"volume_analysis"`
	OBVAnalysis                  string `json:"obv_analysis"`
	RVOLAnalysis                 string `json:"rvol_analysis"`
	ADXResult                    string `json:"adx_result"`
	ADXAnalysis                  string `json:"adx_analysis"`
	MomentumResult               string `json:"momentum_result"`
	MomentumAnalysis             string `json:"momentum_analysis"`
	CCIResult                    string `json:"cci_result"`
	CCIAnalysis                  string `json:"cci_analysis"`
}

// getQuote handles the retrieval of stock quotes
func GetIndex(c *gin.Context) {
	symbol := c.Param("symbol")
	fromParam := c.Query("from")
	intervalParam := c.Query("interval")
	if fromParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query parameter 'from' is required.",
		})
		return
	}

	// Convert the query parameter `from` to an integer
	from, err := strconv.Atoi(fromParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameter. 'from' must be an integer.",
		})
		return
	}

	if from < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameter. 'from' must be greater than 1.",
		})
		return
	}

	if !IsValidInterval(Interval(intervalParam)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameter. 'Interval' must be at least 1h or 1d.",
		})
		return
	}

	indexes, _ := services.FindIndexesBySymbol(symbol, from, intervalParam)
	response := IndexResponse{
		Symbol:                       symbol,
		SMAResult:                    indexes.SMA.TrendType.String(),
		SMAAnalysis:                  indexes.SMA.Result,
		EMAResult:                    indexes.EMA.TrendType.String(),
		EMAAnalysis:                  indexes.EMA.Result,
		MACDResult:                   indexes.MACD.TrendType.String(),
		MACDAnalysis:                 indexes.MACD.Result,
		RSIResult:                    indexes.RSI.TrendType.String(),
		RSIAnalysis:                  indexes.RSI.Result,
		StochasticOscillatorResult:   indexes.Stochastic.TrendType.String(),
		StochasticOscillatorAnalysis: indexes.Stochastic.Result,
		VolumeAnalysis:               indexes.Volume.Result,
		OBVAnalysis:                  indexes.OBV.Result,
		RVOLAnalysis:                 indexes.RVOL.Result,
		ADXResult:                    indexes.ATR.TrendType.String(),
		ADXAnalysis:                  indexes.ATR.Result,
		MomentumResult:               indexes.Momentum.TrendType.String(),
		MomentumAnalysis:             indexes.Momentum.Result,
		CCIResult:                    indexes.CCI.TrendType.String(),
		CCIAnalysis:                  indexes.CCI.Result,
	}

	c.JSON(http.StatusOK, response)
}
