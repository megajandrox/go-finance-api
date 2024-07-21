package models

import (
	"errors"
	"fmt"
)

type EMA struct {
	symbol    string
	EMA12     []float64
	EMA26     []float64
	TrendType TrendType
	Result    string
}

func NewEMA(symbol string) (*EMA, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	return &EMA{symbol: symbol}, nil
}

func (ema *EMA) calculate(closes []float64) (bool, error) {
	// Calculate EMA for 12-day and 26-day periods
	ema12, err := ema.CalculateEMA(closes, 12)
	if err != nil {
		return false, fmt.Errorf("Error calculating EMA 12: %v", err)
	}
	ema.EMA12 = ema12
	ema26, err := ema.CalculateEMA(closes, 26)
	if err != nil {
		return false, fmt.Errorf("Error calculating EMA 26: %v", err)
	}
	ema.EMA26 = ema26
	return true, nil
}

// calculateEMA calculates the Exponential Moving Average (EMA)
func (ema *EMA) CalculateEMA(prices []float64, period int) ([]float64, error) {
	if len(prices) < period {
		return nil, fmt.Errorf("not enough data to calculate EMA for the given period")
	}

	// Initialize the EMA array
	emaArray := make([]float64, len(prices))

	// Calculate the initial SMA to start the EMA
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	emaArray[period-1] = sum / float64(period)

	// Calculate the smoothing factor
	alpha := 2.0 / float64(period+1)

	// Calculate the EMA for each subsequent price
	for i := period; i < len(prices); i++ {
		emaArray[i] = (prices[i] * alpha) + (emaArray[i-1] * (1 - alpha))
	}

	return emaArray[period-1:], nil // Return only the valid EMA values
}

// analyzeEMACrossover analyzes the crossover between EMA12 and EMA26
func (ema *EMA) AnalyzeEMACrossover(closes []float64) {
	var ema12, ema26 []float64
	var trendType TrendType = Neutral
	var result string = "The SMAs are not in a clear order to confirm a specific trend."
	status, err := ema.calculate(closes)
	if !status {
		trendType = Neutral
		result = fmt.Sprintf("It is not possible to calculate EMA due to: %s", err)
		ema.TrendType = trendType
		ema.Result = result
		return
	}
	ema12 = ema.EMA12
	ema26 = ema.EMA26
	if len(ema12) == 0 || len(ema26) == 0 {
		trendType = Neutral
		result = "Not enough data for analysis."
		ema.TrendType = trendType
		ema.Result = result
		return
	}

	// Determine the most recent values
	latestEMA12 := ema12[len(ema12)-1]
	latestEMA26 := ema26[len(ema26)-1]

	// Determine the previous values
	prevEMA12 := ema12[len(ema12)-2]
	prevEMA26 := ema26[len(ema26)-2]

	// Analyze the crossover
	trendType = Neutral
	result = "No significant crossover detected."
	if latestEMA12 > latestEMA26 && prevEMA12 <= prevEMA26 {
		trendType = Potential_Uptrend
		result = "Bullish crossover detected. EMA12 has crossed above EMA26, indicating a potential uptrend."
	} else if latestEMA12 < latestEMA26 && prevEMA12 >= prevEMA26 {
		trendType = Potential_Downtrend
		result = "Bearish crossover detected. EMA12 has crossed below EMA26, indicating a potential downtrend."
	} else if latestEMA12 > latestEMA26 {
		trendType = Potential_Uptrend
		result = "EMA12 is above EMA26, indicating a potential uptrend."
	} else if latestEMA12 < latestEMA26 {
		trendType = Potential_Downtrend
		result = "EMA12 is below EMA26, indicating a potential downtrend."
	}
	ema.TrendType = trendType
	ema.Result = result
}
