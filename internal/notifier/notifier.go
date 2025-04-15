package notifier

import "time"

type Notifier struct {
	articles         ArticleProvider
	summarizer       Summarizer
	bot              *tgbotapi.BotAPI
	fetchInterval    time.Duration
	lookupTimeWindow time.Duration
	chanelID         int64
}
