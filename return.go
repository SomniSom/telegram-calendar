package calendar

import "errors"

// CalendarOrErr represents the result of a calendar interaction,
// containing either a selected date or an error such as cancellation or timeout.
type CalendarOrErr struct {
	Error        error
	CalendarDate *calData
}

// ErrorCancelled is returned when the user cancels a calendar operation.
var ErrorCancelled = errors.New("cancelled")

// ErrorTimeout is returned when the calendar operation times out.
var ErrorTimeout = errors.New("timeout")
