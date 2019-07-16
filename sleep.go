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
