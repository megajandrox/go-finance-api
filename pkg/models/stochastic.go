package models

import (
	"errors"
	"fmt"
)

type Stochastic struct {
	symbol    string
	K         []float64
	D         []float64
	TrendType TrendType
	Result    string
}

func NewStochastic(symbol string) (*Stochastic, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	return &Stochastic{symbol: symbol}, nil
}

func (i *Stochastic) SetIndex(indexes *Indexes) *Indexes {
	indexes.Stochastic = *i
	return indexes
}

func (sto *Stochastic) calculate(marketDataList []BasicMarketData) (bool, error) {
	// Calculate Stochastic Oscillator for 14-day period
	k, d, err := sto.calculateStochasticOscillator(marketDataList, 14)
	if err != nil {
		return false, fmt.Errorf(`Error calculating Stochastic Oscillator: %v`, err)
	}
	sto.K = k
	sto.D = d
	return true, nil
}

// calculateStochasticOscillator calculates the Stochastic Oscillator
func (sto *Stochastic) calculateStochasticOscillator(marketDataList []BasicMarketData, period int) ([]float64, []float64, error) {
	closes, highs, lows, _ := ExtractMarketData(marketDataList)
	if len(closes) < period || len(highs) < period || len(lows) < period {
		return nil, nil, fmt.Errorf("not enough data to calculate Stochastic Oscillator for the given period")
	}

	k := make([]float64, len(closes)-period+1)
	for i := period - 1; i < len(closes); i++ {
		low := lows[i-period+1]
		high := highs[i-period+1]
		for j := i - period + 1; j <= i; j++ {
			if lows[j] < low {
				low = lows[j]
			}
			if highs[j] > high {
				high = highs[j]
			}
		}
		k[i-period+1] = 100 * ((closes[i] - low) / (high - low))
	}

	// Calculate %D as a 3-period SMA of %K
	d := make([]float64, len(k))
	for i := 2; i < len(k); i++ {
		d[i] = (k[i] + k[i-1] + k[i-2]) / 3
	}

	return k, d, nil
}

// analyzeStochasticOscillator analyzes the Stochastic Oscillator values
func (sto *Stochastic) Analyze(marketDataList []BasicMarketData) error {
	status, err := sto.calculate(marketDataList)
	latestK := sto.K[len(sto.K)-1]
	latestD := sto.D[len(sto.D)-1]
	var trendType TrendType = Neutral
	var result string = fmt.Sprintf("Stochastic Oscillator is %.2f/%.2f, indicating normal market conditions.", latestK, latestD)
	if !status {
		trendType = Neutral
		result = fmt.Sprintf("It is not possible to calculate SMA because: %s", err)
		sto.TrendType = trendType
		sto.Result = result
		return fmt.Errorf("It is not possible to calculate SMA because: %s", err)
	}
	if len(sto.K) == 0 || len(sto.D) == 0 {
		sto.TrendType = Neutral
		sto.Result = "Not enough data for Stochastic Oscillator analysis."
		return fmt.Errorf("Not enough data for Stochastic Oscillator analysis.")
	}

	if latestK > 80 && latestD > 80 {
		trendType = Overbought
		result = fmt.Sprintf("Stochastic Oscillator is %.2f/%.2f, indicating the asset is overbought.", latestK, latestD)
	} else if latestK < 20 && latestD < 20 {
		trendType = Oversold
		result = fmt.Sprintf("Stochastic Oscillator is %.2f/%.2f, indicating the asset is oversold.", latestK, latestD)
	}
	sto.TrendType = trendType
	sto.Result = result
	return nil
}
