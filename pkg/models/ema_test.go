package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEMA(t *testing.T) {
	ema, err := NewEMA("AAPL")
	assert.NoError(t, err)
	assert.NotNil(t, ema)
	assert.Equal(t, "AAPL", ema.symbol)

	_, err = NewEMA("")
	assert.Error(t, err)
	assert.Equal(t, "symbol cannot be empty", err.Error())
}

func TestCalculateEMA(t *testing.T) {
	ema, err := NewEMA("AAPL")
	assert.NoError(t, err)
	assert.NotNil(t, ema)

	prices := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	period := 5
	expected := []float64{3, 4, 5, 6, 7, 8, 9, 10}
	result, err := ema.CalculateEMA(prices, period)
	fmt.Println(result)
	assert.NoError(t, err)
	assert.InDeltaSlice(t, expected, result, 1e-9)

	_, err = ema.CalculateEMA(prices, 20)
	assert.Error(t, err)
	assert.Equal(t, "not enough data to calculate EMA for the given period", err.Error())
}

func TestAnalyzeEMACrossover(t *testing.T) {
	ema, err := NewEMA("AAPL")
	assert.NoError(t, err)
	assert.NotNil(t, ema)

	// Simulate some closing prices
	inputValues := [][]float64{
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26},
	}

	err = ema.Analyze(inputValues)
	assert.Error(t, err)
	assert.Equal(t, Neutral, ema.TrendType)
	assert.Equal(t, "Not enough data for analysis.", ema.Result)
}

func TestAnalyzeEMACrossoverNotEnoughData(t *testing.T) {
	ema, err := NewEMA("AAPL")
	assert.NoError(t, err)
	assert.NotNil(t, ema)

	// Simulate some closing prices
	inputValues := [][]float64{
		{1, 2, 3, 4},
	}

	err = ema.Analyze(inputValues)
	assert.Error(t, err)
	assert.Equal(t, Neutral, ema.TrendType)
	assert.Equal(t, "It is not possible to calculate EMA due to: Error calculating EMA 12: not enough data to calculate EMA for the given period", ema.Result)
}

func TestAnalyzeEMACrossoverInvalidData(t *testing.T) {
	ema, err := NewEMA("AAPL")
	assert.NoError(t, err)
	assert.NotNil(t, ema)

	// Simulate invalid data
	inputValues := [][]float64{}

	err = ema.Analyze(inputValues)
	assert.Error(t, err)
	assert.Equal(t, Neutral, ema.TrendType)
	assert.Equal(t, "Not enough data for analysis.", ema.Result)
}
