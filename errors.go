// Copyright © 2024 Timothy E. Peoples

package timetool

type Error string

func (e Error) Error() string {
	return string(e)
}

//╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴
// Ticker related errors.

// ErrTickerActive is returned by NormalTicker's Err method if the ticker is
// still active (i.e. it has not been stopped).
var ErrTickerActive = Error("ticker is active")

//╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴
// Backoff related errors.

// ErrNilReceiver is returned when a method is called using a nil pointer
// reciever.
var ErrNilReceiver = Error("nil reciver")

// ErrMissingDeadline is returned by functions expecting a Deadline
// but are passed a Context devoid of such nature.
var ErrMissingDeadline = Error("context must have a deadline")

// ErrTooFewIterations is returned if the number of requested iterations
// is less than 2.
var ErrTooFewIterations = Error("number of iterations must be at least 2")

// ErrTimeWarp is returned by RetryWithBackoff in the freakishly uncommon
// event that a context ends up with a deadline in the past and a Done
// channel that does not return.
var ErrTimeWarp = Error("valid context has deadline in the past")

// ErrRetriesExhausted is returned by RetryWithBackoff if all retry
// attempts have been unsuccessful.
var ErrRetriesExhausted = Error("all retries exhausted")

// ErrBadJitter is returned when an invalid jitter value has been requested.
var ErrBadJitter = Error("invalid jitter value; must be [0.0, 100.0)")

// ErrNegativeDelay is returned when an applicable time.Duration value is
// negative.
var ErrNegativeDelay = Error("negative delay value; time travel not yet supported")

// ErrZeroCoefficient is returned when a Backoff.Coefficient value is zero.
var ErrZeroCoefficient = Error("coefficient cannot be zero")

//╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴╶╴
