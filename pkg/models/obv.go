package models

import (
	"errors"
	"fmt"
)

type OBV struct {
	symbol    string
	OBVArray  []float64
	TrendType TrendType
	Result    string
}

func NewOBV(symbol string) (*OBV, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	return &OBV{symbol: symbol}, nil
}

func (i *OBV) SetIndex(indexes *Indexes) *Indexes {
	indexes.OBV = *i
	return indexes
}

func (obv *OBV) calculate(marketDataList []BasicMarketData) (bool, error) {
	// Calculate OBV
	obvArray, err := obv.calculateOBV(marketDataList)
	if err != nil {
		return false, fmt.Errorf("Error calculating OBV: %v", err)
	}
	obv.OBVArray = obvArray
	return true, nil
}

// calculateOBV calculates the On-Balance Volume (OBV)
func (obv *OBV) calculateOBV(marketDataList []BasicMarketData) ([]float64, error) {
	closes, _, _, volumes := ExtractMarketData(marketDataList)
	if len(closes) != len(volumes) {
		return nil, fmt.Errorf("length of closes and volumes must be the same")
	}

	obvArray := make([]float64, len(closes))
	obvArray[0] = float64(volumes[0])

	for i := 1; i < len(closes); i++ {
		if closes[i] > closes[i-1] {
			obvArray[i] = obvArray[i-1] + float64(volumes[i])
		} else if closes[i] < closes[i-1] {
			obvArray[i] = obvArray[i-1] - float64(volumes[i])
		} else {
			obvArray[i] = obvArray[i-1]
		}
	}

	return obvArray, nil
}

// analyzeOBV analyzes the OBV data to determine accumulation or distribution
func (obv *OBV) Analyze(marketDataList []BasicMarketData) error {
	status, err := obv.calculate(marketDataList)
	var trendType TrendType
	var result string
	if !status {
		trendType = Neutral
		result = fmt.Sprintf("It is not possible to calculate OBV due to:%s", err)
		obv.TrendType = trendType
		obv.Result = result
		return fmt.Errorf("It is not possible to calculate OBV due to:%s", err)
	}
	if len(obv.OBVArray) < 2 {
		obv.TrendType = Neutral
		obv.Result = "Not enough data to analyze OBV."
		return fmt.Errorf("Not enough data to analyze OBV.")
	}

	isRising := false
	for i := 1; i < len(obv.OBVArray); i++ {
		if obv.OBVArray[i] > obv.OBVArray[i-1] {
			isRising = true
			break
		}
	}
	if obv.OBVArray[len(obv.OBVArray)-1] > obv.OBVArray[0] {
		obv.TrendType = Neutral
		obv.Result = fmt.Sprintf("OBV indicates accumulation with a final value of %.2f. OBV is rising: %v.", obv.OBVArray[len(obv.OBVArray)-1], isRising)
	} else if obv.OBVArray[len(obv.OBVArray)-1] < obv.OBVArray[0] {
		obv.TrendType = Neutral
		obv.Result = fmt.Sprintf("OBV indicates distribution with a final value of %.2f. OBV is rising: %v.", obv.OBVArray[len(obv.OBVArray)-1], isRising)
	} else {
		obv.TrendType = Neutral
		obv.Result = fmt.Sprintf("OBV indicates no significant change with a final value of %.2f. OBV is rising: %v.", obv.OBVArray[len(obv.OBVArray)-1], isRising)
	}
	return nil
}
