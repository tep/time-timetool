package timetool

import (
	"testing"
	"time"
)

var (
	millisValue = int64(1279156356512)
	timeValue   = time.Date(2010, 7, 14, 18, 12, 36, int(512*time.Millisecond), time.Local)
)

func TestToMillis(t *testing.T) {
	want := millisValue
	if got := ToMillis(timeValue); got != want {
		t.Errorf("ToMillis(%v) == %d; Wanted: %d", timeValue, got, want)
	}
}

func TestFromMillis(t *testing.T) {
	want := timeValue
	if got := FromMillis(millisValue); !got.Equal(want) {
		t.Errorf("FromMillis(%d) == %v; Wanted: %v", millisValue, got, want)
	}
}
