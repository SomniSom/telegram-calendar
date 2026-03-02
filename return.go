package calendar

import "errors"

//goland:noinspection GoNameStartsWithPackageName
type CalendarOrErr struct {
	Error        error
	CalendarDate *calData
}

var ErrorCancelled = errors.New("cancelled")
var ErrorTimeout = errors.New("timeout")
