package models

import (
	"errors"
	"fmt"
)

type RVOL struct {
	symbol    string
	RVOLValue float64
	TrendType TrendType
	Result    string
}

func NewRVOL(symbol string) (*RVOL, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	return &RVOL{symbol: symbol}, nil
}

func (rvol *RVOL) calculate(volumes []float64) (bool, error) {
	// Calculate RVOL for a 20-day period
	rvolValue, err := rvol.calculateRVOL(volumes, 20)
	if err != nil {
		return false, fmt.Errorf("Error calculating RVOL: %v", err)
	}
	rvol.RVOLValue = rvolValue
	return true, nil
}

// calculateRVOL calculates the Relative Volume (RVOL)
func (rvol *RVOL) calculateRVOL(volumes []float64, period int) (float64, error) {
	if len(volumes) < period {
		return 0, fmt.Errorf("not enough data to calculate RVOL for the given period")
	}

	// Calculate the average volume over the period
	totalVolume := 0.0
	for i := len(volumes) - period; i < len(volumes); i++ {
		totalVolume += volumes[i]
	}
	averageVolume := totalVolume / float64(period)

	// Get the most recent volume
	//TODO obtener el ultimo distinto de cero.
	currentVolume := volumes[len(volumes)-2]

	// Calculate the relative volume
	rvolValue := currentVolume / averageVolume

	return rvolValue, nil
}

// analyzeRVOL analyzes the RVOL value and returns a descriptive analysis
func (rvol *RVOL) AnalyzeRVOL(volumes []float64) {
	status, err := rvol.calculate(volumes)
	rvol.TrendType = Neutral
	var result string
	if !status {
		rvol.Result = fmt.Sprintf("No es posible calcular RVOL debido a: %s", err)
		return
	}
	if rvol.RVOLValue > 1.0 {
		result = fmt.Sprintf("RVOL is %.2f, indicating a higher than average volume.", rvol.RVOLValue)
	} else if rvol.RVOLValue < 1.0 {
		result = fmt.Sprintf("RVOL is %.2f, indicating a lower than average volume.", rvol.RVOLValue)
	} else {
		result = fmt.Sprintf("RVOL is %.2f, indicating an average volume.", rvol.RVOLValue)
	}
	rvol.Result = result
}
