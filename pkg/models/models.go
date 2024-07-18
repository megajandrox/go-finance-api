package models

import "time"

type PositionType int

// Definimos las constantes que representarán los valores del enumerador
const (
	Sell PositionType = iota
	Buy
)

type MarketType int

// Definimos las constantes que representarán los valores del enumerador
const (
	Equity MarketType = iota
	ETF
)

/*
 * This struct gonna be persisted for long term
 */
type Position struct {
	Symbol       string       // Símbolo del activo financiero
	EntryPrice   float64      // Precio de entrada
	ExitPrice    float64      // Precio de salida
	Quantity     int          // Cantidad de acciones/compras
	EntryTime    time.Time    // Tiempo de entrada
	ExitTime     time.Time    // Tiempo de salida
	PositionType PositionType // Tipo de posición (compra o venta)
	MarketType   MarketType   // Mercado (Acciones , ETFs)
	Balance      float64      // Ganancia o pérdida
	Indicators   Indicators   // Indicadores técnicos utilizados en la decisión
}

/*
 * This struct could be persisted for short term
 */
type TradeHistory struct {
	Positions    []Position // Lista de todas las posiciones
	TotalBalance float64    // Ganancia o pérdida total
}

// AddPosition agrega una posición al historial y actualiza la ganancia/pérdida total.
func (th *TradeHistory) AddPosition(position Position) {
	th.Positions = append(th.Positions, position)
	th.TotalBalance += position.Balance
}

/*
 * This struct gonna be persisted for long term because has
 */
// MarketData guarda información de mercado obtenida periódicamente.
type MarketData struct {
	Symbol     string     // Símbolo del activo financiero
	Timestamp  time.Time  // Marca de tiempo del dato
	Open       float64    // Precio de apertura
	High       float64    // Precio más alto
	Low        float64    // Precio más bajo
	Close      float64    // Precio de cierre
	Volume     int64      // Volumen negociado
	Indicators Indicators // Indicadores técnicos aplicados
}

/*
 * This struct is not gonna be persisted
 */
// Indicators almacena los valores de los indicadores técnicos aplicados.
type Indicators struct {
	SMA40      float64 // Media Móvil Simple de 40 días
	SMA80      float64 // Media Móvil Simple de 80 días
	SMA200     float64 // Media Móvil Simple de 200 días
	EMA12      float64 // Media Móvil Exponencial de 12 días
	EMA26      float64 // Media Móvil Exponencial de 26 días
	RSI14      float64 // Índice de Fuerza Relativa de 14 días
	MACD       float64 // MACD
	MACDSignal float64 // Señal del MACD
	ATR14      float64 // ATR de 14 días
}
