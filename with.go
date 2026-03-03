package calendar

import "time"

// Option is a functional option type used to configure calData instances.
type Option func(c *calData)

// WithTopic returns an Option that sets the topic identifier for the calendar data.
func WithTopic(topic int) Option {
	return func(c *calData) {
		c.topic = &topic
	}
}

// WithTimeLocation sets the time zone used for the calendar date and time.
// If not provided, the local time zone is used.
func WithTimeLocation(tl *time.Location) Option {
	return func(c *calData) {
		c.loc = tl
	}
}

// WithStartDate sets the initial date and time for the calendar.
func WithStartDate(t time.Time) Option {
	return func(c *calData) {
		c.startDate = t
		c.Year = t.Year()
		c.Month = t.Month().String()
		c.Day = t.Day()
		c.Hour = t.Hour()
		c.Min = t.Minute()
		c.loc = t.Location()
	}
}

// WithReplaceDoneTime enables replacing the calendar message with custom text upon final selection.// WithReplaceDoneTime defines a custom message text to replace the calendar message upon user confirmation.// WithReplaceDoneTime enables replacing the calendar message with custom text after final selection.
// The b argument controls whether replacement occurs; msg provides a function to generate the text from the selected time.
func WithReplaceDoneTime(b bool, msg func(time.Time) string) Option {
	return func(c *calData) {
		c.doneTime = b
		c.doneMsg = msg
	}
}

// WithRemoveMessageAfterDone configures removal of the message after the done event.
func WithRemoveMessageAfterDone(b bool) Option {
	return func(c *calData) {
		c.removeMessage = b
	}
}

// WithCancelButton adds a cancel button to the calendar UI, allowing the user to abort selection.
func WithCancelButton() Option {
	return func(c *calData) {
		c.cancelButton = true
	}
}

// WithAcceptBackward enables backward date selection in the calendar.
func WithAcceptBackward() Option {
	return func(c *calData) {
		c.acceptBackward = true
	}
}
