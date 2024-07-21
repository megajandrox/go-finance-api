package models

import (
	"fmt"
	"time"
)

type PositionType int

// Define constants representing the enumerator values
const (
	Sell PositionType = iota
	Buy
)

type MarketType int

// Define constants representing the enumerator values
const (
	Equity MarketType = iota
	ETF
)

/*
 * This struct is going to be persisted for the long term
 */
type Position struct {
	Symbol       string       // Financial asset symbol
	EntryPrice   float64      // Entry price
	ExitPrice    float64      // Exit price
	Quantity     int          // Quantity of shares/buys
	EntryTime    time.Time    // Entry time
	ExitTime     time.Time    // Exit time
	PositionType PositionType // Position type (buy or sell)
	MarketType   MarketType   // Market (Equities, ETFs)
	Balance      float64      // Profit or loss
	Indicators   Indicators   // Technical indicators used in the decision
}

/*
 * This struct could be persisted for the short term
 */
type TradeHistory struct {
	Positions    []Position // List of all positions
	TotalBalance float64    // Total profit or loss
}

// AddPosition adds a position to the history and updates the total profit/loss.
func (th *TradeHistory) AddPosition(position Position) {
	th.Positions = append(th.Positions, position)
	th.TotalBalance += position.Balance
}

/*
 * This struct is going to be persisted for the long term because it contains
 */
// MarketData stores market information obtained periodically.
type MarketData struct {
	Symbol     string     // Financial asset symbol
	Timestamp  time.Time  // Data timestamp
	Open       float64    // Opening price
	High       float64    // Highest price
	Low        float64    // Lowest price
	Close      float64    // Closing price
	Volume     int64      // Traded volume
	Indicators Indicators // Applied technical indicators
}

/*
 * This struct is not going to be persisted
 */
// Indicators stores the values of the applied technical indicators.
type Indicators struct {
	SMA40      float64 // 40-day Simple Moving Average
	SMA80      float64 // 80-day Simple Moving Average
	SMA200     float64 // 200-day Simple Moving Average
	EMA12      float64 // 12-day Exponential Moving Average
	EMA26      float64 // 26-day Exponential Moving Average
	RSI14      float64 // 14-day Relative Strength Index
	MACD       float64 // MACD
	MACDSignal float64 // MACD Signal
	ATR14      float64 // 14-day ATR
}

type TrendType int

// Define constants representing the enumerator values
const (
	Uptrend TrendType = iota
	Potential_Uptrend
	Downtrend
	Potential_Downtrend
	Shortterm_Uptrend_Longterm_Downtrend
	Shortterm_Downtrend_Longterm_Uptrend
	Overbought
	Oversold
	Neutral
)

func (t TrendType) String() string {
	switch t {
	case Uptrend:
		return "Uptrend"
	case Downtrend:
		return "Downtrend"
	case Potential_Uptrend:
		return "Potential Uptrend"
	case Potential_Downtrend:
		return "Potential Downtrend"
	case Shortterm_Uptrend_Longterm_Downtrend:
		return "Short-term Uptrend, Long-term Downtrend"
	case Shortterm_Downtrend_Longterm_Uptrend:
		return "Short-term Downtrend, Long-term Uptrend"
	case Overbought:
		return "Overbought"
	case Oversold:
		return "Oversold"
	case Neutral:
		return "Neutral"
	default:
		return "Unknown TrendType"
	}
}

type Analyzer interface {
	Analyze(inputValues [][]float64) error
	SetIndex(indexes *Indexes) *Indexes
}

// runAnalysis es una función genérica que ejecuta el análisis utilizando la interfaz Analyzer
func RunAnalysis[T Analyzer](symbol string, inputValues [][]float64, indexes *Indexes, newAnalyzer func(string) (T, error)) (*Indexes, error) {
	analyzer, err := newAnalyzer(symbol)
	if err != nil {
		return nil, fmt.Errorf("error creating analyzer: %w", err)
	}
	err = analyzer.Analyze(inputValues)
	if err != nil {
		return nil, fmt.Errorf("error analyzing trend: %w", err)
	}
	ind := analyzer.SetIndex(indexes)
	return ind, nil
}

type Indexes struct {
	symbol     string
	SMA        SMA
	EMA        EMA
	MACD       MACD
	RSI        RSI
	Stochastic Stochastic
	Volume     Volume
	OBV        OBV
	RVOL       RVOL
}

func NewIndexes(symbol string) *Indexes {
	return &Indexes{symbol: symbol}
}

// Adapter functions to convert specific analyzers to Analyzer interface
func NewSMAAdapter(symbol string) (Analyzer, error) {
	return NewSMA(symbol)
}

// Define similar adapter functions for other analyzers
func NewEMAAdapter(symbol string) (Analyzer, error) {
	return NewEMA(symbol)
}

func NewMACDAdapter(symbol string) (Analyzer, error) {
	return NewMACD(symbol)
}

func NewRSIAdapter(symbol string) (Analyzer, error) {
	return NewRSI(symbol)
}

func NewStochasticAdapter(symbol string) (Analyzer, error) {
	return NewStochastic(symbol)
}

func NewVolumeAdapter(symbol string) (Analyzer, error) {
	return NewVolume(symbol)
}

func NewOBVAdapter(symbol string) (Analyzer, error) {
	return NewOBV(symbol)
}

func NewRVOLAdapter(symbol string) (Analyzer, error) {
	return NewRVOL(symbol)
}
