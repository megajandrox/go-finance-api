package models

import (
	"errors"
	"fmt"
	"math"
)

type ATR struct {
	symbol    string
	TrendType TrendType
	Result    string
}

func NewATR(symbol string) (*ATR, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	return &ATR{symbol: symbol}, nil
}

func (i *ATR) SetIndex(indexes *Indexes) *Indexes {
	indexes.ATR = *i
	return indexes
}

// CalculateTR calculates the True Range (TR) for a given day
func (atr *ATR) calculateTR(high, low, prevClose float64) float64 {
	tr1 := high - low
	tr2 := math.Abs(high - prevClose)
	tr3 := math.Abs(low - prevClose)
	return math.Max(tr1, math.Max(tr2, tr3))
}

// CalculateATR calculates the Average True Range (ATR) for a given period
func (atr *ATR) calculateATR(marketDataList []BasicMarketData, period int) ([]float64, error) {
	if len(marketDataList) < period {
		return nil, fmt.Errorf("not enough data to calculate ATR for the given period")
	}

	// Initialize TR array
	trArray := make([]float64, len(marketDataList))

	// Calculate TR for each day
	for i := 1; i < len(marketDataList); i++ {
		trArray[i] = atr.calculateTR(marketDataList[i].High, marketDataList[i].Low, marketDataList[i-1].Close)
	}

	// Initialize ATR array
	atrArray := make([]float64, len(marketDataList)-(period-1))

	// Calculate the initial ATR using the average of the first 'period' TR values
	sumTR := 0.0
	for i := 1; i <= period; i++ {
		sumTR += trArray[i]
	}
	atrArray[0] = sumTR / float64(period)

	// Calculate the ATR for each subsequent day
	for i := period + 1; i < len(marketDataList); i++ {
		atrArray[i-period] = ((atrArray[i-period-1] * float64(period-1)) + trArray[i]) / float64(period)
	}

	return atrArray, nil
}

// analyze ATR analyzes
func (atr *ATR) Analyze(marketDataList []BasicMarketData) error {
	var trendType TrendType
	var result string
	period := 30
	atrArray, err := atr.calculateATR(marketDataList, period)
	if err != nil {
		fmt.Printf("Error calculating ATR: %v\n", err)
		atr.TrendType = None
		atr.Result = "None"
		return fmt.Errorf("Error calculating ATR: %v\n", err)
	}
	if len(atrArray) == 0 {
		atr.TrendType = None
		atr.Result = "None"
		return fmt.Errorf("No ATR values provided.")
	}

	// Calculate the average ATR
	sumATR := 0.0
	for _, value := range atrArray {
		sumATR += value
	}
	averageATR := sumATR / float64(len(atrArray))

	// Determine volatility and risk descriptions based on the average ATR
	switch {
	case averageATR > 2.0:
		trendType = IncreasedTradingRisk
		result = fmt.Sprintf("High volatility with an average ATR of %.2f. Increased trading risk.", averageATR)
	case averageATR > 1.0:
		trendType = ModerateTradingRisk
		result = fmt.Sprintf("Moderate volatility with an average ATR of %.2f. Moderate trading risk.", averageATR)
	default:
		trendType = LowerTradingRisk
		result = fmt.Sprintf("Low volatility with an average ATR of %.2f. Lower trading risk.", averageATR)
	}
	atr.TrendType = trendType
	atr.Result = result
	return nil
}
