package notifier

import (
	"context"
	"io"
	"strings"
	"tg-bot-supchick/internal/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ArticleProvider interface {
	AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]models.Article, error)
	MarkAsPosted(ctx context.Context, id int64) error
}

type Summarizer interface {
	Summarize(ctx context.Context, since time.Time, text string) (string, error)
	ExtractSummary(ctx context.Context, text string) (string, error)
}

type Notifier struct {
	articles         ArticleProvider
	summarizer       Summarizer
	bot              *tgbotapi.BotAPI
	fetchInterval    time.Duration
	lookupTimeWindow time.Duration
	channelID        int64
}

func NewNotifier(
	articles ArticleProvider,
	summarizer Summarizer,
	bot *tgbotapi.BotAPI,
	fetchInterval time.Duration,
	lookupTimeWindow time.Duration,
	channelID int64,
) *Notifier {

	return &Notifier{
		articles:         articles,
		summarizer:       summarizer,
		bot:              bot,
		fetchInterval:    fetchInterval,
		lookupTimeWindow: lookupTimeWindow,
		channelID:        channelID,
	}
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	articles, err := n.articles.AllNotPosted(ctx, time.Now().Add(-n.lookupTimeWindow), 1)
	if err != nil {
		return err
	}

	if len(articles) == 0 {
		return nil
	}

	article := articles[0]
	summary, err := n.summarizer.ExtractSummary(ctx, article.Summary)
	if err != nil {
		return err
	}

	if err := n.sendArticle(ctx, article, summary); err != nil {
		return err
	}

	return n.articles.MarkAsPosted(ctx, article.ID)
}

func (n *Notifier) ExtractSummary(ctx context.Context, a models.Article) (string, error) {
	var r io.Reader

	if a.Summary != "" {
		r = strings.NewReader(a.Summary)
	} else {

	}
}
