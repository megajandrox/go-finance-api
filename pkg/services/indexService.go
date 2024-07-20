package services

import (
	"fmt"
	"log"
	"time"

	"github.com/megajandrox/go-finance-api/pkg/models"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

type Indexes struct {
	symbol     string
	SMA        models.SMA
	EMA        models.EMA
	MACD       models.MACD
	RSI        models.RSI
	Stochastic models.Stochastic
	Volume     models.Volume
	OBV        models.OBV
	RVOL       models.RVOL
}

// getQuote handles the retrieval of stock quotes
func FindIndexesBySymbol(symbol string) (Indexes, error) {
	startDateTime := &datetime.Datetime{
		Month: int(time.Now().AddDate(-1, 0, 0).Month()),
		Day:   1,
		Year:  time.Now().AddDate(-1, 0, 0).Year(),
	}
	endDateTime := &datetime.Datetime{
		Month: int(time.Now().Month()),
		Day:   time.Now().Day(),
		Year:  time.Now().Year(),
	}
	params := &chart.Params{
		Symbol:   symbol,
		Interval: datetime.OneHour,
		Start:    startDateTime,
		End:      endDateTime,
	}
	iter := chart.Get(params)
	indexesResult := Indexes{symbol: symbol}
	var closes []float64
	var highs []float64
	var lows []float64
	var volumes []float64
	for iter.Next() {
		p := iter.Bar()
		close, _ := p.Close.Float64()
		high, _ := p.High.Float64()
		low, _ := p.Low.Float64()
		closes = append(closes, close)
		highs = append(highs, high)
		lows = append(lows, low)
		volumes = append(volumes, float64(p.Volume))
	}
	if err := iter.Err(); err != nil {
		fmt.Println(err)
	}
	// Analyze SMA
	sma, errSMA := models.NewSMA(symbol)
	if errSMA != nil {
		log.Fatalf("Error creating SMA: %v", errSMA)
	}

	sma.AnalyzeSMATrend(closes)
	indexesResult.SMA = *sma
	// Analyze EMA
	ema, errEMA := models.NewEMA(symbol)
	if errEMA != nil {
		log.Fatalf("Error creating EMA: %v", errEMA)
	}
	ema.AnalyzeEMACrossover(closes)
	indexesResult.EMA = *ema
	// Analyze MACD
	macd, errMACD := models.NewMACD(symbol)
	if errMACD != nil {
		log.Fatalf("Error creating MACD: %v", errMACD)
	}
	macd.AnalyzeMACD(closes)
	indexesResult.MACD = *macd
	// Analyze the RSI value
	rsi, errRSI := models.NewRSI(symbol)
	if errRSI != nil {
		log.Fatalf("Error creating RSI: %v", errRSI)
	}
	rsi.AnalyzeRSI(closes)
	indexesResult.RSI = *rsi

	// Analyze the Stochastic Oscillator
	sto, errSTO := models.NewStochastic(symbol)
	if errSTO != nil {
		log.Fatalf("Error creating Stochastic: %v", errSTO)
	}
	sto.AnalyzeStochasticOscillator(closes, lows, highs)
	indexesResult.Stochastic = *sto

	// Analyze the volume trend
	vol, errVol := models.NewVolume(symbol)
	if errVol != nil {
		log.Fatalf("Error creating Volume: %v", errVol)
	}
	vol.AnalyzeVolumeTrend(volumes)
	indexesResult.Volume = *vol

	// Analyze the OBV
	obv, errOBV := models.NewOBV(symbol)
	if errOBV != nil {
		log.Fatalf("Error creating OBV: %v", errOBV)
	}
	obv.AnalyzeOBV(closes, volumes)
	indexesResult.OBV = *obv

	// Analyze the RVOL
	rvol, errVOL := models.NewRVOL(symbol)
	if errVOL != nil {
		log.Fatalf("Error creating RVOL: %v", errVOL)
	}
	rvol.AnalyzeRVOL(volumes)
	indexesResult.RVOL = *rvol

	return indexesResult, nil
}
