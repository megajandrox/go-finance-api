package models

import "time"

// extractDailyCloses extrae el precio de cierre diario de los datos intradía
func ExtractDailyCloses(data []BasicMarketData) []BasicMarketData {
	dailyCloses := []BasicMarketData{}
	var currentDay time.Time
	var lastClose BasicMarketData

	for _, entry := range data {
		entryTime := time.Unix(entry.TimeStamp, 0)
		if !currentDay.IsZero() && entryTime.Day() != currentDay.Day() {
			dailyCloses = append(dailyCloses, lastClose)
		}
		currentDay = entryTime
		lastClose = entry
	}
	// Agregar el último precio de cierre del último día
	if !currentDay.IsZero() {
		dailyCloses = append(dailyCloses, lastClose)
	}

	return dailyCloses
}
