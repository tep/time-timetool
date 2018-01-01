// Copyright © 2017 Tim Peoples <coders@toolman.org>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

// Package timetools provides tools and utilities for dealing with our most
// preceious comodity: time.
package timetools

import (
	"context"
	"errors"
	"math"
	"time"
)

const (
	backoffPower = 3
)

var (
	// MissingDeadlineError is returned by functions expecting a Deadline
	// but are passed a Context devoid of such nature.
	MissingDeadlineError = errors.New("context must have a deadline")

	// TooFewIterationsError is returned by RetryWithBackoff if the number of
	// requested iterations is less than 2.
	TooFewIterationsError = errors.New("number of iterations must be at least 2")

	// TimeWarpError is returned by RetryWithBackoff in the freakishly uncommon
	// event that a context ends up with a deadline in the past and a Done
	// channel that does not return.
	TimeWarpError = errors.New("valid context has deadline in the past!")

	// RetriesExhaustedError is returned by RetryWithBackoff if all retry
	// attempts have been unsuccessful.
	RetriesExhaustedError = errors.New("all retries exhausted")
)

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
// If the provided Context has no defined deadline, MissingDeadlineError is
// returned. TooFewIterationsError will be returned if iters is less than 2.
// If each call to retry returns false, RetriesExhaustedError is returned.
func RetryWithBackoff(ctx context.Context, iters int, retry RetryFunc) error {
	if iters < 2 {
		return TooFewIterationsError
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

	return contextDoneOr(ctx, RetriesExhaustedError)
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
		return nil, MissingDeadlineError
	}

	if ttd := dl.Sub(time.Now()); ttd > 0 {
		return &backoffSession{now, ttd, math.Pow(float64(iters), backoffPower)}, nil
	}

	return nil, contextDoneOr(ctx, TimeWarpError)
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
