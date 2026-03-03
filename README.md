# Calendar for telegram bot

Use for this: github.com/go-telegram/bot

## Install
```shell
go get github.com/SomniSom/telegram-calendar
```

## Example usage

Params:
* `text` - start text message
* `pref` - prefix for callback identity

```go
package main

import (
	"errors"
	"fmt"
	"time"

	calendar "github.com/SomniSom/telegram-calendar"
	"github.com/go-telegram/bot"
)

const botToken = `1111111111:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA`

func main() {
	myLoc, _ := time.LoadLocation(`Europe/Berlin`)
	tm := time.Now().In(myLoc)
	b, _ := bot.New(botToken)
	ch, err := calendar.Calendar(b, 111111111, "test", "calldr",
		calendar.WithReplaceDoneTime(true, func(t time.Time) string {
			return fmt.Sprintf("la-la-la\n%s", t.String())
		}),
		calendar.WithStartDate(tm),
		calendar.WithCancelButton(),
		calendar.WithTimeLocation(myLoc),
	)
	if err != nil {
		panic(err)
	}
	res := <-ch
	if res.Error != nil {
		if errors.Is(res.Error, calendar.ErrorTimeout) {
			//on timeout action
		}
		if errors.Is(res.Error, calendar.ErrorCancelled) {
			//on cancel action
		}
		panic(res.Error)
	}
	fmt.Println(res.CalendarDate.Time().String())
}
```

### All options
* `WithReplaceDoneTime` - defines a custom message text to replace the calendar message upon user confirmation.
* `WithStartDate` - sets the initial date and time for the calendar.
* `WithCancelButton` - adds a cancel button to the calendar UI, allowing the user to abort selection.
* `WithTimeLocation` - sets the time zone used for the calendar date and time. If not provided, the local time zone is used.
* `WithTopic` - hat sets the topic identifier for the calendar data.
* `WithRemoveMessageAfterDone` - configures removal of the message after the done event.
* `WithAcceptBackward` - enables backward date selection in the calendar.