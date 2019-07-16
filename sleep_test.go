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
	"fmt"
	"testing"
	"time"
)

var now = time.Date(2010, 7, 14, 18, 0, 0, 0, time.UTC)

func TestSleep(t *testing.T) {
	defer resetTimeFuncs()
	sd := 1 * time.Second
	ch := make(chan time.Time)

	awake := func() {
		go func() {
			ch <- now
			close(ch)
		}()
	}

	timeAfter = func(d time.Duration) <-chan time.Time {
		if d != sd {
			t.Errorf("bad arg1 to time.After(); Got %v; Wanted %v", d, sd)
		}
		return ch
	}

	t.Run("Cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := doTestSleep(ctx, cancel, sd, context.Canceled); err != nil {
			t.Error(err)
		}
	})

	t.Run("Completed", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := doTestSleep(ctx, awake, sd, nil); err != nil {
			t.Error(err)
		}
	})
}

func doTestSleep(ctx context.Context, prepFunc func(), d time.Duration, want error) error {
	prepFunc()

	if got := Sleep(ctx, d); got != want {
		return fmt.Errorf("Sleep(ctx, %#v) == %v; Wanted %v", d, got, want)
	}

	return nil
}
