package models

import (
	"errors"
	"fmt"
	"math"
)

type ADR struct {
	symbol    string
	TrendType TrendType
	Result    string
}

func NewADR(symbol string) (*ADR, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	return &ADR{symbol: symbol}, nil
}

func (i *ADR) SetIndex(indexes *Indexes) *Indexes {
	indexes.ADR = *i
	return indexes
}

// CalculateTR calculates the True Range (TR) for a given day
func (adr *ADR) calculateTR(current, previous BasicMarketData) float64 {
	tr1 := current.High - current.Low
	tr2 := math.Abs(current.High - previous.Close)
	tr3 := math.Abs(current.Low - previous.Close)
	return math.Max(tr1, math.Max(tr2, tr3))
}

// CalculateDM calculates the +DM and -DM for a given day
func (adr *ADR) calculateDM(current, previous BasicMarketData) (float64, float64) {
	upMove := current.High - previous.High
	downMove := previous.Low - current.Low

	positiveDM := 0.0
	negativeDM := 0.0

	if upMove > downMove && upMove > 0 {
		positiveDM = upMove
	}
	if downMove > upMove && downMove > 0 {
		negativeDM = downMove
	}

	return positiveDM, negativeDM
}

// CalculateADX calculates the ADX for a given period
func (adr *ADR) calculateADX(marketDataList []BasicMarketData, period int) ([]float64, error) {
	if len(marketDataList) < period+1 {
		return nil, fmt.Errorf("not enough data to calculate ADX for the given period")
	}

	// Initialize arrays
	tr := make([]float64, len(marketDataList))
	positiveDM := make([]float64, len(marketDataList))
	negativeDM := make([]float64, len(marketDataList))
	positiveDI := make([]float64, len(marketDataList))
	negativeDI := make([]float64, len(marketDataList))
	dx := make([]float64, len(marketDataList))
	adx := make([]float64, len(marketDataList)-period)

	// Calculate TR, +DM, -DM
	for i := 1; i < len(marketDataList); i++ {
		tr[i] = adr.calculateTR(marketDataList[i], marketDataList[i-1])
		positiveDM[i], negativeDM[i] = adr.calculateDM(marketDataList[i], marketDataList[i-1])
	}

	// Smooth +DM, -DM, TR
	trSmooth := tr[1]
	positiveDMSmooth := positiveDM[1]
	negativeDMSmooth := negativeDM[1]
	for i := 2; i <= period; i++ {
		trSmooth += tr[i]
		positiveDMSmooth += positiveDM[i]
		negativeDMSmooth += negativeDM[i]
	}

	// Calculate +DI, -DI
	for i := period + 1; i < len(marketDataList); i++ {
		trSmooth = trSmooth - (trSmooth / float64(period)) + tr[i]
		positiveDMSmooth = positiveDMSmooth - (positiveDMSmooth / float64(period)) + positiveDM[i]
		negativeDMSmooth = negativeDMSmooth - (negativeDMSmooth / float64(period)) + negativeDM[i]

		positiveDI[i] = 100 * (positiveDMSmooth / trSmooth)
		negativeDI[i] = 100 * (negativeDMSmooth / trSmooth)

		dx[i] = 100 * math.Abs(positiveDI[i]-negativeDI[i]) / (positiveDI[i] + negativeDI[i])
	}

	// Smooth DX to get ADX
	dxSmooth := 0.0
	for i := period + 1; i <= 2*period; i++ {
		dxSmooth += dx[i]
	}
	adx[0] = dxSmooth / float64(period)

	for i := 2*period + 1; i < len(marketDataList); i++ {
		dxSmooth = (dxSmooth*(float64(period)-1) + dx[i]) / float64(period)
		adx[i-period] = dxSmooth
	}

	return adx, nil
}

// AnalyzeADX analyzes the ADX values and provides a description of the trend strength
func (adr *ADR) Analyze(marketDataList []BasicMarketData) error {
	var trendType TrendType
	var result string
	period := 3
	adxValues, err := adr.calculateADX(marketDataList, period)
	if err != nil {
		adr.TrendType = None
		adr.Result = "None"
		return fmt.Errorf("Error calculating ADX: %v\n", err)
	}
	if len(adxValues) == 0 {
		adr.TrendType = None
		adr.Result = "None"
		return fmt.Errorf("No ADX values provided.")
	}

	lastADX := adxValues[len(adxValues)-1]

	switch {
	case lastADX > 25:
		trendType = StrongTrend
		result = fmt.Sprintf("Strong trend with an ADX of %.2f.", lastADX)
	case lastADX > 20:
		trendType = ModerateTrend
		result = fmt.Sprintf("Moderate trend with an ADX of %.2f.", lastADX)
	default:
		trendType = WeakTrend
		result = fmt.Sprintf("Weak or no trend with an ADX of %.2f.", lastADX)
	}
	adr.TrendType = trendType
	adr.Result = result
	return nil
}
