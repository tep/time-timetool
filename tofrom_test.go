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
