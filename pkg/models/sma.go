package models

import (
	"errors"
	"fmt"
)

type SMA struct {
	symbol    string
	SMA40     float64
	SMA80     float64
	SMA200    float64
	TrendType TrendType
	Result    string
}

func NewSMA(symbol string) (*SMA, error) {
	if symbol == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	return &SMA{symbol: symbol}, nil
}

func (sma *SMA) calculate(closes []float64) (bool, error) {
	sma40, err := sma.calculateSMAN(closes, 40)
	if err != nil {
		return false, err
	}
	sma.SMA40 = sma40
	sma80, err := sma.calculateSMAN(closes, 80)
	if err != nil {
		return false, err
	}
	sma.SMA80 = sma80
	sma200, err := sma.calculateSMAN(closes, 200)
	if err != nil {
		return false, err
	}
	sma.SMA200 = sma200
	return true, nil
}

func (sma *SMA) calculateSMAN(closes []float64, n int) (float64, error) {
	// Calculate SMA for the last N days
	if len(closes) < n {
		return 0, fmt.Errorf("not enough data to calculate SMA%d", n)
	}
	sum := 0.0
	for i := len(closes) - n; i < len(closes); i++ {
		sum += closes[i]
	}
	return sum / float64(n), nil
}

// analyzeSMATrend analyzes the relationship between the SMAs to determine if it's an uptrend
func (sma *SMA) AnalyzeSMATrend(closes []float64) {
	status, err := sma.calculate(closes)
	var trendType TrendType = Neutral
	var result string = "Las SMAs no están en un orden claro para confirmar una tendencia específica."
	if !status {
		trendType = Neutral
		result = fmt.Sprintf("No es posible calcular SMA debido a: %s", err)
		sma.TrendType = trendType
		sma.Result = result
		return
	}
	if sma.SMA40 > sma.SMA80 && sma.SMA80 > sma.SMA200 {
		trendType = Uptrend
		result = "SMA 40 > SMA 80 > SMA 200: Esta relación sugiere que el precio de la acción está en una tendencia alcista. Las SMAs más cortas (40 días) están por encima de las SMAs más largas (80 y 200 días), lo que indica que los precios recientes están más altos que los precios pasados."
	} else if sma.SMA40 < sma.SMA80 && sma.SMA80 < sma.SMA200 {
		trendType = Downtrend
		result = "SMA 40 < SMA 80 < SMA 200: Esta relación sugiere que el precio de la acción está en una tendencia bajista. Las SMAs más cortas (40 días) están por debajo de las SMAs más largas (80 y 200 días), lo que indica que los precios recientes están más bajos que los precios pasados."
	} else if sma.SMA40 > sma.SMA80 && sma.SMA80 < sma.SMA200 {
		trendType = Shortterm_Uptrend_Longterm_Downtrend
		result = "SMA 40 > SMA 80 < SMA 200: Esta relación sugiere que el precio de la acción puede estar en una recuperación a corto plazo, pero sigue en una tendencia bajista a largo plazo."
	} else if sma.SMA40 < sma.SMA80 && sma.SMA80 > sma.SMA200 {
		trendType = Shortterm_Downtrend_Longterm_Uptrend
		result = "SMA 40 < SMA 80 > SMA 200: Esta relación sugiere que el precio de la acción puede estar en una corrección a corto plazo, pero sigue en una tendencia alcista a largo plazo."
	}
	sma.TrendType = trendType
	sma.Result = result
}
