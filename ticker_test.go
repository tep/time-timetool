package timetool

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func TestNormalTicker(t *testing.T) {
	rand.Seed(1)
	nt := NewNormalTicker(context.Background(), 75*time.Millisecond, 15.*time.Millisecond)

	var pt time.Time

	want := []int64{0, 73, 68, 110, 80, 84, 78, 90, 64, 86, 99, 88, 95, 83, 86, 59, 86, 82, 90, 52, 71, 104, 92, 60, 90}

	for i := 0; i < 25; i++ {
		tv := <-nt.C

		if !pt.IsZero() {
			if got := tv.Sub(pt).Round(time.Millisecond).Milliseconds(); !plusOrMinusOne(got, want[i]) {
				t.Errorf("iteration %d: got %d; wanted %d (±1)", i, got, want[i])
			}
		}

		pt = tv
	}
}

func plusOrMinusOne(a, b int64) bool {
	return a >= b-1 && a <= b+1
}
