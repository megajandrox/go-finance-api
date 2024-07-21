package models

import (
	"fmt"
	"testing"
	"time"
)

func TestTradeHistory(t *testing.T) {
	// Crear datos de mercado
	marketData := MarketData{
		Symbol:    "AAPL",
		Timestamp: time.Now(),
		Open:      150.0,
		High:      155.0,
		Low:       149.0,
		Close:     154.0,
		Volume:    1000000,
		Indicators: Indicators{
			SMA40:      145.0,
			SMA80:      140.0,
			SMA200:     130.0,
			EMA12:      152.0,
			EMA26:      148.0,
			RSI14:      65.0,
			MACD:       4.0,
			MACDSignal: 3.5,
			ATR14:      2.5,
		},
	}

	// Crear una posición de compra
	position := Position{
		Symbol:       "AAPL",
		EntryPrice:   150.0,
		ExitPrice:    154.0,
		Quantity:     100,
		EntryTime:    time.Now().Add(-24 * time.Hour),
		ExitTime:     time.Now(),
		PositionType: Buy,
		MarketType:   Equity,
		Balance:      (154.0 - 150.0) * 100,
		Indicators:   marketData.Indicators,
	}

	// Crear historial de transacciones
	tradeHistory := TradeHistory{}
	tradeHistory.AddPosition(position)

	// Comprobar que la posición se añadió correctamente
	if len(tradeHistory.Positions) != 1 {
		t.Errorf("I was expecting 1 position, but I got %d", len(tradeHistory.Positions))
	}

	// Comprobar que los datos de la posición son correctos
	pos := tradeHistory.Positions[0]
	fmt.Println(pos)
	if pos.Symbol != "AAPL" || pos.EntryPrice != 150.0 || pos.ExitPrice != 154.0 || pos.Quantity != 100 || pos.PositionType != Buy || pos.Balance != 400.0 || pos.MarketType != Equity {
		t.Errorf("The position data is not as expected.")
	}

	// Comprobar que el cálculo de la ganancia/pérdida total es correcto
	if tradeHistory.TotalBalance != 400.0 {
		t.Errorf("I expected a total profit/loss of 400.0, but I got %.2f", tradeHistory.TotalBalance)
	}
}
