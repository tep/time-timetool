package timetool

import (
	"context"
	"math/rand"
	"time"
)

// Type NormalTicker holds a channel that delivers "ticks" of a clock over
// a normally distributed time interval.
type NormalTicker struct {
	C      chan time.Time
	done   chan struct{}
	mean   time.Duration
	stddev time.Duration
	err    error
}

// NewNormalTicker returns a new NormalTicker containing a channel that will
// send the current time on the channel after each tick. The period of the
// ticks is over a normal distribution as specified by the mean and stddev
// arguments. The ticker will drop ticks to make up for slow receivers and
// will continue to send values to its channel until the Stop method is called
// or the given context is expired.
func NewNormalTicker(ctx context.Context, mean, stddev time.Duration) *NormalTicker {
	nt := &NormalTicker{
		C:      make(chan time.Time),
		done:   make(chan struct{}),
		mean:   mean,
		stddev: stddev,
	}

	go nt.run(ctx)

	return nt
}

// Stop turns off the ticker. After Stop, no more ticks will be sent. Stop does
// not close the channel, to prevent a concurrent goroutine reading from the
// channel from seeing an erroneous "tick". If Stop is called before the
// constructor's Context has expired, the Err method will return a nil error.
func (nt *NormalTicker) Stop() {
	close(nt.done)
}

// Err returns an error indicating how the ticker was stopped. If the Stop
// method was called, a nil error returned. If the constructor's Context has
// expired, ctx.Err() is returned. If the ticker has not been stopped,
// ErrTickerActive is returned.
func (nt *NormalTicker) Err() error {
	return nt.err
}

func (nt *NormalTicker) run(ctx context.Context) {
	t := time.NewTimer(nt.duration())

	defer stopAndFlush(t)

	for {
		if done, err := nt.onePass(ctx, t); done || err != nil {
			nt.err = err
			return
		}

		t.Reset(nt.duration())
	}
}

func (nt *NormalTicker) onePass(ctx context.Context, tt *time.Timer) (bool, error) {
	var tv time.Time

	select {
	case <-ctx.Done():
		return true, ctx.Err()

	case <-nt.done:
		return true, nil

	case tv = <-tt.C:
	}

	select {
	case <-ctx.Done():
		return true, ctx.Err()

	case <-nt.done:
		return true, nil

	case nt.C <- tv:
		return false, nil
	}
}

func (nt *NormalTicker) duration() time.Duration {
	return time.Duration(rand.NormFloat64()*float64(nt.stddev) + float64(nt.mean))
}

func stopAndFlush(t *time.Timer) {
	if t == nil || t.Stop() {
		return
	}

	select {
	case <-t.C:
	default:
	}
}
