package models

import (
	"errors"
	"fmt"
)

type Momentum struct {
	symbol    string
	TrendType TrendType
	Result    string
}

func NewMomentum(symbol string) (*Momentum, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	return &Momentum{symbol: symbol}, nil
}

func (i *Momentum) SetIndex(indexes *Indexes) *Indexes {
	indexes.Momentum = *i
	return indexes
}

func (m *Momentum) calculateMomentum(inputValues [][]float64, period int) ([]float64, error) {
	if len(inputValues[0]) < period {
		return nil, fmt.Errorf("not enough data to calculate momentum for the given period")
	}

	// Initialize the momentum array
	momentum := make([]float64, len(inputValues[0])-period)

	// Calculate momentum for each subsequent price
	for i := period; i < len(inputValues[0]); i++ {
		momentum[i-period] = inputValues[0][i] - inputValues[0][i-period]
	}

	return momentum, nil
}

// analyze ATR analyzes
func (m *Momentum) Analyze(inputValues [][]float64) error {
	var trendType TrendType
	var result string
	period := 10
	momentum, err := m.calculateMomentum(inputValues, period)
	if err != nil {
		return fmt.Errorf("Error calculating momentum: %v", err)
	}
	if len(momentum) == 0 {
		return fmt.Errorf("No momentum values provided.")
	}
	lastMomentum := momentum[len(momentum)-1]
	switch {
	case lastMomentum > 0:
		trendType = Potential_Uptrend
		result = "The momentum is positive, indicating a potential uptrend."
	case lastMomentum < 0:
		trendType = Potential_Downtrend
		result = "The momentum is negative, indicating a potential downtrend."
	default:
		trendType = Neutral
		result = "The momentum is neutral."
	}
	m.Result = result
	m.TrendType = trendType
	return nil
}
