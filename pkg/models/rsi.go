package models

import (
	"errors"
	"fmt"
)

type RSI struct {
	symbol    string
	LatestRSI float64
	TrendType TrendType
	Result    string
}

func NewRSI(symbol string) (*RSI, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	return &RSI{symbol: symbol}, nil
}

func (rsi *RSI) calculate(closes []float64) (bool, error) {
	// Calculate RSI for 14-day period
	rsiArray, err := rsi.calculateRSI(closes, 14)
	if err != nil {
		return false, fmt.Errorf(`Error calculating RSI: %v`, err)
	}
	// Get the last RSI value
	rsi.LatestRSI = rsiArray[len(rsiArray)-1]
	return true, nil
}

// calculateRSI calculates the Relative Strength Index (RSI)
func (rsi *RSI) calculateRSI(prices []float64, period int) ([]float64, error) {
	if len(prices) < period {
		return nil, fmt.Errorf("not enough data to calculate RSI for the given period")
	}

	// Initialize gains and losses
	gains := make([]float64, len(prices))
	losses := make([]float64, len(prices))

	// Calculate daily gains and losses
	for i := 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains[i] = change
		} else {
			losses[i] = -change
		}
	}

	// Calculate average gains and losses for the first period
	avgGain := 0.0
	avgLoss := 0.0
	for i := 1; i <= period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Initialize RSI array
	rsiArray := make([]float64, len(prices))

	// Calculate RSI for the first period
	rs := avgGain / avgLoss
	rsiArray[period-1] = 100 - (100 / (1 + rs))

	// Calculate RSI for subsequent periods
	for i := period; i < len(prices); i++ {
		avgGain = (avgGain*float64(period-1) + gains[i]) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + losses[i]) / float64(period)
		rs = avgGain / avgLoss
		rsiArray[i] = 100 - (100 / (1 + rs))
	}

	return rsiArray[period-1:], nil // Return only the valid RSI values
}

// analyzeRSI analyzes the RSI value and returns a descriptive analysis
func (rsi *RSI) AnalyzeRSI(closes []float64) {
	status, err := rsi.calculate(closes)
	var trendType TrendType = Neutral
	var result string = "Las RSI no están en un orden claro para confirmar una tendencia específica."
	if !status {
		trendType = Neutral
		result = fmt.Sprintf("No es posible calcular RSI debido a: %s", err)
		rsi.TrendType = trendType
		rsi.Result = result
		return
	}
	if rsi.LatestRSI > 70 {
		trendType = Overbought
		result = fmt.Sprintf("RSI is %.2f, indicating the asset is overbought.", rsi.LatestRSI)
	} else if rsi.LatestRSI < 30 {
		trendType = Oversold
		result = fmt.Sprintf("RSI is %.2f, indicating the asset is oversold.", rsi.LatestRSI)
	} else {
		trendType = Neutral
		result = fmt.Sprintf("RSI is %.2f, indicating normal market conditions.", rsi.LatestRSI)
	}
	rsi.TrendType = trendType
	rsi.Result = result
}
