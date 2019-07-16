// Copyright 2018 Timothy E. Peoples
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package timetool

import (
	"context"
	"errors"
	"math"
	"time"
)

const (
	backoffPower = 3
)

// ErrMissingDeadline is returned by functions expecting a Deadline
// but are passed a Context devoid of such nature.
var ErrMissingDeadline = errors.New("context must have a deadline")

// ErrTooFewIterations is returned by RetryWithBackoff if the number of
// requested iterations is less than 2.
var ErrTooFewIterations = errors.New("number of iterations must be at least 2")

// ErrTimeWarp is returned by RetryWithBackoff in the freakishly uncommon
// event that a context ends up with a deadline in the past and a Done
// channel that does not return.
var ErrTimeWarp = errors.New("valid context has deadline in the past")

// ErrRetriesExhausted is returned by RetryWithBackoff if all retry
// attempts have been unsuccessful.
var ErrRetriesExhausted = errors.New("all retries exhausted")

// RetryFunc is a function passed to RetryWithBackoff that should be retried
// until successful. It should return true if the operation was successful or
// false if it should be retried after a brief delay. The function will be
// passed a single int argument which is the current (zero based) iteration
// number.
type RetryFunc func(i int) bool

// RetryWithBackoff calls the RetryFunc retry a maximum of iters times until it
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

func contextDoneOr(ctx context.Context, err error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return err
	}
	return nil
}
