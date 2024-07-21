package models

import (
	"errors"
	"fmt"
)

type Volume struct {
	symbol    string
	K         []float64
	D         []float64
	TrendType TrendType
	Result    string
}

func NewVolume(symbol string) (*Volume, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	return &Volume{symbol: symbol}, nil
}

func (i *Volume) SetIndex(indexes *Indexes) *Indexes {
	indexes.Volume = *i
	return indexes
}

// analyzeVolumeTrend analyzes the volume data to confirm the strength of a trend
func (v *Volume) Analyze(inputValues [][]float64) error {
	if len(inputValues[3]) < 2 {
		v.TrendType = Neutral
		v.Result = "Not enough data to analyze volume trend."
		return fmt.Errorf("Not enough data to analyze volume trend.")
	}

	totalVolume := 0.0
	increasingDays := 0
	for i := 1; i < len(inputValues[3]); i++ {
		totalVolume += inputValues[3][i]
		if inputValues[3][i] > inputValues[3][i-1] {
			increasingDays++
		}
	}

	averageVolume := totalVolume / float64(len(inputValues[3])-1)
	increasePercentage := float64(increasingDays) / float64(len(inputValues[3])-1) * 100
	v.Result = fmt.Sprintf("Average Volume: %.2f\nPercentage of Increasing Volume Days: %.2f%%", averageVolume, increasePercentage)
	v.TrendType = Neutral
	return nil
}
