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
	"testing"
	"time"
)

var (
	millisValue = int64(1279156356512)
	timeValue   = time.Date(2010, 7, 14, 18, 12, 36, int(512*time.Millisecond), time.Local)
)

func TestToMillis(t *testing.T) {
	want := millisValue
	if got := ToMillis(timeValue); got != want {
		t.Errorf("ToMillis(%v) == %d; Wanted: %d", timeValue, got, want)
	}
}

func TestFromMillis(t *testing.T) {
	want := timeValue
	if got := FromMillis(millisValue); !got.Equal(want) {
		t.Errorf("FromMillis(%d) == %v; Wanted: %v", millisValue, got, want)
	}
}
