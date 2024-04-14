// Copyright © 2018 Timothy E. Peoples

// XXX This file contains deprecated cod that will be removed in the future. XXX

package timetool

import (
	"context"
	"math"
	"time"
)

const (
	backoffPower = 3
)

// RetryWithBackoff calls the given RetryFunc a maximum of iters times until it
// returns true. The provided context must have a defined deadline and the
// number of iterations requested must be at least 2. Nil is returned if retry
// returns true before the deadline expires and within the stated number of
// iterations.
//
// Note that, even if retry returns true, the deadline is checked a final time
// and, if it has expired, ctx.Err() is returned.
//
// If the provided Context has no defined deadline, ErrMissingDeadline is
// returned. ErrTooFewIterations will be returned if iters is less than 2.
// If each call to retry returns false, ErrRetriesExhausted is returned.
//
// Deprecated: Please use *Backoff.Retry instead.
func RetryWithBackoff(ctx context.Context, iters int, retry RetryFunc) error {
	if iters < 2 {
		return ErrTooFewIterations
	}

	bos, err := newBackoffSession(ctx, iters)
	if err != nil {
		return err
	}

	for i := 0; i < iters; i++ {
		if err := bos.sleep(ctx, i); err != nil {
			return err
		}

		if retry(i) {
			return contextDoneOr(ctx, nil)
		}
	}

	return contextDoneOr(ctx, ErrRetriesExhausted)
}

// RetryWithBackoffDuration is a wrapper around RetryWithBackoff accepting
// a Duration in leu of a Context.
func RetryWithBackoffDuration(dur time.Duration, iters int, retry RetryFunc) error {
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()
	return RetryWithBackoff(ctx, iters, retry)
}

type backoffSession struct {
	start       time.Time
	timeout     time.Duration
	denominator float64
}

func newBackoffSession(ctx context.Context, iters int) (*backoffSession, error) {
	now := time.Now()

	dl, ok := ctx.Deadline()
	if !ok {
		return nil, ErrMissingDeadline
	}

	if ttd := dl.Sub(time.Now()); ttd > 0 {
		return &backoffSession{now, ttd, math.Pow(float64(iters), backoffPower)}, nil
	}

	return nil, contextDoneOr(ctx, ErrTimeWarp)
}

func (bo *backoffSession) sleep(ctx context.Context, i int) error {
	if i == 0 {
		return contextDoneOr(ctx, nil)
	}

	// START + (timeout * fraction) - NOW   =>  SLEEP
	fx := math.Pow(float64(i), backoffPower) / bo.denominator
	dt := time.Duration(float64(bo.timeout) * fx)
	sd := bo.start.Add(dt).Sub(time.Now())

	ch := time.After(sd)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
	}

	return nil
}
