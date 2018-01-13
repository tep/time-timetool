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
