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
	"math/rand"
	"time"
)

const (
	minIterations = 2
)

// RetryFunc is the function provided to a retry operation that should be
// executed until it succeeds, as indicated by its return value. i.e. If the
// function returns false, it will be retried after a brief delay.
//
// The function will be passed a single int argument indicating the current
// (zero based) iteration number.
type RetryFunc func(i int) bool

// Backoff defines the parameters for a set of retries with exponential
// backoff.
type Backoff struct {
	// Iterations declares the maximum number execution attempts.
	Iterations int

	// Coefficient indicates the initial delay between attempts.
	Coefficient time.Duration

	// Jitter is a random modifier percentage applied to each delay period.
	Jitter float64

	startWait time.Duration
	initWait  time.Duration
}

// StdBackoff provides a Backoff with common parameters.
var StdBackoff = &Backoff{Iterations: 5, Coefficient: time.Second, Jitter: 0.1}

// CalculateBackoff returns a new Backoff having a Coefficient value
// calculated for the given number of iterations and desired total
// delay time.
//
// If a non-zero jitter value is provided, it will be included in the return
// value but it will not be used to calculate the resulting Coefficient.
//
// An error is returned if iters < 2, total < 0, or jitter is outside [0, 100).
func CalculateBackoff(iters int, total time.Duration, jitter float64) (*Backoff, error) {
	b := &Backoff{
		Iterations: iters,
		Jitter:     jitter,
	}

	return b.WithTotalDelay(total)
}

// MustCalculateBackoff is a wrapper around CalculateBackoff that will panic if
// an error is returned.
func MustCalculateBackoff(iters int, total time.Duration, jitter float64) *Backoff {
	b, err := CalculateBackoff(iters, total, jitter)
	if err != nil {
		panic(err)
	}
	return b
}

// WithTotalDelay returns a pointer to its receiver with a new Coefficient
// value calculated to consume the provided Duration based on the receiver's
// other field values.
//
// This is similar to the CalculateBackoff function but also considers the
// optional "startup wait time" and/or "initial wait time" (if any) when
// calculating a new Coefficient value.
//
// A Coefficient value will only be calculated if the receiver's Iterations
// field is greater than 2 -- since a call to Retry with only 2 iterations
// never refers to the Coefficient field. However, if Iterations is exactly
// 2, this method has the same effect as calling WithInitialWait.
//
// On error, a nil pointer will always be returned.
//
// ErrNegativeDelay is retured if the resultant total delay time, after any
// "startup wait time" and/or "initial wait time" values are considered, is
// calculated to be a negative value.
//
// If the receiver's Iterations field is less than 2, ErrTooFewIterations is
// returned.
func (b Backoff) WithTotalDelay(d time.Duration) (*Backoff, error) {
	if b.Iterations < minIterations {
		return nil, ErrTooFewIterations
	}

	if b.Iterations == minIterations {
		b.Coefficient = 1 // To prevent future validation failure
		return b.WithInitialWait(d), nil
	}

	if d -= (b.startWait + b.initWait); d < 0 {
		return nil, ErrNegativeDelay
	}

	b.Coefficient = time.Duration(float64(d) / float64(int((1<<(b.Iterations-2))-1)))

	if err := b.validate(); err != nil {
		return nil, err
	}

	return &b, nil
}

// MustTotalDelay is a wrapper around WithTotalDelay that will panic if an
// error is returned.
func (b Backoff) MustTotalDelay(d time.Duration) *Backoff {
	out, err := b.WithTotalDelay(d)
	if err != nil {
		panic(err)
	}
	return out
}

// WithStartWait returns a pointer to its receiver that adds a "startup wait
// time" to Retry execution (for those situations where you're quite certain
// the RetryFunc will not succeed right away).
//
// Normally, the Retry method will execute its first attempt immediately.
// However, with a defined "startup wait time", Retry will Sleep for the
// given duration before executing its first attempt.
//
// If the Context provided to Retry becomes done during this startup wait
// time, Retry will immediately return ctx.Err() without executing any
// attempts.
func (b Backoff) WithStartWait(d time.Duration) *Backoff {
	b.startWait = d
	return &b
}

// WithInitialWait returns a pointer to its receiver and adds an "initial
// wait time" to Retry execution, after a failed first attempt.
//
// Normally, when a RetryFunc does not succeed on its first attempt,
// the next attempt is executed immediately without delay. Adding this
// "initial wait time" will inject an additional (static) delay in this
// situation.
//
// No other inter-attempt delays are effected by this value; those delays
// will continue to be calculated using the current iteration number and
// the reciever's Coefficient and Jitter fields.
//
// If the Context provided to Retry becomes done during this initial wait
// time, Retry will immediately return ctx.Err() without executing its next
// attempt.
func (b Backoff) WithInitialWait(d time.Duration) *Backoff {
	b.initWait = d
	return &b
}

// Retry calls the given RetryFunc up to b.Iterations times until it returns
// true or the provided Context is cancelled, whichever comes first.
//
// If the initial call to RetryFunc returns false (indicating it should be
// reattempted), the RetryFunc is (by default) rerun immediately. Subsequent
// execution attempts (after the first) are interleaved with an exponentially
// increasing delay based on the receiver such that each delay is calculated
// as:
//
//	multiple  =  2 ** (attempt - 1)
//	delay     =  b.Coefficient * multiple
//
// ...or, if b.Jitter is non zero:
//
//	multiple  =  2^(attempt - 1)
//	jitter    =  (b.Jitter * random) - (b.Jitter / 2) / 100
//	multiple +=  multiple * jitter
//	delay     =  b.Coefficient * multiple
//
// If the receiver declares fewer than 2 iterations an error will be returned.
//
// The receiver's delay Coefficient must be a positive, non-zero value or
// an error will be returned.
//
// If b.Jitter is 0, no Jitter will be applied. Otherwise, the Jitter value
// must be in the range (0,100) or an error will be returned.
func (b *Backoff) Retry(ctx context.Context, retry RetryFunc) error {
	if err := b.validate(); err != nil {
		return err
	}

	// We'll give attempt #0 special handling with an optional Startup Delay...
	if err := Sleep(ctx, b.startWait); err != nil {
		return err
	}

	// ...before running the 'retry' func for the first time...
	if retry(0) {
		return contextDoneOr(ctx, nil)
	}

	// ...before entering our retry loop on attempt #1.
	for attempt := 1; attempt < b.Iterations; attempt++ {
		if attempt == 1 {
			if err := Sleep(ctx, b.initWait); err != nil {
				return err
			}
		}

		if retry(attempt) {
			return contextDoneOr(ctx, nil)
		}

		multiple := float64(uint(1) << (uint(attempt) - 1))

		if b.Jitter != 0 {
			j := (((b.Jitter * rand.Float64()) - (b.Jitter / 2)) / 100)
			multiple += multiple * j
		}

		if err := Sleep(ctx, b.Coefficient*time.Duration(multiple)); err != nil {
			return err
		}
	}

	return contextDoneOr(ctx, ErrRetriesExhausted)
}

func (b *Backoff) validate() error {
	switch {
	case b == nil:
		return ErrNilReceiver

	case b.Iterations < minIterations:
		return ErrTooFewIterations

	case b.Coefficient == 0:
		return ErrZeroCoefficient

	case b.Coefficient < 0:
		return ErrNegativeDelay

	case b.Jitter != 0 && (b.Jitter < 0 || b.Jitter >= 100):
		return ErrBadJitter

	default:
		return nil
	}
}

func contextDoneOr(ctx context.Context, err error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return err
	}
}
