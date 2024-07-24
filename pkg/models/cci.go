package models

import (
	"errors"
	"fmt"
	"math"
)

type CCI struct {
	symbol    string
	TrendType TrendType
	Result    string
}

func NewCCI(symbol string) (*CCI, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	return &CCI{symbol: symbol}, nil
}

func (i *CCI) SetIndex(indexes *Indexes) *Indexes {
	indexes.CCI = *i
	return indexes
}

// CalculateCCI calculates the Commodity Channel Index (CCI) for a given period
func (i *CCI) calculateCCI(data []BasicMarketData, period int) ([]float64, error) {
	if len(data) < period {
		return nil, fmt.Errorf("not enough data to calculate CCI for the given period")
	}

	ccis := make([]float64, len(data)-period+1)
	tps := make([]float64, len(data))

	// Calculate the Typical Price (TP)
	for i := 0; i < len(data); i++ {
		tps[i] = (data[i].High + data[i].Low + data[i].Close) / 3
	}

	// Calculate the CCI
	for i := period - 1; i < len(data); i++ {
		// Calculate the Simple Moving Average (SMA) of TP
		sumTP := 0.0
		for j := i - period + 1; j <= i; j++ {
			sumTP += tps[j]
		}
		sma := sumTP / float64(period)

		// Calculate the Mean Deviation (MD)
		sumMD := 0.0
		for j := i - period + 1; j <= i; j++ {
			sumMD += math.Abs(tps[j] - sma)
		}
		md := sumMD / float64(period)

		// Calculate the CCI
		cci := (tps[i] - sma) / (0.015 * md)
		ccis[i-period+1] = cci
	}

	return ccis, nil
}

// AnalyzeCCI analyzes the CCI values and provides a description of the overbought or oversold conditions
func (cci *CCI) Analyze(marketDataList []BasicMarketData) error {
	period := 3
	cciValues, err := cci.calculateCCI(marketDataList, period)
	if err != nil {
		fmt.Printf("Error calculating CCI: %v\n", err)
		cci.TrendType = None
		cci.Result = "None"
		return fmt.Errorf("Error calculating ATR: %v\n", err)
	}
	if len(cciValues) == 0 {
		cci.TrendType = None
		cci.Result = "None"
		return fmt.Errorf("No CCI values provided.")
	}

	lastCCI := cciValues[len(cciValues)-1]

	switch {
	case lastCCI > 100:
		cci.TrendType = Overbought
		cci.Result = fmt.Sprintf("Overbought condition with a CCI of %.2f.", lastCCI)
	case lastCCI < -100:
		cci.TrendType = Oversold
		cci.Result = fmt.Sprintf("Oversold condition with a CCI of %.2f.", lastCCI)
	default:
		cci.TrendType = Neutral
		cci.Result = fmt.Sprintf("Neutral condition with a CCI of %.2f.", lastCCI)
	}
	return nil
}
