package timetools

import "time"

func init() {
	resetTimeFuncs()
}

var (
	timeAfter func(time.Duration) <-chan time.Time
	timeNow   func() time.Time
	timeSleep func(time.Duration)
)

func resetTimeFuncs() {
	timeAfter = time.After
	timeNow = time.Now
	timeSleep = time.Sleep
}
