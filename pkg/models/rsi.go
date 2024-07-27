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

func (i *RSI) SetIndex(indexes *Indexes) *Indexes {
	indexes.RSI = *i
	return indexes
}

func (rsi *RSI) calculate(marketDataList []BasicMarketData) (bool, error) {
	// Calculate RSI for 14-day period
	dailyCloses := ExtractDailyCloses(marketDataList)
	rsiArray, err := rsi.calculateRSI(dailyCloses, 14)
	if err != nil {
		return false, fmt.Errorf(`Error calculating RSI: %v`, err)
	}
	// Get the last RSI value
	rsi.LatestRSI = rsiArray[len(rsiArray)-1]
	return true, nil
}

// CalculateRSI calculates the RSI for a given period
func (rsi *RSI) calculateRSI(marketDataList []BasicMarketData, period int) ([]float64, error) {
	if len(marketDataList) < period {
		return nil, fmt.Errorf("not enough data to calculate RSI for the given period")
	}

	// Initialize gains and losses
	gains := make([]float64, len(marketDataList))
	losses := make([]float64, len(marketDataList))

	// Calculate daily gains and losses
	for i := 1; i < len(marketDataList); i++ {
		change := marketDataList[i].Close - marketDataList[i-1].Close
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
	rsiArray := make([]float64, len(marketDataList))

	// Calculate RSI for the first period
	rs := avgGain / avgLoss
	rsiArray[period-1] = 100 - (100 / (1 + rs))

	// Calculate RSI for subsequent periods
	for i := period; i < len(marketDataList); i++ {
		avgGain = (avgGain*float64(period-1) + gains[i]) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + losses[i]) / float64(period)
		rs = avgGain / avgLoss
		rsiArray[i] = 100 - (100 / (1 + rs))
	}

	return rsiArray[period-1:], nil // Return only the valid RSI values
}

// analyzeRSI analyzes the RSI value and returns a descriptive analysis
func (rsi *RSI) Analyze(marketDataList []BasicMarketData) error {
	status, err := rsi.calculate(marketDataList)
	var trendType TrendType = Neutral
	var result string = "The RSIs are not in a clear order to confirm a specific trend."
	if !status {
		trendType = Neutral
		result = fmt.Sprintf("It is not possible to calculate RSI due to: %s", err)
		rsi.TrendType = trendType
		rsi.Result = result
		return fmt.Errorf("It is not possible to calculate RSI due to: %s", err)
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
	return nil
}
