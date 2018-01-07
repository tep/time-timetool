
# timetool
`import "toolman.org/time/timetool"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package timetool provides tools and utilities for dealing with our most
preceious comodity: time.

## Install

``` sh
  go get toolman.org/time/timetool
```

## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func FromMillis(millis int64) time.Time](#FromMillis)
* [func RetryWithBackoff(ctx context.Context, iters int, retry RetryFunc) error](#RetryWithBackoff)
* [func Sleep(ctx context.Context, d time.Duration) error](#Sleep)
* [func SleepUntil(ctx context.Context, t time.Time) error](#SleepUntil)
* [func ToMillis(t time.Time) int64](#ToMillis)
* [type RetryFunc](#RetryFunc)


#### <a name="pkg-files">Package files</a>
[backoff.go](/src/toolman.org/time/timetool/backoff.go) [common.go](/src/toolman.org/time/timetool/common.go) [sleep.go](/src/toolman.org/time/timetool/sleep.go) [tofrom.go](/src/toolman.org/time/timetool/tofrom.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (
    // ErrMissingDeadline is returned by functions expecting a Deadline
    // but are passed a Context devoid of such nature.
    ErrMissingDeadline = errors.New("context must have a deadline")

    // ErrTooFewIterations is returned by RetryWithBackoff if the number of
    // requested iterations is less than 2.
    ErrTooFewIterations = errors.New("number of iterations must be at least 2")

    // ErrTimeWarp is returned by RetryWithBackoff in the freakishly uncommon
    // event that a context ends up with a deadline in the past and a Done
    // channel that does not return.
    ErrTimeWarp = errors.New("valid context has deadline in the past")

    // ErrRetriesExhausted is returned by RetryWithBackoff if all retry
    // attempts have been unsuccessful.
    ErrRetriesExhausted = errors.New("all retries exhausted")
)
```


## <a name="FromMillis">func</a> [FromMillis](/tofrom.go?s=144:183#L1)
``` go
func FromMillis(millis int64) time.Time
```
FromMillis interprets `millis` as milliseconds since the Epoch and returns
the equivalent `time.Time` value.



## <a name="RetryWithBackoff">func</a> [RetryWithBackoff](/backoff.go?s=2745:2821#L59)
``` go
func RetryWithBackoff(ctx context.Context, iters int, retry RetryFunc) error
```
RetryWithBackoff calls the RetryFunc `retry` a maximum of `iters` times until
it returns `true`. The provided `Context` must have a defined deadline and the
number of iterations requested must be at least 2. Nil is returned if `retry`
returns `true` before the deadline expires and within the stated number of
iterations.

Note that, even if `retry` returns `true`, the deadline is checked a final
time and, if it has expired, `ctx.Err()` is returned.

If the provided `Context` has no defined deadline, `ErrMissingDeadline` is
returned. `ErrTooFewIterations` will be returned if `iters` is less than 2.
If each call to `retry` returns `false`, `ErrRetriesExhausted` is returned.



## <a name="Sleep">func</a> [Sleep](/sleep.go?s=232:286#L1)
``` go
func Sleep(ctx context.Context, d time.Duration) error
```
Sleep is a wrapper around `time.Sleep` that may be interrupted by the
cancellation of a `Context`. Sleep returns `ctx.Err()` if cancelled by
the Context, otherwise it returns `nil`.



## <a name="SleepUntil">func</a> [SleepUntil](/sleep.go?s=512:567#L17)
``` go
func SleepUntil(ctx context.Context, t time.Time) error
```
SleepUntil is a wrapper around Sleep that accepts a time.Time instead
of a time.Duration.



## <a name="ToMillis">func</a> [ToMillis](/tofrom.go?s=400:432#L6)
``` go
func ToMillis(t time.Time) int64
```
ToMillis returns the number of milliseconds since the Epoch for the provided
time.Time value t.




## <a name="RetryFunc">type</a> [RetryFunc](/backoff.go?s=2030:2061#L45)
``` go
type RetryFunc func(i int) bool
```
RetryFunc is a function passed to RetryWithBackoff that should be retried
until successful. It should return true if the operation was successful or
false if it should be retried after a brief delay. The function will be
passed a single int argument which is the current (zero based) iteration
number.












