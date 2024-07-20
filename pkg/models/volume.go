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

// analyzeVolumeTrend analyzes the volume data to confirm the strength of a trend
func (v *Volume) AnalyzeVolumeTrend(volumes []float64) {
	if len(volumes) < 2 {
		v.TrendType = Neutral
		v.Result = "Not enough data to analyze volume trend."
		return
	}

	totalVolume := 0.0
	increasingDays := 0
	for i := 1; i < len(volumes); i++ {
		totalVolume += volumes[i]
		if volumes[i] > volumes[i-1] {
			increasingDays++
		}
	}

	averageVolume := totalVolume / float64(len(volumes)-1)
	increasePercentage := float64(increasingDays) / float64(len(volumes)-1) * 100
	v.Result = fmt.Sprintf("Average Volume: %.2f\nPercentage of Increasing Volume Days: %.2f%%", averageVolume, increasePercentage)
	v.TrendType = Neutral
}
