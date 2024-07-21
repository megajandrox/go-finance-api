package models

import (
	"errors"
	"fmt"
)

type MACD struct {
	symbol     string
	MACDArray  []float64
	MACDSignal []float64
	TrendType  TrendType
	Result     string
}

func NewMACD(symbol string) (*MACD, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	return &MACD{symbol: symbol}, nil
}

func (i *MACD) SetIndex(indexes *Indexes) *Indexes {
	indexes.MACD = *i
	return indexes
}

// calculateMACD calculates the MACD and Signal line
func (macd *MACD) calculate(closes []float64) (bool, error) {
	// Calculate EMA for 12-day and 26-day periods
	ema, errEMA := NewEMA(macd.symbol)
	if errEMA != nil {
		return false, fmt.Errorf("Error creating EMA: %v", errEMA)
	}
	var inputValues [][]float64
	inputValues = append(inputValues, closes)
	ema.Analyze(inputValues)
	var ema12, ema26 []float64
	ema12 = ema.EMA12
	ema26 = ema.EMA26
	// Ensure we have enough data for the MACD calculation
	minLength := min(len(ema12), len(ema26))
	ema12 = ema12[len(ema12)-minLength:]
	ema26 = ema26[len(ema26)-minLength:]

	macdArray := make([]float64, minLength)
	for i := range macdArray {
		macdArray[i] = ema12[i] - ema26[i]
	}

	signal, err := ema.CalculateEMA(macdArray, 9)
	if err != nil {
		return false, fmt.Errorf("error calculating Signal line: %v", err)
	}
	macd.MACDArray = macdArray
	macd.MACDSignal = signal
	return true, nil
}

// analyzeMACD analyzes the MACD and Signal line for crossovers
func (macd *MACD) Analyze(inputValues [][]float64) error {
	var macdArray, signal []float64
	var trendType TrendType = Neutral
	var result string = "The SMAs are not in a clear order to confirm a specific trend."
	status, err := macd.calculate(inputValues[0])
	if !status {
		trendType = Neutral
		result = fmt.Sprintf("It is not possible to calculate MACD due to: %s", err)
		macd.TrendType = trendType
		macd.Result = result
		return fmt.Errorf("It is not possible to calculate MACD due to: %s", err)
	}
	macdArray = macd.MACDArray
	signal = macd.MACDSignal
	if len(macdArray) == 0 || len(signal) == 0 {
		macd.TrendType = Neutral
		macd.Result = "Not enough data for MACD analysis."
		return fmt.Errorf("Not enough data for MACD analysis.")
	}
	// Determine the most recent values
	latestMACD := macdArray[len(macdArray)-1]
	latestSignal := signal[len(signal)-1]

	// Determine the previous values
	prevMACD := macdArray[len(macdArray)-2]
	prevSignal := signal[len(signal)-2]

	// Analyze the crossover
	macd.TrendType = Neutral
	macd.Result = "No significant crossover detected."
	if latestMACD > latestSignal && prevMACD <= prevSignal {
		macd.TrendType = Potential_Uptrend
		macd.Result = "Bullish crossover detected. MACD has crossed above the Signal line, indicating a potential uptrend."
	} else if latestMACD < latestSignal && prevMACD >= prevSignal {
		macd.TrendType = Potential_Downtrend
		macd.Result = "Bearish crossover detected. MACD has crossed below the Signal line, indicating a potential downtrend."
	} else if latestMACD > latestSignal {
		macd.TrendType = Potential_Uptrend
		macd.Result = "MACD is above the Signal line, indicating a potential uptrend."
	} else if latestMACD < latestSignal {
		macd.TrendType = Potential_Downtrend
		macd.Result = "MACD is below the Signal line, indicating a potential downtrend."
	}
	return nil
}
