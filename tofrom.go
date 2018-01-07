package timetool

import "time"

func FromMillis(millis int64) time.Time {
	sse := millis / 1000
	mso := int64(time.Millisecond) * (millis - (sse * 1000))

	return time.Unix(sse, mso)
}

func ToMillis(t time.Time) int64 {
	t = t.Round(time.Millisecond)
	return t.UnixNano() / int64(time.Millisecond)
}
