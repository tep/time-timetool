package timetool

import "time"

// FromMillis interprets millis as milliseconds since the Epoch and returns
// the equivalent time.Time value.
func FromMillis(millis int64) time.Time {
	sse := millis / 1000
	mso := int64(time.Millisecond) * (millis - (sse * 1000))

	return time.Unix(sse, mso)
}

// ToMillis returns the number of milliseconds since the Epoch for the provided
// time.Time value t.
func ToMillis(t time.Time) int64 {
	t = t.Round(time.Millisecond)
	return t.UnixNano() / int64(time.Millisecond)
}
