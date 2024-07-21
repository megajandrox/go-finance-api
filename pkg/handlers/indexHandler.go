package handlers

import (
	"net/http"

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
}

// getQuote handles the retrieval of stock quotes
func GetIndex(c *gin.Context) {
	symbol := c.Param("symbol")
	indexes, _ := services.FindIndexesBySymbol(symbol)
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
	}

	c.JSON(http.StatusOK, response)
}