package handlers

type Interval string

const (
	// Define intervals
	SixtyMins Interval = "60m"
	OneHour   Interval = "1h"
	OneDay    Interval = "1d"
)

var ValidIntervals = []Interval{
	SixtyMins, OneHour, OneDay,
}

// IsValidInterval checks if the interval is valid.
func IsValidInterval(interval Interval) bool {
	for _, v := range ValidIntervals {
		if v == interval {
			return true
		}
	}
	return false
}
