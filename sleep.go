package timetool

import (
	"context"
	"time"
)

// Sleep is a wrapper around time.Sleep that may be interrupted by the
// cancellation of a Context. Sleep returns ctx.Err() if cancelled by
// the Context, otherwise it returns nil.
func Sleep(ctx context.Context, d time.Duration) error {
	var err error
	ch := timeAfter(d)
	select {
	case <-ctx.Done():
		err = ctx.Err()

	case <-ch:
		err = nil
	}

	return err
}

// SleepUntil is a wrapper around Sleep that accepts a time.Time instead
// of a time.Duration.
func SleepUntil(ctx context.Context, t time.Time) error {
	return Sleep(ctx, t.Sub(timeNow()))
}
