package api

import "time"

// nowMs returns the current time in milliseconds using the monotonic clock.
func nowMs() float64 {
	return float64(time.Now().UnixNano()) / 1e6
}
