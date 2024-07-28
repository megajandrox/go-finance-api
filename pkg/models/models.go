package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PositionType int

// Define constants representing the enumerator values
const (
	Bought PositionType = iota
	Sold
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
	gorm.Model
	Symbol       string       // Financial asset symbol
	EntryPrice   float64      // Entry price
	ExitPrice    float64      // Exit price
	Quantity     int          // Quantity of shares/buys
	EntryTime    time.Time    // Entry time
	ExitTime     time.Time    // Exit time
	PositionType PositionType // Position type (buy or sell)
	MarketType   MarketType   // Market (Equities, ETFs)
	Balance      float64      // Profit or loss
}

func NewPosition(symb string, price float64, qty int, marketType MarketType) (*Position, error) {
	if symb == "" {
		return nil, errors.New("symbol cannot be empty")
	}
	return &Position{Symbol: symb, EntryPrice: price, Quantity: qty, MarketType: marketType, EntryTime: time.Now(), PositionType: Bought}, nil
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
	IncreasedTradingRisk
	ModerateTradingRisk
	LowerTradingRisk
	StrongTrend
	WeakTrend
	ModerateTrend
	None
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
	case IncreasedTradingRisk:
		return "Increased Trading Risk"
	case ModerateTradingRisk:
		return "Moderate Trading Risk"
	case LowerTradingRisk:
		return "Lower Trading Risk"
	case None:
		return "None"
	default:
		return "Unknown TrendType"
	}
}

type BasicMarketData struct {
	High      float64
	Low       float64
	Close     float64
	Volume    int64
	TimeStamp int64
}

func ExtractMarketData(data []BasicMarketData) ([]float64, []float64, []float64, []int64) {
	closes := make([]float64, len(data))
	highs := make([]float64, len(data))
	lows := make([]float64, len(data))
	volumes := make([]int64, len(data))

	for i, d := range data {
		closes[i] = d.Close
		highs[i] = d.High
		lows[i] = d.Low
		volumes[i] = d.Volume
	}

	return closes, highs, lows, volumes
}

type Analyzer interface {
	Analyze(marketDataList []BasicMarketData) error
	SetIndex(indexes *Indexes) *Indexes
}

// runAnalysis es una función genérica que ejecuta el análisis utilizando la interfaz Analyzer
func RunAnalysis[T Analyzer](symbol string, marketDataList []BasicMarketData, indexes *Indexes, newAnalyzer func(string) (T, error)) (*Indexes, error) {
	analyzer, err := newAnalyzer(symbol)
	if err != nil {
		return nil, fmt.Errorf("error creating analyzer: %w", err)
	}
	err = analyzer.Analyze(marketDataList)
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
	ATR        ATR
	Momentum   Momentum
	ADX        ADX
	CCI        CCI
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

func NewATRAdapter(symbol string) (Analyzer, error) {
	return NewATR(symbol)
}

func NewMomentumAdapter(symbol string) (Analyzer, error) {
	return NewMomentum(symbol)
}

func NewCCIAdapter(symbol string) (Analyzer, error) {
	return NewCCI(symbol)
}

func NewADXAdapter(symbol string) (Analyzer, error) {
	return NewADX(symbol)
}
