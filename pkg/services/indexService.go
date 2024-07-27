package services

import (
	"fmt"
	"log"
	"time"

	"github.com/megajandrox/go-finance-api/pkg/models"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

func ConvertTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04")
}

// getQuote handles the retrieval of stock quotes
func FindIndexesBySymbol(symbol string, from int, interval string) (models.Indexes, error) {
	startDateTime := &datetime.Datetime{
		Month: int(time.Now().AddDate(0, from*-1, 0).Month()),
		Day:   1,
		Year:  time.Now().AddDate(0, from*-1, 0).Year(),
	}
	endDateTime := &datetime.Datetime{
		Month: int(time.Now().Month()),
		Day:   time.Now().Day(),
		Year:  time.Now().Year(),
	}
	params := &chart.Params{
		Symbol:   symbol,
		Interval: datetime.Interval(interval),
		Start:    startDateTime,
		End:      endDateTime,
	}
	iter := chart.Get(params)
	indexesResult := models.NewIndexes(symbol)
	var marketDataList []models.BasicMarketData
	for iter.Next() {
		p := iter.Bar()
		close, _ := p.Close.Float64()
		high, _ := p.High.Float64()
		low, _ := p.Low.Float64()
		var marketData = models.BasicMarketData{Close: close, High: high, Low: low, Volume: int64(p.Volume), TimeStamp: int64(p.Timestamp)}
		fmt.Println(marketData)
		fmt.Printf(" FECHA: %s", ConvertTimestamp(marketData.TimeStamp))
		marketDataList = append(marketDataList, marketData)
	}
	if err := iter.Err(); err != nil {
		fmt.Println(err)
	}

	analyzers := []func(string) (models.Analyzer, error){
		models.NewSMAAdapter,
		models.NewEMAAdapter,
		models.NewMACDAdapter,
		models.NewRSIAdapter,
		models.NewStochasticAdapter,
		models.NewVolumeAdapter,
		models.NewOBVAdapter,
		models.NewRVOLAdapter,
		models.NewATRAdapter,
		models.NewADXAdapter,
		models.NewMomentumAdapter,
		models.NewCCIAdapter,
	}

	for _, newAnalyzer := range analyzers {
		var err error
		indexesResult, err = models.RunAnalysis(symbol, marketDataList, indexesResult, newAnalyzer)
		if err != nil {
			log.Println("Error running analysis: %v", err)
		}
	}
	return *indexesResult, nil
}
