package calendar

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jinzhu/now"
)

type calData struct {
	removeMessage  bool
	cancelButton   bool
	doneTime       bool
	doneMsg        string
	loc            *time.Location
	topic          *int
	Year           int
	Month          string
	Day            int
	Hour           int
	Min            int
	startDate      time.Time
	acceptBackward bool
}

func (d calData) Time() time.Time {
	return time.Date(d.Year, time.Month(idxMonth(d.Month)+1), d.Day, d.Hour, d.Min, 0, 0, d.loc)
}

var months = []string{"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

func nextMonth(curM string) (string, int) {
	idx := slices.Index(months, curM)
	if (idx + 1) == len(months) {
		return months[0], 1
	}
	return months[idx+1], 0
}
func prevMonth(curM string) (string, int) {
	idx := slices.Index(months, curM)
	if idx == 0 {
		return months[len(months)-1], -1
	}
	return months[idx-1], 0
}

func idxMonth(curM string) int {
	return slices.Index(months, curM)
}

func makeKB(cd *calData, pref string) *models.InlineKeyboardMarkup {
	date := func(from int, to int, sp string) []models.InlineKeyboardButton {
		ikb := make([]models.InlineKeyboardButton, 0, to-from)
		for i := from; i <= to; i++ {
			txt := fmt.Sprintf("%d", i)
			if cd.Day == i {
				txt = fmt.Sprintf("%d ✅", i)
			}
			ikb = append(ikb, models.InlineKeyboardButton{
				Text:         txt,
				CallbackData: fmt.Sprintf("%s_%s%d", pref, sp, i),
			})
		}
		return ikb
	}
	hours := func(from int, to int, sp string) []models.InlineKeyboardButton {
		ikb := make([]models.InlineKeyboardButton, 0, to-from)
		for i := from; i <= to; i++ {
			txt := fmt.Sprintf("%d", i)
			if cd.Hour == i {
				txt = fmt.Sprintf("%d ✅", i)
			}
			ikb = append(ikb, models.InlineKeyboardButton{
				Text:         txt,
				CallbackData: fmt.Sprintf("%s_%s%d", pref, sp, i),
			})
		}
		return ikb
	}
	minutes := func(from int, to int, sp string) []models.InlineKeyboardButton {
		ikb := make([]models.InlineKeyboardButton, 0, to-from)
		for i := from; i <= to; i += 15 {
			txt := fmt.Sprintf("%d", i)
			if cd.Min == i {
				txt = fmt.Sprintf("%d ✅", i)
			}
			ikb = append(ikb, models.InlineKeyboardButton{
				Text:         txt,
				CallbackData: fmt.Sprintf("%s_%s%d", pref, sp, i),
			})
		}
		return ikb
	}
	tn := time.Now()
	var iks [][]models.InlineKeyboardButton

	//region Month
	var kbMonth []models.InlineKeyboardButton
	if cd.Year == tn.Year() {
		if cd.Month != tn.Month().String() || cd.acceptBackward {
			kbMonth = append(kbMonth, models.InlineKeyboardButton{Text: "<", CallbackData: pref + "_left"})
		}
	} else {
		kbMonth = append(kbMonth, models.InlineKeyboardButton{Text: "<", CallbackData: pref + "_left"})
	}
	kbMonth = append(kbMonth, models.InlineKeyboardButton{Text: fmt.Sprintf("%s %d", cd.Month, cd.Year), CallbackData: pref + "_name"})
	kbMonth = append(kbMonth, models.InlineKeyboardButton{Text: ">", CallbackData: pref + "_right"})
	iks = append(iks, kbMonth)
	//endregion

	if cd.Year == tn.Year() && cd.Month == tn.Month().String() {
		startDay := cd.Day
		if cd.acceptBackward {
			startDay = 1
		}
		lastDay := now.With(tn).EndOfMonth().Day()
		if (lastDay-startDay)/6 == 0 {
			iks = append(iks, date(startDay, lastDay, "d"))
		}
		wks := (lastDay - startDay) / 6
		if (lastDay-startDay)%6 > 0 {
			wks++
		}
		for i := 0; i < wks; i++ {
			if startDay+6 >= lastDay {
				iks = append(iks, date(startDay, lastDay, "d"))
				continue
			}
			iks = append(iks, date(startDay, startDay+5, "d"))
			startDay = startDay + 6
		}
	} else {
		startDay := 1
		mnt := time.Month(idxMonth(cd.Month) + 1)
		if mnt+1 > 12 {
			mnt = 1
		}
		lastDay := time.Date(cd.Year, mnt+1, 1, 0, 0, 0, 0, time.Local).Add(-time.Hour * 24).Day()
		cnt := (lastDay - startDay) / 6
		if (lastDay-startDay)%6 > 0 {
			cnt++
		}

		for i := 0; i < cnt; i++ {
			if startDay+6 >= lastDay {
				iks = append(iks, date(startDay, lastDay, "d"))
				continue
			}
			iks = append(iks, date(startDay, startDay+5, "d"))
			startDay = startDay + 6
		}
	}

	finalButtons := []models.InlineKeyboardButton{{Text: "Done", CallbackData: pref + "_final"}}
	if cd.cancelButton {
		finalButtons = []models.InlineKeyboardButton{{Text: "Done", CallbackData: pref + "_final"}, {Text: "Cancel", CallbackData: pref + "_cancel"}}
	}

	iks = append(iks,
		[]models.InlineKeyboardButton{{Text: "Hours", CallbackData: pref + "_name"}},
		hours(0, 6, "h"),
		hours(7, 12, "h"),
		hours(13, 18, "h"),
		hours(19, 23, "h"),
		[]models.InlineKeyboardButton{{Text: "Minutes", CallbackData: pref + "_name"}},
		minutes(0, 45, "m"),
		finalButtons)

	kb := &models.InlineKeyboardMarkup{InlineKeyboard: iks}
	return kb
}

func Calendar(b *bot.Bot, chatID int64, text string, pref string, options ...Option) (chan CalendarOrErr, error) {
	ch := make(chan CalendarOrErr)
	cd := new(calData)
	cd.startDate = time.Now()
	for _, opt := range options {
		opt(cd)
	}

	if cd.loc == nil {
		cd.loc = time.Local
	}
	cd.Year = cd.startDate.Year()
	cd.Month = cd.startDate.Month().String()
	cd.Day = cd.startDate.Day()
	cd.Hour = cd.startDate.Hour()
	cd.Min = cd.startDate.Minute()

	kb := makeKB(cd, pref)
	smp := &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: kb,
	}
	if cd.topic != nil {
		smp.MessageThreadID = *cd.topic
	}
	msg, err := b.SendMessage(context.Background(), smp)
	if err != nil {
		return ch, err
	}

	var es string
	var handlerString = &es
	var ctx, cancel = context.WithCancel(context.Background())
	*handlerString = b.RegisterHandlerMatchFunc(
		func(u *models.Update) bool {
			if u.CallbackQuery == nil {
				return false
			}
			if strings.HasPrefix(u.CallbackQuery.Data, pref) {
				return true
			}
			return false
		},
		func(ctx context.Context, b *bot.Bot, u *models.Update) {
			log := slog.With("method", "Calendar", "prefix", pref)
			var final bool
			var sendAnswer = true
			defer func(sa *bool) {
				if *sa {
					_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: u.CallbackQuery.ID})
				}
			}(&sendAnswer)
			dt := strings.TrimPrefix(u.CallbackQuery.Data, pref+"_")
			switch dt {
			case "left":
				mn, yr := prevMonth(cd.Month)
				cd.Month = mn
				cd.Year += yr
			case "name":
				sendAnswer = false
				_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: u.CallbackQuery.ID, Text: cd.Time().String(), ShowAlert: true})
			case "right":
				mn, yr := nextMonth(cd.Month)
				cd.Month = mn
				cd.Year += yr
			case "final":
				if time.Since(cd.Time()) > 0 && !cd.acceptBackward {
					_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
						CallbackQueryID: u.CallbackQuery.ID,
						Text:            t,
						ShowAlert:       true,
					})
					return
				}
				final = true
				b.UnregisterHandler(*handlerString)
				cancel()
				ch <- CalendarOrErr{CalendarDate: cd}
				close(ch)
			case "cancel":
				final = true
				b.UnregisterHandler(*handlerString)
				cancel()
				ch <- CalendarOrErr{Error: ErrorCancelled}
				close(ch)
			default:
				switch {
				case strings.HasPrefix(dt, "d"):
					// date
					d, err := strconv.Atoi(strings.TrimPrefix(dt, "d"))
					if err != nil {
						log.Error("Date incorrect", "date", dt, "error", err)
						return
					}
					cd.Day = d
				case strings.HasPrefix(dt, "h"):
					d, err := strconv.Atoi(strings.TrimPrefix(dt, "h"))
					if err != nil {
						log.Error("Hour incorrect", "data", dt, "err", err)
						return
					}
					cd.Hour = d
				case strings.HasPrefix(dt, "m"):
					d, err := strconv.Atoi(strings.TrimPrefix(dt, "m"))
					if err != nil {
						log.Error("Minute incorrect", "data", dt, "err", err)
						return
					}
					cd.Min = d
				}
			}
			if final {
				if cd.removeMessage {
					_, err = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
						ChatID:    msg.Chat.ID,
						MessageID: msg.ID,
					})
					if err != nil {
						log.Error("DeleteMessage", "err", err, "msg-id", msg.ID, "chat-id", msg.Chat.ID)
					}
				} else if cd.doneTime {
					_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
						ChatID:      msg.Chat.ID,
						MessageID:   msg.ID,
						Text:        fmt.Sprintf(t1+" %s\n%s", cd.Time().Format("2006-01-02 15:04:05"), cd.doneMsg),
						ReplyMarkup: nil,
					})
					if err != nil {
						log.Error("Edit message text", "err", err, "msg-id", msg.ID, "chat-id", msg.Chat.ID)
					}
				} else {
					_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
						ChatID:      u.CallbackQuery.Message.Message.Chat.ID,
						MessageID:   u.CallbackQuery.Message.Message.ID,
						ReplyMarkup: nil,
					})
					if err != nil {
						log.Error("Edit message buttons", "err", err, "chat-id", msg.Chat.ID)
					}
				}

				return
			}
			_, _ = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
				ChatID:      u.CallbackQuery.Message.Message.Chat.ID,
				MessageID:   u.CallbackQuery.Message.Message.ID,
				ReplyMarkup: makeKB(cd, pref),
			})
		},
	)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.Tick(time.Hour):
				if cd.removeMessage {
					_, _ = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
						ChatID:    msg.Chat.ID,
						MessageID: msg.ID,
					})
				} else {
					_, _ = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
						ChatID:    msg.Chat.ID,
						MessageID: msg.ID,
						ReplyMarkup: &models.InlineKeyboardMarkup{
							InlineKeyboard: [][]models.InlineKeyboardButton{
								{{Text: t2, CallbackData: "none"}},
							},
						},
					})
				}

				b.UnregisterHandler(*handlerString)
				cancel()
				ch <- CalendarOrErr{Error: ErrorTimeout}
				close(ch)
			}
		}
	}()

	return ch, nil
}
