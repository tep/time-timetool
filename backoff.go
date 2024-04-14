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
}

// StdBackoff provides a Backoff with common parameters.
var StdBackoff = &Backoff{5, time.Second, 0.1}

// CalculateBackoff calculates a Backoff value for the given number of
// iterations and a desired total delay time. An error is returned if iters
// < 3, total < 0, or jitter is outside [0, 100).
func CalculateBackoff(iters int, total time.Duration, jitter float64) (*Backoff, error) {
	switch {
	case iters < 3:
		return nil, ErrTooFewIterations
	case total < 0:
		return nil, ErrNegativeDelay
	case jitter < 0 || jitter >= 100:
		return nil, ErrBadJitter
	}

	coef := time.Duration(float64(total) / float64(int((1<<(iters-2))-1)))

	return &Backoff{Iterations: iters, Coefficient: coef, Jitter: jitter}, nil
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

// Retry calls the given RetryFunc up to b.Iterations times until it returns
// true or the provided Context is cancelled, whichever comes first. If the
// initial call to RetryFunc returns false (meaning it failed and a retry is
// needed), the RetryFunc is rerun immediately. Subsequent execution attempts
// (after the first) are interleaved with an exponentially increasing delay
// based on the receiver such that each delay is calculated as:
//
//	multiple  =  2^(attempt - 1)
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
//
func (b *Backoff) Retry(ctx context.Context, retry RetryFunc) error {
	if err := b.validate(); err != nil {
		return err
	}

	for attempt := 0; attempt < b.Iterations; attempt++ {
		if retry(attempt) {
			return contextDoneOr(ctx, nil)
		}

		if attempt == 0 {
			continue
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
