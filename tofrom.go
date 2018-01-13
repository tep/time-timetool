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

import "time"

// FromMillis interprets millis as milliseconds since the Epoch and returns
// the equivalent time.Time value.
func FromMillis(millis int64) time.Time {
	sse := millis / 1000
	mso := int64(time.Millisecond) * (millis - (sse * 1000))

	return time.Unix(sse, mso)
}

// ToMillis returns the number of milliseconds since the Epoch for the provided
// time.Time value t.
func ToMillis(t time.Time) int64 {
	t = t.Round(time.Millisecond)
	return t.UnixNano() / int64(time.Millisecond)
}
