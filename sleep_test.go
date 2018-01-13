// Copyright © 2018 Tim Peoples <coders@toolman.org>
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
