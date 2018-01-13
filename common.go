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

// Package timetool provides tools and utilities for dealing with our most
// preceious comodity: time.
package timetool // import "toolman.org/time/timetool"

import "time"

func init() {
	resetTimeFuncs()
}

var (
	timeAfter func(time.Duration) <-chan time.Time
	timeNow   func() time.Time
	timeSleep func(time.Duration)
)

func resetTimeFuncs() {
	timeAfter = time.After
	timeNow = time.Now
	timeSleep = time.Sleep
}
