package utils

import (
	"fmt"
	"time"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

type SMAs struct {
	SMA40  float64
	SMA80  float64
	SMA200 float64
}

func CalculateSMAN(closes []float64, n int) (float64, error) {
	// Calculate SMA for the last N days
	if len(closes) < n {
		return 0, fmt.Errorf("not enough data to calculate SMA%d", n)
	}

	sum := 0.0
	for i := len(closes) - n; i < len(closes); i++ {
		sum += closes[i]
	}
	sma := sum / float64(n)
	return sma, nil
}

// calculateEMA calculates the Exponential Moving Average (EMA)
func CalculateEMA(prices []float64, period int) ([]float64, error) {
	if len(prices) < period {
		return nil, fmt.Errorf("not enough data to calculate EMA for the given period")
	}

	// Initialize the EMA array
	ema := make([]float64, len(prices))

	// Calculate the initial SMA to start the EMA
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema[period-1] = sum / float64(period)

	// Calculate the smoothing factor
	alpha := 2.0 / float64(period+1)

	// Calculate the EMA for each subsequent price
	for i := period; i < len(prices); i++ {
		ema[i] = (prices[i] * alpha) + (ema[i-1] * (1 - alpha))
	}

	return ema[period-1:], nil // Return only the valid EMA values
}

// calculateMACD calculates the MACD and Signal line
func CalculateMACD(prices []float64) ([]float64, []float64, error) {
	ema12, err := CalculateEMA(prices, 12)
	if err != nil {
		return nil, nil, fmt.Errorf("error calculating EMA 12: %v", err)
	}

	ema26, err := CalculateEMA(prices, 26)
	if err != nil {
		return nil, nil, fmt.Errorf("error calculating EMA 26: %v", err)
	}

	// Ensure we have enough data for the MACD calculation
	minLength := min(len(ema12), len(ema26))
	ema12 = ema12[len(ema12)-minLength:]
	ema26 = ema26[len(ema26)-minLength:]

	macd := make([]float64, minLength)
	for i := range macd {
		macd[i] = ema12[i] - ema26[i]
	}

	signal, err := CalculateEMA(macd, 9)
	if err != nil {
		return nil, nil, fmt.Errorf("error calculating Signal line: %v", err)
	}

	return macd, signal, nil
}

// analyzeSMATrend analyzes the relationship between the SMAs to determine if it's an uptrend
func AnalyzeSMATrend(smas SMAs) (string, string) {
	if smas.SMA40 > smas.SMA80 && smas.SMA80 > smas.SMA200 {
		return "Uptrend", "SMA 40 > SMA 80 > SMA 200: Esta relación sugiere que el precio de la acción está en una tendencia alcista. Las SMAs más cortas (40 días) están por encima de las SMAs más largas (80 y 200 días), lo que indica que los precios recientes están más altos que los precios pasados."
	} else if smas.SMA40 < smas.SMA80 && smas.SMA80 < smas.SMA200 {
		return "Downtrend", "SMA 40 < SMA 80 < SMA 200: Esta relación sugiere que el precio de la acción está en una tendencia bajista. Las SMAs más cortas (40 días) están por debajo de las SMAs más largas (80 y 200 días), lo que indica que los precios recientes están más bajos que los precios pasados."
	} else if smas.SMA40 > smas.SMA80 && smas.SMA80 < smas.SMA200 {
		return "Shor-term Uptrend Long-term Downtrend", "SMA 40 > SMA 80 < SMA 200: Esta relación sugiere que el precio de la acción puede estar en una recuperación a corto plazo, pero sigue en una tendencia bajista a largo plazo."
	} else if smas.SMA40 < smas.SMA80 && smas.SMA80 > smas.SMA200 {
		return "Short term Downtrend Long-term Uptrend", "SMA 40 < SMA 80 > SMA 200: Esta relación sugiere que el precio de la acción puede estar en una corrección a corto plazo, pero sigue en una tendencia alcista a largo plazo."
	}
	return "Neutral", "Las SMAs no están en un orden claro para confirmar una tendencia específica."
}

// analyzeMACD analyzes the MACD and Signal line for crossovers
func AnalyzeMACD(macd, signal []float64) string {
	if len(macd) == 0 || len(signal) == 0 {
		return "Not enough data for MACD analysis."
	}

	// Determine the most recent values
	latestMACD := macd[len(macd)-1]
	latestSignal := signal[len(signal)-1]

	// Determine the previous values
	prevMACD := macd[len(macd)-2]
	prevSignal := signal[len(signal)-2]

	// Analyze the crossover
	if latestMACD > latestSignal && prevMACD <= prevSignal {
		return "Bullish crossover detected. MACD has crossed above the Signal line, indicating a potential uptrend."
	} else if latestMACD < latestSignal && prevMACD >= prevSignal {
		return "Bearish crossover detected. MACD has crossed below the Signal line, indicating a potential downtrend."
	} else if latestMACD > latestSignal {
		return "MACD is above the Signal line, indicating a potential uptrend."
	} else if latestMACD < latestSignal {
		return "MACD is below the Signal line, indicating a potential downtrend."
	}

	return "No significant crossover detected."
}

// analyzeEMACrossover analyzes the crossover between EMA12 and EMA26
func AnalyzeEMACrossover(ema12, ema26 []float64) (string, string) {
	if len(ema12) == 0 || len(ema26) == 0 {
		return "Neutral", "Not enough data for analysis."
	}

	// Determine the most recent values
	latestEMA12 := ema12[len(ema12)-1]
	latestEMA26 := ema26[len(ema26)-1]

	// Determine the previous values
	prevEMA12 := ema12[len(ema12)-2]
	prevEMA26 := ema26[len(ema26)-2]

	// Analyze the crossover
	if latestEMA12 > latestEMA26 && prevEMA12 <= prevEMA26 {
		return "Potential uptrend", "Bullish crossover detected. EMA12 has crossed above EMA26, indicating a potential uptrend."
	} else if latestEMA12 < latestEMA26 && prevEMA12 >= prevEMA26 {
		return "Potential downtrend", "Bearish crossover detected. EMA12 has crossed below EMA26, indicating a potential downtrend."
	} else if latestEMA12 > latestEMA26 {
		return "Potential uptrend", "EMA12 is above EMA26, indicating a potential uptrend."
	} else if latestEMA12 < latestEMA26 {
		return "Potential downtrend", "EMA12 is below EMA26, indicating a potential downtrend."
	}
	return "Neutral", "No significant crossover detected."
}

// min returns the smaller of x or y.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// calculateRSI calculates the Relative Strength Index (RSI)
func CalculateRSI(prices []float64, period int) ([]float64, error) {
	if len(prices) < period {
		return nil, fmt.Errorf("not enough data to calculate RSI for the given period")
	}

	// Initialize gains and losses
	gains := make([]float64, len(prices))
	losses := make([]float64, len(prices))

	// Calculate daily gains and losses
	for i := 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains[i] = change
		} else {
			losses[i] = -change
		}
	}

	// Calculate average gains and losses for the first period
	avgGain := 0.0
	avgLoss := 0.0
	for i := 1; i <= period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Initialize RSI array
	rsi := make([]float64, len(prices))

	// Calculate RSI for the first period
	rs := avgGain / avgLoss
	rsi[period-1] = 100 - (100 / (1 + rs))

	// Calculate RSI for subsequent periods
	for i := period; i < len(prices); i++ {
		avgGain = (avgGain*float64(period-1) + gains[i]) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + losses[i]) / float64(period)
		rs = avgGain / avgLoss
		rsi[i] = 100 - (100 / (1 + rs))
	}

	return rsi[period-1:], nil // Return only the valid RSI values
}

// analyzeRSI analyzes the RSI value and returns a descriptive analysis
func AnalyzeRSI(rsiValue float64) string {
	if rsiValue > 70 {
		return fmt.Sprintf("RSI is %.2f, indicating the asset is overbought.", rsiValue)
	} else if rsiValue < 30 {
		return fmt.Sprintf("RSI is %.2f, indicating the asset is oversold.", rsiValue)
	} else {
		return fmt.Sprintf("RSI is %.2f, indicating normal market conditions.", rsiValue)
	}
}

// fetchHighLowPrices fetches the high and low prices for a given symbol
func FetchHighLowPrices(symbol string) ([]float64, []float64, error) {
	params := &chart.Params{
		Symbol:   symbol,
		Interval: datetime.OneDay, // Daily data
		Start: &datetime.Datetime{
			Month: int(time.Now().AddDate(-1, 0, 0).Month()),
			Day:   1,
			Year:  time.Now().AddDate(-1, 0, 0).Year(),
		},
		End: &datetime.Datetime{
			Month: int(time.Now().Month()),
			Day:   time.Now().Day(),
			Year:  time.Now().Year(),
		},
	}

	iter := chart.Get(params)
	var highs []float64
	var lows []float64
	for iter.Next() {
		p := iter.Bar()
		high, _ := p.High.Float64()
		low, _ := p.Low.Float64()
		highs = append(highs, high)
		lows = append(lows, low)
	}

	if iter.Err() != nil {
		return nil, nil, iter.Err()
	}

	return highs, lows, nil
}

// calculateStochasticOscillator calculates the Stochastic Oscillator
func CalculateStochasticOscillator(closes, highs, lows []float64, period int) ([]float64, []float64, error) {
	if len(closes) < period || len(highs) < period || len(lows) < period {
		return nil, nil, fmt.Errorf("not enough data to calculate Stochastic Oscillator for the given period")
	}

	k := make([]float64, len(closes)-period+1)
	for i := period - 1; i < len(closes); i++ {
		low := lows[i-period+1]
		high := highs[i-period+1]
		for j := i - period + 1; j <= i; j++ {
			if lows[j] < low {
				low = lows[j]
			}
			if highs[j] > high {
				high = highs[j]
			}
		}
		k[i-period+1] = 100 * ((closes[i] - low) / (high - low))
	}

	// Calculate %D as a 3-period SMA of %K
	d := make([]float64, len(k))
	for i := 2; i < len(k); i++ {
		d[i] = (k[i] + k[i-1] + k[i-2]) / 3
	}

	return k, d, nil
}

// analyzeStochasticOscillator analyzes the Stochastic Oscillator values
func AnalyzeStochasticOscillator(k, d []float64) string {
	if len(k) == 0 || len(d) == 0 {
		return "Not enough data for Stochastic Oscillator analysis."
	}

	latestK := k[len(k)-1]
	latestD := d[len(d)-1]

	if latestK > 80 && latestD > 80 {
		return fmt.Sprintf("Stochastic Oscillator is %.2f/%.2f, indicating the asset is overbought.", latestK, latestD)
	} else if latestK < 20 && latestD < 20 {
		return fmt.Sprintf("Stochastic Oscillator is %.2f/%.2f, indicating the asset is oversold.", latestK, latestD)
	} else {
		return fmt.Sprintf("Stochastic Oscillator is %.2f/%.2f, indicating normal market conditions.", latestK, latestD)
	}
}
