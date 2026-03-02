package calendar

import "time"

type Option func(c *calData)

func WithTopic(topic int) Option {
	return func(c *calData) {
		c.topic = &topic
	}
}

func WithTimeLocation(tl *time.Location) Option {
	return func(c *calData) {
		c.loc = tl
	}
}

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

func WithReplaceDoneTime(b bool, msg string) Option {
	return func(c *calData) {
		c.doneTime = b
		c.doneMsg = msg
	}
}

func WithRemoveMessageAfterDone(b bool) Option {
	return func(c *calData) {
		c.removeMessage = b
	}
}

func WithCancelButton() Option {
	return func(c *calData) {
		c.cancelButton = true
	}
}

func WithAcceptBackward() Option {
	return func(c *calData) {
		c.acceptBackward = true
	}
}
